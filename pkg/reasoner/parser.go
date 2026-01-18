package reasoner

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// TurtleParser parses Turtle format RDF
type TurtleParser struct {
	prefixes map[string]string
	base     string
	input    string
	pos      int
}

// NewTurtleParser creates a new Turtle parser
func NewTurtleParser() *TurtleParser {
	return &TurtleParser{
		prefixes: make(map[string]string),
		base:     "",
	}
}

// Parse parses Turtle content and returns triples
func (p *TurtleParser) Parse(content string) ([]Triple, error) {
	// Reset parser state
	p.prefixes = make(map[string]string)
	p.base = ""
	p.input = content
	p.pos = 0

	var triples []Triple

	// Preprocess: remove BOM, normalize line endings
	p.input = strings.TrimPrefix(p.input, "\ufeff")
	p.input = strings.ReplaceAll(p.input, "\r\n", "\n")

	for p.pos < len(p.input) {
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			break
		}

		// Check for prefix declaration
		if p.lookingAt("@prefix") || p.lookingAtCaseInsensitive("PREFIX") {
			if err := p.parsePrefix(); err != nil {
				return nil, err
			}
			continue
		}

		// Check for base declaration
		if p.lookingAt("@base") || p.lookingAtCaseInsensitive("BASE") {
			if err := p.parseBase(); err != nil {
				return nil, err
			}
			continue
		}

		// Parse triple(s)
		newTriples, err := p.parseTriples()
		if err != nil {
			// Try to skip to next statement on error
			p.skipToNextStatement()
			continue
		}
		triples = append(triples, newTriples...)
	}

	return triples, nil
}

func (p *TurtleParser) skipWhitespaceAndComments() {
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
			p.pos++
		} else if ch == '#' {
			// Skip comment until end of line
			for p.pos < len(p.input) && p.input[p.pos] != '\n' {
				p.pos++
			}
		} else {
			break
		}
	}
}

func (p *TurtleParser) lookingAt(s string) bool {
	if p.pos+len(s) > len(p.input) {
		return false
	}
	return p.input[p.pos:p.pos+len(s)] == s
}

func (p *TurtleParser) lookingAtCaseInsensitive(s string) bool {
	if p.pos+len(s) > len(p.input) {
		return false
	}
	return strings.EqualFold(p.input[p.pos:p.pos+len(s)], s)
}

func (p *TurtleParser) parsePrefix() error {
	// Skip @prefix or PREFIX
	if p.lookingAt("@prefix") {
		p.pos += 7
	} else {
		p.pos += 6
	}

	p.skipWhitespaceAndComments()

	// Parse prefix name (e.g., "ex:")
	prefixName := p.parsePrefixName()
	if prefixName == "" {
		return fmt.Errorf("expected prefix name at position %d", p.pos)
	}

	p.skipWhitespaceAndComments()

	// Parse IRI
	iri, err := p.parseIRI()
	if err != nil {
		return err
	}

	p.skipWhitespaceAndComments()

	// Skip optional '.'
	if p.pos < len(p.input) && p.input[p.pos] == '.' {
		p.pos++
	}

	// Store prefix (without the colon)
	prefix := strings.TrimSuffix(prefixName, ":")
	p.prefixes[prefix] = iri

	return nil
}

func (p *TurtleParser) parseBase() error {
	// Skip @base or BASE
	if p.lookingAt("@base") {
		p.pos += 5
	} else {
		p.pos += 4
	}

	p.skipWhitespaceAndComments()

	// Parse IRI
	iri, err := p.parseIRI()
	if err != nil {
		return err
	}

	p.skipWhitespaceAndComments()

	// Skip optional '.'
	if p.pos < len(p.input) && p.input[p.pos] == '.' {
		p.pos++
	}

	p.base = iri
	return nil
}

func (p *TurtleParser) parsePrefixName() string {
	start := p.pos

	// Read until ':'
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ':' {
			p.pos++
			return p.input[start:p.pos]
		}
		if unicode.IsSpace(rune(ch)) {
			break
		}
		p.pos++
	}

	return p.input[start:p.pos]
}

func (p *TurtleParser) parseIRI() (string, error) {
	p.skipWhitespaceAndComments()

	if p.pos >= len(p.input) || p.input[p.pos] != '<' {
		return "", fmt.Errorf("expected '<' at position %d", p.pos)
	}

	p.pos++ // skip '<'
	start := p.pos

	for p.pos < len(p.input) && p.input[p.pos] != '>' {
		p.pos++
	}

	if p.pos >= len(p.input) {
		return "", fmt.Errorf("unterminated IRI")
	}

	iri := p.input[start:p.pos]
	p.pos++ // skip '>'

	return iri, nil
}

