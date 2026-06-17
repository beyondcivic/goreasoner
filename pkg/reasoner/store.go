package reasoner

import (
	"fmt"
	"strings"
)

// Triple represents an RDF triple (subject, predicate, object).
type Triple struct {
	Subject   string
	Predicate string
	Object    string
}

// String returns the triple in N-Triples format.
func (t Triple) String() string {
	subj := formatTerm(t.Subject)
	pred := formatTerm(t.Predicate)
	obj := formatTerm(t.Object)
	return fmt.Sprintf("%s %s %s .", subj, pred, obj)
}

// formatTerm formats a term for output.
func formatTerm(term string) string {
	if strings.HasPrefix(term, "http://") || strings.HasPrefix(term, "https://") {
		return "<" + term + ">"
	}
	if strings.HasPrefix(term, "<") && strings.HasSuffix(term, ">") {
		return term
	}
	if strings.HasPrefix(term, "\"") {
		return term
	}
	if strings.Contains(term, ":") && !strings.HasPrefix(term, "_:") {
		return term
	}
	return term
}

// wellKnownIDs caches interned IDs for RDF/RDFS/OWL URI constants
// so the fast reasoning path never has to hash these strings.
type wellKnownIDs struct {
	RDFType        uint32
	SubClassOf     uint32
	SubPropertyOf  uint32
	Domain         uint32
	Range          uint32
	EquivClass     uint32
	SameAs         uint32
	InverseOf      uint32
	TransitiveProp uint32
	SymmetricProp  uint32
}

// TripleStore is an in-memory store for RDF triples.
// It wraps FastStore + Intern for compact storage while keeping
// the original string-based public API.
type TripleStore struct {
	intern *Intern
	fast   *FastStore
	ids    wellKnownIDs
}

// NewTripleStore creates a new empty triple store.
func NewTripleStore() *TripleStore {
	return newTripleStoreWithCapacity(256)
}

// newTripleStoreWithCapacity creates a triple store pre-sized for n triples.
func newTripleStoreWithCapacity(n int) *TripleStore {
	ts := &TripleStore{
		intern: NewIntern(),
		fast:   NewFastStore(n),
	}
	ts.initWellKnownIDs()
	return ts
}

// initWellKnownIDs pre-interns all URI constants from rules.go.
func (ts *TripleStore) initWellKnownIDs() {
	ts.ids = wellKnownIDs{
		RDFType:        ts.intern.ID(RDFType),
		SubClassOf:     ts.intern.ID(RDFSSubClassOf),
		SubPropertyOf:  ts.intern.ID(RDFSSubPropertyOf),
		Domain:         ts.intern.ID(RDFSDomain),
		Range:          ts.intern.ID(RDFSRange),
		EquivClass:     ts.intern.ID(OWLEquivalentClass),
		SameAs:         ts.intern.ID(OWLSameAs),
		InverseOf:      ts.intern.ID(OWLInverseOf),
		TransitiveProp: ts.intern.ID(OWLTransitiveProperty),
		SymmetricProp:  ts.intern.ID(OWLSymmetricProperty),
	}
}

// Add adds a triple to the store, returns true if it was new.
func (ts *TripleStore) Add(t Triple) bool {
	ct := compactTriple{
		S: ts.intern.ID(t.Subject),
		P: ts.intern.ID(t.Predicate),
		O: ts.intern.ID(t.Object),
	}
	return ts.fast.Add(ct)
}

// Contains checks if a triple exists in the store.
func (ts *TripleStore) Contains(t Triple) bool {
	s, okS := ts.intern.toID[t.Subject]
	p, okP := ts.intern.toID[t.Predicate]
	o, okO := ts.intern.toID[t.Object]
	if !okS || !okP || !okO {
		return false
	}
	return ts.fast.Contains(compactTriple{S: s, P: p, O: o})
}

// FindBySubject returns all triples with the given subject.
func (ts *TripleStore) FindBySubject(subject string) []Triple {
	id, ok := ts.intern.toID[subject]
	if !ok {
		return nil
	}
	idxs := ts.fast.ByS(id)
	result := make([]Triple, len(idxs))
	for i, idx := range idxs {
		result[i] = ts.toTriple(ts.fast.all[idx])
	}
	return result
}

// FindByPredicate returns all triples with the given predicate.
func (ts *TripleStore) FindByPredicate(predicate string) []Triple {
	id, ok := ts.intern.toID[predicate]
	if !ok {
		return nil
	}
	idxs := ts.fast.ByP(id)
	result := make([]Triple, len(idxs))
	for i, idx := range idxs {
		result[i] = ts.toTriple(ts.fast.all[idx])
	}
	return result
}

// FindByObject returns all triples with the given object.
func (ts *TripleStore) FindByObject(object string) []Triple {
	id, ok := ts.intern.toID[object]
	if !ok {
		return nil
	}
	// No dedicated byO index in FastStore; scan byP entries.
	// This keeps the hot path (BySP, ByPO, ByP, ByS) fast.
	var result []Triple
	for _, ct := range ts.fast.all {
		if ct.O == id {
			result = append(result, ts.toTriple(ct))
		}
	}
	return result
}

// FindBySubjectPredicate returns all triples matching subject and predicate.
func (ts *TripleStore) FindBySubjectPredicate(subject, predicate string) []Triple {
	s, okS := ts.intern.toID[subject]
	p, okP := ts.intern.toID[predicate]
	if !okS || !okP {
		return nil
	}
	idxs := ts.fast.BySP(s, p)
	result := make([]Triple, len(idxs))
	for i, idx := range idxs {
		result[i] = ts.toTriple(ts.fast.all[idx])
	}
	return result
}

// FindByPredicateObject returns all triples matching predicate and object.
func (ts *TripleStore) FindByPredicateObject(predicate, object string) []Triple {
	p, okP := ts.intern.toID[predicate]
	o, okO := ts.intern.toID[object]
	if !okP || !okO {
		return nil
	}
	idxs := ts.fast.ByPO(p, o)
	result := make([]Triple, len(idxs))
	for i, idx := range idxs {
		result[i] = ts.toTriple(ts.fast.all[idx])
	}
	return result
}

// All returns all triples in the store.
func (ts *TripleStore) All() []Triple {
	result := make([]Triple, len(ts.fast.all))
	for i, ct := range ts.fast.all {
		result[i] = ts.toTriple(ct)
	}
	return result
}

// Size returns the number of triples in the store.
func (ts *TripleStore) Size() int {
	return ts.fast.Size()
}

// toTriple converts a compactTriple back to a public Triple.
func (ts *TripleStore) toTriple(ct compactTriple) Triple {
	return Triple{
		Subject:   ts.intern.Str(ct.S),
		Predicate: ts.intern.Str(ct.P),
		Object:    ts.intern.Str(ct.O),
	}
}
