package reasoner

// Common RDF/RDFS/OWL URIs
const (
	RDFType           = "http://www.w3.org/1999/02/22-rdf-syntax-ns#type"
	RDFSSubClassOf    = "http://www.w3.org/2000/01/rdf-schema#subClassOf"
	RDFSSubPropertyOf = "http://www.w3.org/2000/01/rdf-schema#subPropertyOf"
	RDFSDomain        = "http://www.w3.org/2000/01/rdf-schema#domain"
	RDFSRange         = "http://www.w3.org/2000/01/rdf-schema#range"
	OWLClass          = "http://www.w3.org/2002/07/owl#Class"
	OWLThing          = "http://www.w3.org/2002/07/owl#Thing"
	OWLEquivalentClass = "http://www.w3.org/2002/07/owl#equivalentClass"
	OWLSameAs         = "http://www.w3.org/2002/07/owl#sameAs"
	OWLInverseOf      = "http://www.w3.org/2002/07/owl#inverseOf"
	OWLTransitiveProperty = "http://www.w3.org/2002/07/owl#TransitiveProperty"
	OWLSymmetricProperty  = "http://www.w3.org/2002/07/owl#SymmetricProperty"
)

// Rule represents a forward reasoning rule
type Rule interface {
	// Name returns the rule name
	Name() string
	// Apply applies the rule to the store and returns new inferred triples
	Apply(store *TripleStore) []Triple
}

// SubClassTransitivity implements rdfs:subClassOf transitivity
// If A rdfs:subClassOf B and B rdfs:subClassOf C, then A rdfs:subClassOf C
type SubClassTransitivity struct{}

func (r *SubClassTransitivity) Name() string {
	return "rdfs:subClassOf-transitivity"
}

