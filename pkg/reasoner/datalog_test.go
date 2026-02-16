package reasoner

import (
	"testing"
)

func TestDatalog(t *testing.T) {
	datalogContent := `
Disjoint(Gemeinde,Kanton).
Type(Zürich_Stadt,Gemeinde).
Type(Kanton_Zürich,Kanton).
Type(Basel_Stadt,Gemeinde).
Type(Kanton_Basel,Kanton).
Disjoint(X,Y).
Type(X,Y).
Is_A_Kanton(X).
Not_Kanton(X).
Is_A_Kanton(X). :- Type(X,Kanton).
Not_Kanton(X). :- Disjoint(Gemeinde,Kanton), Type(X,Gemeinde).
`

	tests := []struct {
		query    string
		expected bool
	}{
		{"?‑ Type(Kanton_Zürich, Gemeinde).", false},
		{"?‑ Type(Kanton_Zürich, Kanton).", true},
		{"?‑ Is_A_Kanton(Kanton_Zürich).", true},
		{"?‑ Not_Kanton(Zürich_Stadt).", true},
		{"?‑ Not_Kanton(Kanton_Zürich).", false},
	}

	for _, tt := range tests {
		result, err := DLQuery(datalogContent, tt.query)
		if err != nil {
			t.Errorf("DLQuery error for %s: %v", tt.query, err)
			continue
		}
		if result != tt.expected {
			t.Errorf("DLQuery(%s) = %v, expected %v", tt.query, result, tt.expected)
		}
	}
}

func TestDatalogMultiCharVar(t *testing.T) {
	datalogContent := `
Parent(john, mary).
Ancestor(X, Y) :- Parent(X, Y).
Ancestor(X, Z) :- Parent(X, Y), Ancestor(Y, Z).
`
	result, _ := DLQuery(datalogContent, "?- Ancestor(john, mary).")
	if !result {
		t.Errorf("Ancestor(john, mary) should be true")
	}

	datalogContent2 := `
Parent(john, mary).
Parent(mary, jane).
Ancestor(VAR_X, VAR_Y) :- Parent(VAR_X, VAR_Y).
Ancestor(VAR_X, VAR_Z) :- Parent(VAR_X, VAR_Y), Ancestor(VAR_Y, VAR_Z).
`
	result2, _ := DLQuery(datalogContent2, "?- Ancestor(john, jane).")
	if !result2 {
		t.Errorf("Ancestor(john, jane) should be true with multi-char variables")
	}
}

func TestParser(t *testing.T) {
	input := "Parent(john, mary). Human(X) :- Parent(X, Y)."
	program, err := ParseDatalog(input)
	if err != nil {
		t.Fatalf("ParseDatalog failed: %v", err)
	}

	if len(program.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(program.Facts))
	}
	if len(program.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(program.Rules))
	}
}

func TestParserWithComments(t *testing.T) {
	input := `
% This is a comment.
Parent(john, mary). % Another comment.
Human(X) :- Parent(X, Y). // C++ style comment.
`
	program, err := ParseDatalog(input)
	if err != nil {
		t.Fatalf("ParseDatalog failed: %v", err)
	}

	if len(program.Facts) != 1 {
		t.Errorf("Expected 1 fact, got %d", len(program.Facts))
	}
	if len(program.Rules) != 1 {
		t.Errorf("Expected 1 rule, got %d", len(program.Rules))
	}
}