func (p *TurtleParser) parseTriples() ([]Triple, error) {
	var triples []Triple

	// Parse subject
	subject, err := p.parseSubject()
	if err != nil {
		return nil, err
	}

	// Parse predicate-object list
	for {
		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			break
		}

		// Check for end of statement
		if p.input[p.pos] == '.' {
			p.pos++
			break
		}

		// Parse predicate
		predicate, err := p.parsePredicate()
		if err != nil {
			return nil, err
		}

		// Parse object list
		for {
			p.skipWhitespaceAndComments()

			object, err := p.parseObject()
			if err != nil {
				return nil, err
			}

			triples = append(triples, Triple{
				Subject:   subject,
				Predicate: predicate,
				Object:    object,
			})

			p.skipWhitespaceAndComments()
			if p.pos >= len(p.input) {
				break
			}

			// Check for object list continuation
			if p.input[p.pos] == ',' {
				p.pos++
				continue
			}
			break
		}

		p.skipWhitespaceAndComments()
		if p.pos >= len(p.input) {
			break
		}

		// Check for predicate-object list continuation
		if p.input[p.pos] == ';' {
			p.pos++
			p.skipWhitespaceAndComments()
			// Check if followed by '.' (empty predicate-object after semicolon)
			if p.pos < len(p.input) && p.input[p.pos] == '.' {
				p.pos++
				break
			}
			continue
		}

		if p.input[p.pos] == '.' {
			p.pos++
			break
		}

		break
	}

	return triples, nil
}

func (p *TurtleParser) parseSubject() (string, error) {
	p.skipWhitespaceAndComments()

	if p.pos >= len(p.input) {
		return "", fmt.Errorf("unexpected end of input")
	}

	// IRI
	if p.input[p.pos] == '<' {
		iri, err := p.parseIRI()
		if err != nil {
			return "", err
		}
		return p.resolveIRI(iri), nil
	}

	// Blank node
	if p.lookingAt("_:") {
		return p.parseBlankNode()
	}

	// Prefixed name
	return p.parsePrefixedName()
}

func (p *TurtleParser) parsePredicate() (string, error) {
	p.skipWhitespaceAndComments()

	if p.pos >= len(p.input) {
		return "", fmt.Errorf("unexpected end of input")
	}

	// 'a' keyword for rdf:type
	if p.pos+1 <= len(p.input) && p.input[p.pos] == 'a' {
		// Check it's standalone 'a' not part of another token
		if p.pos+1 >= len(p.input) || !isNameChar(rune(p.input[p.pos+1])) {
			p.pos++
			return RDFType, nil
		}
	}

	// IRI
	if p.input[p.pos] == '<' {
		iri, err := p.parseIRI()
		if err != nil {
			return "", err
		}
		return p.resolveIRI(iri), nil
	}

	// Prefixed name
	return p.parsePrefixedName()
}

func (p *TurtleParser) parseObject() (string, error) {
	p.skipWhitespaceAndComments()

	if p.pos >= len(p.input) {
		return "", fmt.Errorf("unexpected end of input")
	}

	// IRI
	if p.input[p.pos] == '<' {
		iri, err := p.parseIRI()
		if err != nil {
			return "", err
		}
		return p.resolveIRI(iri), nil
	}

	// Blank node
	if p.lookingAt("_:") {
		return p.parseBlankNode()
	}

	// Literal
	if p.input[p.pos] == '"' {
		return p.parseLiteral()
	}

	// Prefixed name
	return p.parsePrefixedName()
}

func (p *TurtleParser) parseBlankNode() (string, error) {
	start := p.pos
	p.pos += 2 // skip "_:"

	for p.pos < len(p.input) && isNameChar(rune(p.input[p.pos])) {
		p.pos++
	}

	return p.input[start:p.pos], nil
}

func (p *TurtleParser) parsePrefixedName() (string, error) {
	start := p.pos

	// Read prefix part
	for p.pos < len(p.input) {
		ch := p.input[p.pos]
		if ch == ':' {
			break
		}
		if !isNameChar(rune(ch)) {
			break
		}
		p.pos++
	}

	if p.pos >= len(p.input) || p.input[p.pos] != ':' {
		return "", fmt.Errorf("expected ':' in prefixed name at position %d", p.pos)
	}

	prefix := p.input[start:p.pos]
	p.pos++ // skip ':'

	// Read local part
	localStart := p.pos
	for p.pos < len(p.input) && isNameChar(rune(p.input[p.pos])) {
		p.pos++
	}

	local := p.input[localStart:p.pos]

	// Resolve prefix
	if base, ok := p.prefixes[prefix]; ok {
		return base + local, nil
	}

	return prefix + ":" + local, nil
}

