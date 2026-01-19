package reasoner

import (
	"fmt"
	"strings"
)

// ConvertTriplesToDatalog converts a list of N-Triple strings to Datalog format
// Each RDF triple (subject, predicate, object) becomes a Datalog fact: predicate(subject, object)
// IRIs are converted to simplified names by extracting the local part after # or /
func ConvertTriplesToDatalog(triples []string) []string {
	datalogFacts := make([]string, 0, len(triples))

	for _, triple := range triples {
		// Parse the N-Triple format: <subject> <predicate> <object> .
		triple = strings.TrimSpace(triple)
		if !strings.HasSuffix(triple, " .") {
			continue // Skip malformed triples
		}

		// Remove the trailing " ."
		triple = strings.TrimSuffix(triple, " .")

		// Split into parts, handling quoted literals
		parts := parseNTripleParts(triple)
		if len(parts) != 3 {
			continue // Skip malformed triples
		}

		subject := simplifyIRI(parts[0])
		predicate := simplifyIRI(parts[1])
		object := simplifyIRI(parts[2])

		// Format as Datalog fact: predicate(subject, object)
		datalogFact := fmt.Sprintf("%s(%s, %s).", predicate, subject, object)
		datalogFacts = append(datalogFacts, datalogFact)
	}

	return datalogFacts
}

// String returns a human-readable summary of the reasoning result
func (r *ReasoningResult) String() string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("Reasoning Result:\n"))
	sb.WriteString(fmt.Sprintf("  Original triples: %d\n", r.OriginalCount))
	sb.WriteString(fmt.Sprintf("  Inferred triples: %d\n", r.InferredCount))
	sb.WriteString(fmt.Sprintf("  Total triples: %d\n", r.TotalCount))
	return sb.String()
}

// parseNTripleParts splits an N-Triple into subject, predicate, and object parts
func parseNTripleParts(triple string) []string {
	var parts []string
	var current strings.Builder
	inQuotes := false
	escaped := false

	for i, r := range triple {
		switch {
		case escaped:
			current.WriteRune(r)
			escaped = false
		case r == '\\':
			current.WriteRune(r)
			escaped = true
		case r == '"':
			current.WriteRune(r)
			inQuotes = !inQuotes
		case r == ' ' && !inQuotes:
			if current.Len() > 0 {
				parts = append(parts, strings.TrimSpace(current.String()))
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}

		// Handle the last part
		if i == len(triple)-1 && current.Len() > 0 {
			parts = append(parts, strings.TrimSpace(current.String()))
		}
	}

	return parts
}

// simplifyIRI converts an IRI to a simplified name for Datalog
func simplifyIRI(iri string) string {
	// Remove angle brackets for IRIs
	if strings.HasPrefix(iri, "<") && strings.HasSuffix(iri, ">") {
		iri = iri[1 : len(iri)-1]
	}

	// Handle quoted literals
	if strings.HasPrefix(iri, "\"") {
		// For literals, extract the value and make it safe for Datalog
		if strings.Contains(iri, "\"^^") {
			// Typed literal: "value"^^<type>
			parts := strings.Split(iri, "\"^^")
			if len(parts) >= 2 {
				value := strings.Trim(parts[0], "\"")
				return makeDatalogSafe(value)
			}
		} else if strings.Contains(iri, "\"@") {
			// Language tagged literal: "value"@lang
			parts := strings.Split(iri, "\"@")
			if len(parts) >= 2 {
				value := strings.Trim(parts[0], "\"")
				return makeDatalogSafe(value)
			}
		} else {
			// Simple literal: "value"
			value := strings.Trim(iri, "\"")
			return makeDatalogSafe(value)
		}
	}

	// Extract local name from IRI
	if strings.Contains(iri, "#") {
		parts := strings.Split(iri, "#")
		if len(parts) > 1 {
			return makeDatalogSafe(parts[len(parts)-1])
		}
	}

	if strings.Contains(iri, "/") {
		parts := strings.Split(iri, "/")
		if len(parts) > 1 {
			return makeDatalogSafe(parts[len(parts)-1])
		}
	}

	return makeDatalogSafe(iri)
}

// makeDatalogSafe converts a string to be safe for use in Datalog
func makeDatalogSafe(s string) string {
	// Replace common characters that might cause issues in Datalog
	s = strings.ReplaceAll(s, "-", "_")
	s = strings.ReplaceAll(s, ":", "_")
	s = strings.ReplaceAll(s, " ", "_")
	s = strings.ReplaceAll(s, ".", "_")

	// If it starts with a number, prefix with underscore
	if len(s) > 0 && s[0] >= '0' && s[0] <= '9' {
		s = "_" + s
	}

	// If empty or problematic, use a default
	if s == "" || s == "_" {
		s = "unknown"
	}

	return s
}
