// Package reasoner provides forward reasoning capabilities for RDF/OWL ontologies.
// It parses Turtle format inputs and applies RDFS/OWL inference rules to derive
// new facts from the given TBox (terminology) and ABox (assertions).
package reasoner

import (
	"fmt"
	"sort"
)

// Reasoner performs forward reasoning on RDF data
type Reasoner struct {
	store  *TripleStore
	rules  []Rule
	parser *TurtleParser
}

// NewReasoner creates a new reasoner with default rules
func NewReasoner() *Reasoner {
	return &Reasoner{
		store:  NewTripleStore(),
		rules:  DefaultRules(),
		parser: NewTurtleParser(),
	}
}

// NewReasonerWithRules creates a new reasoner with custom rules
func NewReasonerWithRules(rules []Rule) *Reasoner {
	return &Reasoner{
		store:  NewTripleStore(),
		rules:  rules,
		parser: NewTurtleParser(),
	}
}

// LoadTurtle parses and loads Turtle content into the store
func (r *Reasoner) LoadTurtle(content string) error {
	triples, err := r.parser.Parse(content)
	if err != nil {
		return fmt.Errorf("failed to parse Turtle: %w", err)
	}

	for _, t := range triples {
		r.store.Add(t)
	}

	return nil
}

// RunForwardReasoning applies all rules until no new facts are derived
// Returns the number of new triples inferred
func (r *Reasoner) RunForwardReasoning() int {
	totalInferred := 0

	for {
		newInThisRound := 0

		for _, rule := range r.rules {
			inferred := rule.Apply(r.store)
			for _, t := range inferred {
				if r.store.Add(t) {
					newInThisRound++
				}
			}
		}

		if newInThisRound == 0 {
			break
		}

		totalInferred += newInThisRound
	}

	return totalInferred
}

// GetAllTriples returns all triples in the store as strings
func (r *Reasoner) GetAllTriples() []string {
	triples := r.store.All()
	result := make([]string, len(triples))
	for i, t := range triples {
		result[i] = t.String()
	}
	sort.Strings(result)
	return result
}

// GetInferredTypes returns all rdf:type assertions for a given subject
func (r *Reasoner) GetInferredTypes(subject string) []string {
	var types []string
	for _, t := range r.store.FindBySubjectPredicate(subject, RDFType) {
		types = append(types, t.Object)
	}
	sort.Strings(types)
	return types
}

// Query returns all triples matching the given pattern
// Use empty string "" as wildcard
func (r *Reasoner) Query(subject, predicate, object string) []Triple {
	var results []Triple

	if subject != "" && predicate != "" {
		for _, t := range r.store.FindBySubjectPredicate(subject, predicate) {
			if object == "" || t.Object == object {
				results = append(results, t)
			}
		}
	} else if subject != "" {
		for _, t := range r.store.FindBySubject(subject) {
			if (predicate == "" || t.Predicate == predicate) &&
				(object == "" || t.Object == object) {
				results = append(results, t)
			}
		}
	} else if predicate != "" {
		for _, t := range r.store.FindByPredicate(predicate) {
			if object == "" || t.Object == object {
				results = append(results, t)
			}
		}
	} else if object != "" {
		for _, t := range r.store.FindByObject(object) {
			results = append(results, t)
		}
	} else {
		results = r.store.All()
	}

	return results
}

// GetStore returns the underlying triple store
func (r *Reasoner) GetStore() *TripleStore {
	return r.store
}

// ForwardReason is the main public API function.
// It accepts TBox (terminology/schema) and ABox (assertions/instances) in Turtle format,
// performs forward reasoning, and returns all inferred triples as strings.
//
// Parameters:
//   - abox: Turtle string containing instance data (assertions)
//   - tbox: Turtle string containing schema/ontology definitions
//
// Returns:
//   - []string: List of all triples (including inferred) in N-Triples format
//   - error: Any parsing or processing errors
func ForwardReason(abox, tbox string) ([]string, error) {
	reasoner := NewReasoner()

	// Load TBox first (schema/ontology)
	if tbox != "" {
		if err := reasoner.LoadTurtle(tbox); err != nil {
			return nil, fmt.Errorf("failed to load TBox: %w", err)
		}
	}

	// Load ABox (instances/assertions)
	if abox != "" {
		if err := reasoner.LoadTurtle(abox); err != nil {
			return nil, fmt.Errorf("failed to load ABox: %w", err)
		}
	}

	// Run forward reasoning
	reasoner.RunForwardReasoning()

	return reasoner.GetAllTriples(), nil
}

// ForwardReasonWithDetails returns both original and inferred triples separately
func ForwardReasonWithDetails(abox, tbox string) (*ReasoningResult, error) {
	reasoner := NewReasoner()

	// Load TBox first
	if tbox != "" {
		if err := reasoner.LoadTurtle(tbox); err != nil {
			return nil, fmt.Errorf("failed to load TBox: %w", err)
		}
	}

	// Load ABox
	if abox != "" {
		if err := reasoner.LoadTurtle(abox); err != nil {
			return nil, fmt.Errorf("failed to load ABox: %w", err)
		}
	}

	// Get original triples count
	originalCount := reasoner.store.Size()
	originalTriples := reasoner.GetAllTriples()

	// Run forward reasoning
	inferredCount := reasoner.RunForwardReasoning()

	// Get all triples after reasoning
	allTriples := reasoner.GetAllTriples()

	// Separate inferred triples
	originalSet := make(map[string]bool)
	for _, t := range originalTriples {
		originalSet[t] = true
	}

	var inferredTriples []string
	for _, t := range allTriples {
		if !originalSet[t] {
			inferredTriples = append(inferredTriples, t)
		}
	}

	return &ReasoningResult{
		OriginalTriples: originalTriples,
		InferredTriples: inferredTriples,
		AllTriples:      allTriples,
		OriginalCount:   originalCount,
		InferredCount:   inferredCount,
		TotalCount:      len(allTriples),
	}, nil
}

// ReasoningResult contains detailed results from forward reasoning
type ReasoningResult struct {
	OriginalTriples []string // Triples from input
	InferredTriples []string // Newly inferred triples
	AllTriples      []string // All triples combined
	OriginalCount   int      // Number of original triples
	InferredCount   int      // Number of inferred triples
	TotalCount      int      // Total number of triples
}