func (p *TurtleParser) parseLiteral() (string, error) {
	var sb strings.Builder

	// Check for triple-quoted string
	if p.lookingAt(`"""`) {
		p.pos += 3
		sb.WriteString(`"`)

		for p.pos < len(p.input) {
			if p.lookingAt(`"""`) {
				p.pos += 3
				sb.WriteString(`"`)
				break
			}
			sb.WriteByte(p.input[p.pos])
			p.pos++
		}
	} else {
		// Single-quoted string
		p.pos++ // skip opening quote
		sb.WriteString(`"`)

		for p.pos < len(p.input) && p.input[p.pos] != '"' {
			if p.input[p.pos] == '\\' && p.pos+1 < len(p.input) {
				sb.WriteByte(p.input[p.pos])
				p.pos++
				sb.WriteByte(p.input[p.pos])
				p.pos++
				continue
			}
			sb.WriteByte(p.input[p.pos])
			p.pos++
		}

		if p.pos < len(p.input) {
			p.pos++ // skip closing quote
		}
		sb.WriteString(`"`)
	}

	// Check for language tag or datatype
	if p.pos < len(p.input) && p.input[p.pos] == '@' {
		// Language tag
		p.pos++
		tagStart := p.pos
		for p.pos < len(p.input) && (isAlphaNum(rune(p.input[p.pos])) || p.input[p.pos] == '-') {
			p.pos++
		}
		sb.WriteString("@")
		sb.WriteString(p.input[tagStart:p.pos])
	} else if p.lookingAt("^^") {
		// Datatype
		p.pos += 2
		sb.WriteString("^^")

		if p.pos < len(p.input) && p.input[p.pos] == '<' {
			iri, _ := p.parseIRI()
			sb.WriteString("<")
			sb.WriteString(p.resolveIRI(iri))
			sb.WriteString(">")
		} else {
			// Prefixed datatype
			dt, _ := p.parsePrefixedName()
			sb.WriteString("<")
			sb.WriteString(dt)
			sb.WriteString(">")
		}
	}

	return sb.String(), nil
}

func (p *TurtleParser) resolveIRI(iri string) string {
	if p.base != "" && !strings.Contains(iri, "://") && !strings.HasPrefix(iri, "#") {
		return p.base + iri
	}
	return iri
}

func (p *TurtleParser) skipToNextStatement() {
	for p.pos < len(p.input) && p.input[p.pos] != '.' {
		p.pos++
	}
	if p.pos < len(p.input) {
		p.pos++
	}
}

func isNameChar(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '-' || r == '.'
}

func isAlphaNum(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r)
}

// Alternative fallback parser using regex (for complex cases)
func (p *TurtleParser) FallbackParse(content string) ([]Triple, error) {
	var triples []Triple

	// Parse prefixes
	prefixRe := regexp.MustCompile(`@prefix\s+([a-zA-Z_][\w-]*):\s*<([^>]+)>\s*\.`)
	for _, match := range prefixRe.FindAllStringSubmatch(content, -1) {
		p.prefixes[match[1]] = match[2]
	}

	// Parse base
	baseRe := regexp.MustCompile(`@base\s*<([^>]+)>\s*\.`)
	if match := baseRe.FindStringSubmatch(content); match != nil {
		p.base = match[1]
	}

	// Remove declarations
	content = prefixRe.ReplaceAllString(content, "")
	content = baseRe.ReplaceAllString(content, "")

	// Remove comments
	commentRe := regexp.MustCompile(`#[^\n]*`)
	content = commentRe.ReplaceAllString(content, "")

	// Split into statements
	statements := p.splitStatements(content)

	for _, stmt := range statements {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}

		ts := p.parseStatement(stmt)
		triples = append(triples, ts...)
	}

	return triples, nil
}

