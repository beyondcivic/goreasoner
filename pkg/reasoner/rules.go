package reasoner

// Common RDF/RDFS/OWL URIs
const (
	RDFType               = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
	RDFSSubClassOf        = "http://www.w3.org/2000/01/rdf-schema#subClassOf"
	RDFSSubPropertyOf     = "http://www.w3.org/2000/01/rdf-schema#subPropertyOf"
	RDFSDomain            = "http://www.w3.org/2000/01/rdf-schema#domain"
	RDFSRange             = "http://www.w3.org/2000/01/rdf-schema#range"
	OWLClass              = "http://www.w3.org/2002/07/owl#Class"
	OWLThing              = "http://www.w3.org/2002/07/owl#Thing"
	OWLEquivalentClass    = "http://www.w3.org/2002/07/owl#equivalentClass"
	OWLSameAs             = "http://www.w3.org/2002/07/owl#sameAs"
	OWLInverseOf          = "http://www.w3.org/2002/07/owl#inverseOf"
	OWLTransitiveProperty = "http://www.w3.org/2002/07/owl#TransitiveProperty"
	OWLSymmetricProperty  = "http://www.w3.org/2002/07/owl#SymmetricProperty"
)

// applyAllRules implements semi-naive forward chaining over delta triples.
// For each new triple in delta, it matches against the full store to derive
// further inferred triples. New triples are added to the store and appended
// to *out (the next round's delta).
func applyAllRules(store *TripleStore, delta []compactTriple, out *[]compactTriple) {
	id := store.ids
	fs := store.fast

	for _, d := range delta {
		// ── rdfs:subClassOf transitivity ──
		if d.P == id.SubClassOf {
			for _, idx := range fs.BySP(d.O, id.SubClassOf) {
				c := fs.all[idx].O
				if d.S != c {
					emitCompact(store, compactTriple{d.S, id.SubClassOf, c}, out)
				}
			}
			for _, idx := range fs.ByPO(id.SubClassOf, d.S) {
				x := fs.all[idx].S
				if x != d.O {
					emitCompact(store, compactTriple{x, id.SubClassOf, d.O}, out)
				}
			}
		}

		// ── rdf:type inheritance ──
		if d.P == id.RDFType {
			for _, idx := range fs.BySP(d.O, id.SubClassOf) {
				b := fs.all[idx].O
				emitCompact(store, compactTriple{d.S, id.RDFType, b}, out)
			}
		}
		if d.P == id.SubClassOf {
			for _, idx := range fs.ByPO(id.RDFType, d.S) {
				x := fs.all[idx].S
				emitCompact(store, compactTriple{x, id.RDFType, d.O}, out)
			}
		}

		// ── rdfs:subPropertyOf transitivity ──
		if d.P == id.SubPropertyOf {
			for _, idx := range fs.BySP(d.O, id.SubPropertyOf) {
				p3 := fs.all[idx].O
				if d.S != p3 {
					emitCompact(store, compactTriple{d.S, id.SubPropertyOf, p3}, out)
				}
			}
			for _, idx := range fs.ByPO(id.SubPropertyOf, d.S) {
				x := fs.all[idx].S
				if x != d.O {
					emitCompact(store, compactTriple{x, id.SubPropertyOf, d.O}, out)
				}
			}
		}

		// ── rdfs:subPropertyOf inheritance ──
		for _, idx := range fs.BySP(d.P, id.SubPropertyOf) {
			p2 := fs.all[idx].O
			emitCompact(store, compactTriple{d.S, p2, d.O}, out)
		}
		if d.P == id.SubPropertyOf {
			for _, idx := range fs.ByP(d.S) {
				t := fs.all[idx]
				emitCompact(store, compactTriple{t.S, d.O, t.O}, out)
			}
		}

		// ── rdfs:domain ──
		for _, idx := range fs.BySP(d.P, id.Domain) {
			c := fs.all[idx].O
			emitCompact(store, compactTriple{d.S, id.RDFType, c}, out)
		}
		if d.P == id.Domain {
			for _, idx := range fs.ByP(d.S) {
				t := fs.all[idx]
				emitCompact(store, compactTriple{t.S, id.RDFType, d.O}, out)
			}
		}

		// ── rdfs:range ──
		for _, idx := range fs.BySP(d.P, id.Range) {
			c := fs.all[idx].O
			emitCompact(store, compactTriple{d.O, id.RDFType, c}, out)
		}
		if d.P == id.Range {
			for _, idx := range fs.ByP(d.S) {
				t := fs.all[idx]
				emitCompact(store, compactTriple{t.O, id.RDFType, d.O}, out)
			}
		}

		// ── owl:SymmetricProperty ──
		if fs.Contains(compactTriple{d.P, id.RDFType, id.SymmetricProp}) {
			emitCompact(store, compactTriple{d.O, d.P, d.S}, out)
		}
		if d.P == id.RDFType && d.O == id.SymmetricProp {
			for _, idx := range fs.ByP(d.S) {
				t := fs.all[idx]
				emitCompact(store, compactTriple{t.O, d.S, t.S}, out)
			}
		}

		// ── owl:TransitiveProperty ──
		if fs.Contains(compactTriple{d.P, id.RDFType, id.TransitiveProp}) {
			for _, idx := range fs.BySP(d.O, d.P) {
				z := fs.all[idx].O
				if d.S != z {
					emitCompact(store, compactTriple{d.S, d.P, z}, out)
				}
			}
			for _, idx := range fs.ByPO(d.P, d.S) {
				w := fs.all[idx].S
				if w != d.O {
					emitCompact(store, compactTriple{w, d.P, d.O}, out)
				}
			}
		}

		// ── owl:inverseOf ──
		for _, idx := range fs.BySP(d.P, id.InverseOf) {
			p2 := fs.all[idx].O
			emitCompact(store, compactTriple{d.O, p2, d.S}, out)
		}
		for _, idx := range fs.ByPO(id.InverseOf, d.P) {
			p1 := fs.all[idx].S
			emitCompact(store, compactTriple{d.O, p1, d.S}, out)
		}
		if d.P == id.InverseOf {
			for _, idx := range fs.ByP(d.S) {
				t := fs.all[idx]
				emitCompact(store, compactTriple{t.O, d.O, t.S}, out)
			}
			for _, idx := range fs.ByP(d.O) {
				t := fs.all[idx]
				emitCompact(store, compactTriple{t.O, d.S, t.S}, out)
			}
		}

		// ── owl:equivalentClass symmetry + transitivity ──
		if d.P == id.EquivClass {
			emitCompact(store, compactTriple{d.O, id.EquivClass, d.S}, out)
			for _, idx := range fs.BySP(d.O, id.EquivClass) {
				c := fs.all[idx].O
				if d.S != c {
					emitCompact(store, compactTriple{d.S, id.EquivClass, c}, out)
				}
			}
		}

		// ── owl:sameAs symmetry + transitivity ──
		if d.P == id.SameAs {
			emitCompact(store, compactTriple{d.O, id.SameAs, d.S}, out)
			for _, idx := range fs.BySP(d.O, id.SameAs) {
				c := fs.all[idx].O
				if d.S != c {
					emitCompact(store, compactTriple{d.S, id.SameAs, c}, out)
				}
			}
		}
	}
}

func emitCompact(store *TripleStore, t compactTriple, out *[]compactTriple) {
	if store.fast.Add(t) {
		*out = append(*out, t)
	}
}

// ── Legacy Rule interface (kept for backward compatibility) ──

// Deprecated: Rule is retained for API compatibility but is no longer used
// by the default reasoning engine. Use RunForwardReasoning() which employs
// semi-naive evaluation directly.
type Rule interface {
	Name() string
	Apply(store *TripleStore) []Triple
}

// Deprecated: DefaultRules returns nil. The new semi-naive engine (applyAllRules)
// is used by default and is significantly faster.
func DefaultRules() []Rule { return nil }