func (r *SubClassTransitivity) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	subClassTriples := store.FindByPredicate(RDFSSubClassOf)

	for _, t1 := range subClassTriples {
		// t1: A subClassOf B
		a := t1.Subject
		b := t1.Object

		// Find all: B subClassOf C
		for _, t2 := range store.FindBySubjectPredicate(b, RDFSSubClassOf) {
			c := t2.Object
			// Infer: A subClassOf C
			newTriple := Triple{Subject: a, Predicate: RDFSSubClassOf, Object: c}
			if !store.Contains(newTriple) && a != c {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// TypeInheritance implements rdf:type inheritance through subClassOf
// If X rdf:type A and A rdfs:subClassOf B, then X rdf:type B
type TypeInheritance struct{}

func (r *TypeInheritance) Name() string {
	return "rdf:type-inheritance"
}

func (r *TypeInheritance) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	typeTriples := store.FindByPredicate(RDFType)

	for _, t := range typeTriples {
		// t: X rdf:type A
		x := t.Subject
		a := t.Object

		// Find all: A subClassOf B
		for _, sc := range store.FindBySubjectPredicate(a, RDFSSubClassOf) {
			b := sc.Object
			// Infer: X rdf:type B
			newTriple := Triple{Subject: x, Predicate: RDFType, Object: b}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// DomainInference implements rdfs:domain inference
// If P rdfs:domain C and X P Y, then X rdf:type C
type DomainInference struct{}

func (r *DomainInference) Name() string {
	return "rdfs:domain-inference"
}

func (r *DomainInference) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	domainTriples := store.FindByPredicate(RDFSDomain)

	for _, dt := range domainTriples {
		// dt: P rdfs:domain C
		p := dt.Subject
		c := dt.Object

		// Find all: X P Y
		for _, t := range store.FindByPredicate(p) {
			x := t.Subject
			// Infer: X rdf:type C
			newTriple := Triple{Subject: x, Predicate: RDFType, Object: c}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// RangeInference implements rdfs:range inference
// If P rdfs:range C and X P Y, then Y rdf:type C
type RangeInference struct{}

func (r *RangeInference) Name() string {
	return "rdfs:range-inference"
}

func (r *RangeInference) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	rangeTriples := store.FindByPredicate(RDFSRange)

	for _, rt := range rangeTriples {
		// rt: P rdfs:range C
		p := rt.Subject
		c := rt.Object

		// Find all: X P Y
		for _, t := range store.FindByPredicate(p) {
			y := t.Object
			// Skip literals
			if len(y) > 0 && y[0] == '"' {
				continue
			}
			// Infer: Y rdf:type C
			newTriple := Triple{Subject: y, Predicate: RDFType, Object: c}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// SubPropertyTransitivity implements rdfs:subPropertyOf transitivity
// If P1 rdfs:subPropertyOf P2 and P2 rdfs:subPropertyOf P3, then P1 rdfs:subPropertyOf P3
type SubPropertyTransitivity struct{}

func (r *SubPropertyTransitivity) Name() string {
	return "rdfs:subPropertyOf-transitivity"
}

func (r *SubPropertyTransitivity) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	subPropTriples := store.FindByPredicate(RDFSSubPropertyOf)

	for _, t1 := range subPropTriples {
		p1 := t1.Subject
		p2 := t1.Object

		for _, t2 := range store.FindBySubjectPredicate(p2, RDFSSubPropertyOf) {
			p3 := t2.Object
			newTriple := Triple{Subject: p1, Predicate: RDFSSubPropertyOf, Object: p3}
			if !store.Contains(newTriple) && p1 != p3 {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// SubPropertyInheritance implements property inheritance
// If P1 rdfs:subPropertyOf P2 and X P1 Y, then X P2 Y
type SubPropertyInheritance struct{}

func (r *SubPropertyInheritance) Name() string {
	return "rdfs:subPropertyOf-inheritance"
}

func (r *SubPropertyInheritance) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	subPropTriples := store.FindByPredicate(RDFSSubPropertyOf)

	for _, sp := range subPropTriples {
		p1 := sp.Subject
		p2 := sp.Object

		for _, t := range store.FindByPredicate(p1) {
			newTriple := Triple{Subject: t.Subject, Predicate: p2, Object: t.Object}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// EquivalentClassSymmetry implements owl:equivalentClass symmetry
// If A owl:equivalentClass B, then B owl:equivalentClass A
type EquivalentClassSymmetry struct{}

func (r *EquivalentClassSymmetry) Name() string {
	return "owl:equivalentClass-symmetry"
}

func (r *EquivalentClassSymmetry) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	eqTriples := store.FindByPredicate(OWLEquivalentClass)

	for _, t := range eqTriples {
		newTriple := Triple{Subject: t.Object, Predicate: OWLEquivalentClass, Object: t.Subject}
		if !store.Contains(newTriple) {
			inferred = append(inferred, newTriple)
		}
	}

	return inferred
}

// EquivalentClassTransitivity implements owl:equivalentClass transitivity
// If A owl:equivalentClass B and B owl:equivalentClass C, then A owl:equivalentClass C
type EquivalentClassTransitivity struct{}

func (r *EquivalentClassTransitivity) Name() string {
	return "owl:equivalentClass-transitivity"
}

func (r *EquivalentClassTransitivity) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	eqTriples := store.FindByPredicate(OWLEquivalentClass)

	for _, t1 := range eqTriples {
		a := t1.Subject
		b := t1.Object

		for _, t2 := range store.FindBySubjectPredicate(b, OWLEquivalentClass) {
			c := t2.Object
			newTriple := Triple{Subject: a, Predicate: OWLEquivalentClass, Object: c}
			if !store.Contains(newTriple) && a != c {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// SameAsSymmetry implements owl:sameAs symmetry
type SameAsSymmetry struct{}

func (r *SameAsSymmetry) Name() string {
	return "owl:sameAs-symmetry"
}

func (r *SameAsSymmetry) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	sameAsTriples := store.FindByPredicate(OWLSameAs)

	for _, t := range sameAsTriples {
		newTriple := Triple{Subject: t.Object, Predicate: OWLSameAs, Object: t.Subject}
		if !store.Contains(newTriple) {
			inferred = append(inferred, newTriple)
		}
	}

	return inferred
}

// SameAsTransitivity implements owl:sameAs transitivity
type SameAsTransitivity struct{}

func (r *SameAsTransitivity) Name() string {
	return "owl:sameAs-transitivity"
}

func (r *SameAsTransitivity) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	sameAsTriples := store.FindByPredicate(OWLSameAs)

	for _, t1 := range sameAsTriples {
		a := t1.Subject
		b := t1.Object

		for _, t2 := range store.FindBySubjectPredicate(b, OWLSameAs) {
			c := t2.Object
			newTriple := Triple{Subject: a, Predicate: OWLSameAs, Object: c}
			if !store.Contains(newTriple) && a != c {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// InversePropertyInference implements owl:inverseOf
// If P1 owl:inverseOf P2 and X P1 Y, then Y P2 X
type InversePropertyInference struct{}

func (r *InversePropertyInference) Name() string {
	return "owl:inverseOf-inference"
}

func (r *InversePropertyInference) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	inverseTriples := store.FindByPredicate(OWLInverseOf)

	for _, inv := range inverseTriples {
		p1 := inv.Subject
		p2 := inv.Object

		// For X P1 Y, infer Y P2 X
		for _, t := range store.FindByPredicate(p1) {
			newTriple := Triple{Subject: t.Object, Predicate: p2, Object: t.Subject}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}

		// For X P2 Y, infer Y P1 X
		for _, t := range store.FindByPredicate(p2) {
			newTriple := Triple{Subject: t.Object, Predicate: p1, Object: t.Subject}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// TransitivePropertyInference implements owl:TransitiveProperty
// If P is transitive and X P Y and Y P Z, then X P Z
type TransitivePropertyInference struct{}

func (r *TransitivePropertyInference) Name() string {
	return "owl:TransitiveProperty-inference"
}

func (r *TransitivePropertyInference) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	// Find all transitive properties
	transitiveProps := make(map[string]bool)
	for _, t := range store.FindByPredicateObject(RDFType, OWLTransitiveProperty) {
		transitiveProps[t.Subject] = true
	}

	for prop := range transitiveProps {
		propTriples := store.FindByPredicate(prop)

		for _, t1 := range propTriples {
			x := t1.Subject
			y := t1.Object

			for _, t2 := range store.FindBySubjectPredicate(y, prop) {
				z := t2.Object
				newTriple := Triple{Subject: x, Predicate: prop, Object: z}
				if !store.Contains(newTriple) && x != z {
					inferred = append(inferred, newTriple)
				}
			}
		}
	}

	return inferred
}

// SymmetricPropertyInference implements owl:SymmetricProperty
// If P is symmetric and X P Y, then Y P X
type SymmetricPropertyInference struct{}

func (r *SymmetricPropertyInference) Name() string {
	return "owl:SymmetricProperty-inference"
}

func (r *SymmetricPropertyInference) Apply(store *TripleStore) []Triple {
	var inferred []Triple

	// Find all symmetric properties
	symmetricProps := make(map[string]bool)
	for _, t := range store.FindByPredicateObject(RDFType, OWLSymmetricProperty) {
		symmetricProps[t.Subject] = true
	}

	for prop := range symmetricProps {
		for _, t := range store.FindByPredicate(prop) {
			newTriple := Triple{Subject: t.Object, Predicate: prop, Object: t.Subject}
			if !store.Contains(newTriple) {
				inferred = append(inferred, newTriple)
			}
		}
	}

	return inferred
}

// DefaultRules returns the default set of reasoning rules
func DefaultRules() []Rule {
	return []Rule{
		&SubClassTransitivity{},
		&TypeInheritance{},
		&DomainInference{},
		&RangeInference{},
		&SubPropertyTransitivity{},
		&SubPropertyInheritance{},
		&EquivalentClassSymmetry{},
		&EquivalentClassTransitivity{},
		&SameAsSymmetry{},
		&SameAsTransitivity{},
		&InversePropertyInference{},
		&TransitivePropertyInference{},
		&SymmetricPropertyInference{},
	}
}