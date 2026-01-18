package reasoner

import (
	"fmt"
	"strings"
)

// Triple represents an RDF triple (subject, predicate, object)
type Triple struct {
	Subject   string
	Predicate string
	Object    string
}

// String returns the triple in N-Triples format
func (t Triple) String() string {
	subj := formatTerm(t.Subject)
	pred := formatTerm(t.Predicate)
	obj := formatTerm(t.Object)
	return fmt.Sprintf("%s %s %s .", subj, pred, obj)
}

// formatTerm formats a term for output
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

// TripleStore is an in-memory store for RDF triples
type TripleStore struct {
	triples    map[string]bool
	tripleList []Triple

	// Indexes for fast lookup
	bySubject   map[string][]int
	byPredicate map[string][]int
	byObject    map[string][]int
}

// NewTripleStore creates a new empty triple store
func NewTripleStore() *TripleStore {
	return &TripleStore{
		triples:     make(map[string]bool),
		tripleList:  make([]Triple, 0),
		bySubject:   make(map[string][]int),
		byPredicate: make(map[string][]int),
		byObject:    make(map[string][]int),
	}
}

// tripleKey generates a unique key for a triple
func tripleKey(t Triple) string {
	return t.Subject + "|" + t.Predicate + "|" + t.Object
}

// Add adds a triple to the store, returns true if it was new
func (ts *TripleStore) Add(t Triple) bool {
	key := tripleKey(t)
	if ts.triples[key] {
		return false
	}

	ts.triples[key] = true
	idx := len(ts.tripleList)
	ts.tripleList = append(ts.tripleList, t)

	ts.bySubject[t.Subject] = append(ts.bySubject[t.Subject], idx)
	ts.byPredicate[t.Predicate] = append(ts.byPredicate[t.Predicate], idx)
	ts.byObject[t.Object] = append(ts.byObject[t.Object], idx)

	return true
}

// Contains checks if a triple exists in the store
func (ts *TripleStore) Contains(t Triple) bool {
	return ts.triples[tripleKey(t)]
}

// FindBySubject returns all triples with the given subject
func (ts *TripleStore) FindBySubject(subject string) []Triple {
	var result []Triple
	for _, idx := range ts.bySubject[subject] {
		result = append(result, ts.tripleList[idx])
	}
	return result
}

// FindByPredicate returns all triples with the given predicate
func (ts *TripleStore) FindByPredicate(predicate string) []Triple {
	var result []Triple
	for _, idx := range ts.byPredicate[predicate] {
		result = append(result, ts.tripleList[idx])
	}
	return result
}

// FindByObject returns all triples with the given object
func (ts *TripleStore) FindByObject(object string) []Triple {
	var result []Triple
	for _, idx := range ts.byObject[object] {
		result = append(result, ts.tripleList[idx])
	}
	return result
}

// FindBySubjectPredicate returns all triples matching subject and predicate
func (ts *TripleStore) FindBySubjectPredicate(subject, predicate string) []Triple {
	var result []Triple
	for _, idx := range ts.bySubject[subject] {
		t := ts.tripleList[idx]
		if t.Predicate == predicate {
			result = append(result, t)
		}
	}
	return result
}

// FindByPredicateObject returns all triples matching predicate and object
func (ts *TripleStore) FindByPredicateObject(predicate, object string) []Triple {
	var result []Triple
	for _, idx := range ts.byPredicate[predicate] {
		t := ts.tripleList[idx]
		if t.Object == object {
			result = append(result, t)
		}
	}
	return result
}

// All returns all triples in the store
func (ts *TripleStore) All() []Triple {
	result := make([]Triple, len(ts.tripleList))
	copy(result, ts.tripleList)
	return result
}

// Size returns the number of triples in the store
func (ts *TripleStore) Size() int {
	return len(ts.tripleList)
}