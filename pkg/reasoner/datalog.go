package reasoner

import (
	"fmt"
	"strings"
	"unicode"
)

// DLTerm represents a Datalog term (either a constant or a variable)
type DLTerm struct {
	Value      string
	IsVariable bool
}

// DLAtom represents a Datalog atom: Predicate(DLTerm1, DLTerm2, ...)
type DLAtom struct {
	Predicate string
	Terms     []DLTerm
}

// DLRule represents a Datalog rule: Head :- Body1, Body2, ...
type DLRule struct {
	Head DLAtom
	Body []DLAtom
}

// DatalogProgram represents a collection of facts and rules
type DatalogProgram struct {
	Facts []DLAtom
	Rules []DLRule
}

func (a DLAtom) String() string {
	var terms []string
	for _, t := range a.Terms {
		terms = append(terms, t.Value)
	}
	return fmt.Sprintf("%s(%s)", a.Predicate, strings.Join(terms, ", "))
}

// ParseDatalog parses a Datalog program from a string
func ParseDatalog(input string) (*DatalogProgram, error) {
	program := &DatalogProgram{}
	statements := splitDatalogStatements(input)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		if strings.Contains(stmt, ":-") {
			// It's a rule
			rule, err := parseRule(stmt)
			if err != nil {
				return nil, err
			}
			program.Rules = append(program.Rules, rule)
		} else {
			// It's a fact (or declaration)
			atom, err := parseAtom(stmt)
			if err != nil {
				return nil, err
			}
			program.Facts = append(program.Facts, atom)
		}
	}

	return program, nil
}

func splitDatalogStatements(input string) []string {
	var statements []string
	var current strings.Builder
	parenCount := 0

	lines := strings.Split(input, "\n")
	for _, line := range lines {
		// Remove comments
		if idx := strings.Index(line, "%"); idx != -1 {
			line = line[:idx]
		}
		if idx := strings.Index(line, "//"); idx != -1 {
			line = line[:idx]
		}
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		runes := []rune(line)
		for i := 0; i < len(runes); i++ {
			r := runes[i]
			if r == '(' {
				parenCount++
			} else if r == ')' {
				parenCount--
			}

			// A dot splits statements ONLY if it's not inside parentheses
			// AND it's not immediately followed by ":-"
			if r == '.' && parenCount == 0 {
				isFollowedByRule := false
				for j := i + 1; j < len(runes); j++ {
					if unicode.IsSpace(runes[j]) {
						continue
					}
					if j+1 < len(runes) && runes[j] == ':' && runes[j+1] == '-' {
						isFollowedByRule = true
					}
					break
				}

				if !isFollowedByRule {
					statements = append(statements, current.String())
					current.Reset()
					continue
				}
			}
			current.WriteRune(r)
		}
		current.WriteRune(' ') // Space between lines
	}

	if strings.TrimSpace(current.String()) != "" {
		statements = append(statements, current.String())
	}

	return statements
}

func parseRule(line string) (DLRule, error) {
	parts := strings.Split(line, ":-")
	if len(parts) != 2 {
		return DLRule{}, fmt.Errorf("invalid rule format: %s", line)
	}

	headStr := strings.TrimSpace(parts[0])
	headStr = strings.TrimSuffix(headStr, ".")
	head, err := parseAtom(headStr)
	if err != nil {
		return DLRule{}, err
	}

	bodyStr := strings.TrimSpace(parts[1])
	bodyStr = strings.TrimSuffix(bodyStr, ".")
	bodyParts := splitAtoms(bodyStr)
	var body []DLAtom
	for _, bp := range bodyParts {
		atom, err := parseAtom(strings.TrimSpace(bp))
		if err != nil {
			return DLRule{}, err
		}
		body = append(body, atom)
	}

	return DLRule{Head: head, Body: body}, nil
}

func parseAtom(s string) (DLAtom, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, ".")
	openParen := strings.Index(s, "(")
	closeParen := strings.LastIndex(s, ")")

	if openParen == -1 {
		// Arity 0 atom or just a predicate name
		return DLAtom{Predicate: s}, nil
	}

	if closeParen == -1 || closeParen < openParen {
		return DLAtom{}, fmt.Errorf("invalid atom format: %s", s)
	}

	predicate := strings.TrimSpace(s[:openParen])
	termsStr := s[openParen+1 : closeParen]
	termParts := strings.Split(termsStr, ",")

	var terms []DLTerm
	for _, tp := range termParts {
		tp = strings.TrimSpace(tp)
		if tp == "" {
			continue
		}
		isVar := isVariable(tp)
		terms = append(terms, DLTerm{Value: tp, IsVariable: isVar})
	}

	return DLAtom{Predicate: predicate, Terms: terms}, nil
}