func (p *TurtleParser) splitStatements(content string) []string {
	var statements []string
	var current strings.Builder
	inString := false
	inTripleString := false
	inIRI := false

	runes := []rune(content)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Check for triple-quoted string
		if i+2 < len(runes) && runes[i] == '"' && runes[i+1] == '"' && runes[i+2] == '"' {
			if !inTripleString && !inString {
				inTripleString = true
			} else if inTripleString {
				inTripleString = false
			}
			current.WriteRune(r)
			continue
		}

		if r == '<' && !inString && !inTripleString {
			inIRI = true
		} else if r == '>' && inIRI {
			inIRI = false
		} else if r == '"' && !inTripleString {
			inString = !inString
		}

		if r == '.' && !inString && !inTripleString && !inIRI {
			statements = append(statements, current.String())
			current.Reset()
		} else {
			current.WriteRune(r)
		}
	}

	if current.Len() > 0 {
		statements = append(statements, current.String())
	}

	return statements
}

func (p *TurtleParser) parseStatement(stmt string) []Triple {
	var triples []Triple

	tokens := p.tokenize(stmt)
	if len(tokens) < 3 {
		return triples
	}

	subject := p.resolveToken(tokens[0])

	i := 1
	for i < len(tokens) {
		if i >= len(tokens) {
			break
		}

		predicate := p.resolveToken(tokens[i])
		i++

		for i < len(tokens) {
			token := tokens[i]
			if token == ";" {
				i++
				break
			}
			if token == "," {
				i++
				continue
			}

			object := p.resolveToken(token)
			triples = append(triples, Triple{
				Subject:   subject,
				Predicate: predicate,
				Object:    object,
			})
			i++
		}
	}

	return triples
}

func (p *TurtleParser) tokenize(stmt string) []string {
	var tokens []string
	var current strings.Builder
	inString := false
	inTripleString := false
	inIRI := false

	runes := []rune(stmt)
	for i := 0; i < len(runes); i++ {
		r := runes[i]

		// Check for triple-quoted string
		if i+2 < len(runes) && runes[i] == '"' && runes[i+1] == '"' && runes[i+2] == '"' {
			if !inTripleString && !inString {
				inTripleString = true
				current.WriteString(`"""`)
				i += 2
				continue
			}
		}
		if inTripleString && i+2 < len(runes) && runes[i] == '"' && runes[i+1] == '"' && runes[i+2] == '"' {
			inTripleString = false
			current.WriteString(`"""`)
			i += 2
			tokens = append(tokens, current.String())
			current.Reset()
			continue
		}

		if r == '<' && !inString && !inTripleString {
			inIRI = true
			current.WriteRune(r)
			continue
		}
		if r == '>' && inIRI {
			inIRI = false
			current.WriteRune(r)
			tokens = append(tokens, current.String())
			current.Reset()
			continue
		}
		if r == '"' && !inTripleString {
			inString = !inString
			current.WriteRune(r)
			if !inString {
				// Check for suffix
				for i+1 < len(runes) && (runes[i+1] == '@' || runes[i+1] == '^') {
					i++
					current.WriteRune(runes[i])
					if runes[i] == '^' && i+1 < len(runes) && runes[i+1] == '^' {
						i++
						current.WriteRune(runes[i])
						for i+1 < len(runes) && !unicode.IsSpace(runes[i+1]) && runes[i+1] != ';' && runes[i+1] != ',' {
							i++
							current.WriteRune(runes[i])
							if runes[i] == '>' {
								break
							}
						}
					} else if runes[i] == '@' {
						for i+1 < len(runes) && (isAlphaNum(runes[i+1]) || runes[i+1] == '-') {
							i++
							current.WriteRune(runes[i])
						}
					}
				}
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		if inIRI || inString || inTripleString {
			current.WriteRune(r)
			continue
		}

		if unicode.IsSpace(r) {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			continue
		}

		if r == ';' || r == ',' {
			if current.Len() > 0 {
				tokens = append(tokens, current.String())
				current.Reset()
			}
			tokens = append(tokens, string(r))
			continue
		}

		current.WriteRune(r)
	}

	if current.Len() > 0 {
		tokens = append(tokens, current.String())
	}

	return tokens
}

func (p *TurtleParser) resolveToken(token string) string {
	token = strings.TrimSpace(token)

	if strings.HasPrefix(token, "<") && strings.HasSuffix(token, ">") {
		return p.resolveIRI(strings.Trim(token, "<>"))
	}

	if strings.HasPrefix(token, `"`) {
		return token
	}

	if token == "a" {
		return RDFType
	}

	if strings.HasPrefix(token, "_:") {
		return token
	}

	if strings.Contains(token, ":") {
		parts := strings.SplitN(token, ":", 2)
		if len(parts) == 2 {
			if base, ok := p.prefixes[parts[0]]; ok {
				return base + parts[1]
			}
		}
	}

	return token
}