func isVariable(s string) bool {
	if len(s) == 0 {
		return false
	}
	// Support single-character uppercase variables (like X, Y) or multi-character starting with ?
	if len(s) == 1 && s[0] >= 'A' && s[0] <= 'Z' {
		return true
	}
	if strings.HasPrefix(s, "?") {
		return true
	}
	// Support multi-character variables if they are all uppercase
	allUpper := true
	hasLetter := false
	for _, r := range s {
		if unicode.IsLetter(r) {
			hasLetter = true
			if !unicode.IsUpper(r) {
				allUpper = false
				break
			}
		} else if !unicode.IsDigit(r) && r != '_' {
			allUpper = false
			break
		}
	}
	return hasLetter && allUpper
}

func splitAtoms(s string) []string {
	var atoms []string
	var current strings.Builder
	parenCount := 0

	for _, r := range s {
		if r == '(' {
			parenCount++
		} else if r == ')' {
			parenCount--
		}

		if r == ',' && parenCount == 0 {
			atoms = append(atoms, current.String())
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		atoms = append(atoms, current.String())
	}

	return atoms
}

// ParseQuery parses a Datalog query
func ParseQuery(s string) (DLAtom, error) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "?-")
	s = strings.TrimPrefix(s, "?‑") // Handle non-standard hyphen
	s = strings.TrimSpace(s)
	s = strings.TrimSuffix(s, ".")
	return parseAtom(s)
}

// Reason evaluates the Datalog program and returns all derived facts
func (p *DatalogProgram) Reason() []DLAtom {
	factMap := make(map[string]DLAtom)
	var factList []DLAtom

	addFact := func(f DLAtom) bool {
		s := f.String()
		if _, ok := factMap[s]; !ok {
			factMap[s] = f
			factList = append(factList, f)
			return true
		}
		return false
	}

	for _, f := range p.Facts {
		if !hasVariables(f) {
			addFact(f)
		}
	}

	for {
		newFactsCount := 0
		for _, rule := range p.Rules {
			substitutions := p.findSubstitutions(rule.Body, factList, make(map[string]string))
			for _, sub := range substitutions {
				head := applySubstitution(rule.Head, sub)
				if !hasVariables(head) {
					if addFact(head) {
						newFactsCount++
					}
				}
			}
		}

		if newFactsCount == 0 {
			break
		}
	}

	return factList
}

func hasVariables(a DLAtom) bool {
	for _, t := range a.Terms {
		if t.IsVariable {
			return true
		}
	}
	return false
}

// DLQuery is the main public API function for Datalog querying.
// It accepts a Datalog program and a query, performs reasoning,
// and returns true if the query is satisfied.
func DLQuery(datalogContent, queryStr string) (bool, error) {
	program, err := ParseDatalog(datalogContent)
	if err != nil {
		return false, fmt.Errorf("failed to parse Datalog: %w", err)
	}

	query, err := ParseQuery(queryStr)
	if err != nil {
		return false, fmt.Errorf("failed to parse query: %w", err)
	}

	derivedFacts := program.Reason()
	return program.EvaluateQuery(query, derivedFacts), nil
}

func (p *DatalogProgram) findSubstitutions(body []DLAtom, facts []DLAtom, currentSub map[string]string) []map[string]string {
	if len(body) == 0 {
		return []map[string]string{currentSub}
	}

	var results []map[string]string
	first := body[0]
	rest := body[1:]

	// Find all facts that match 'first' under 'currentSub'
	for _, f := range facts {
		if f.Predicate != first.Predicate || len(f.Terms) != len(first.Terms) {
			continue
		}

		newSub := make(map[string]string)
		for k, v := range currentSub {
			newSub[k] = v
		}

		match := true
		for i, t := range first.Terms {
			factTerm := f.Terms[i].Value
			if t.IsVariable {
				if val, ok := newSub[t.Value]; ok {
					if val != factTerm {
						match = false
						break
					}
				} else {
					newSub[t.Value] = factTerm
				}
			} else {
				if t.Value != factTerm {
					match = false
					break
				}
			}
		}

		if match {
			results = append(results, p.findSubstitutions(rest, facts, newSub)...)
		}
	}

	return results
}

func applySubstitution(a DLAtom, sub map[string]string) DLAtom {
	newTerms := make([]DLTerm, len(a.Terms))
	for i, t := range a.Terms {
		if t.IsVariable {
			if val, ok := sub[t.Value]; ok {
				newTerms[i] = DLTerm{Value: val, IsVariable: false}
			} else {
				newTerms[i] = t
			}
		} else {
			newTerms[i] = t
		}
	}
	return DLAtom{Predicate: a.Predicate, Terms: newTerms}
}

// EvaluateQuery checks if a query matches any derived facts
func (p *DatalogProgram) EvaluateQuery(query DLAtom, derivedFacts []DLAtom) bool {
	for _, f := range derivedFacts {
		if f.Predicate != query.Predicate || len(f.Terms) != len(query.Terms) {
			continue
		}

		match := true
		for i, qt := range query.Terms {
			ft := f.Terms[i]
			if !qt.IsVariable && qt.Value != ft.Value {
				match = false
				break
			}
		}
		if match {
			return true
		}
	}
	return false
}
