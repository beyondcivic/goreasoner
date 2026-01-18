# goreasoner

Minimalistic Forward Reasoner written in Golang.

# Go RDF Reasoner Library

A lightweight in-memory forward reasoning library for RDF/OWL ontologies in Go. Parses Turtle format and applies RDFS/OWL inference rules.

## Features

- **Turtle Parser**: Custom parser (no external dependencies) supporting:
  - Prefixes and base IRIs
  - Full IRIs and prefixed names
  - Blank nodes
  - Literals with language tags and datatypes
  - Predicate-object lists (`;`) and object lists (`,`)

- **Forward Reasoning Rules**:
  - `rdfs:subClassOf` transitivity
  - `rdf:type` inheritance through class hierarchy
  - `rdfs:domain` and `rdfs:range` inference
  - `rdfs:subPropertyOf` transitivity and inheritance
  - `owl:equivalentClass` symmetry and transitivity
  - `owl:sameAs` symmetry and transitivity
  - `owl:inverseOf` inference
  - `owl:TransitiveProperty` inference
  - `owl:SymmetricProperty` inference

## Installation

```bash
go get github.com/example/reasoner
```

## Usage

### Simple API

```go
package main

import (
    "fmt"
    "github.com/example/reasoner"
)

func main() {
    tbox := `
@prefix ex: <http://example.org/> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .

ex:Car rdfs:subClassOf ex:Vehicle .
ex:Vehicle rdfs:subClassOf ex:Transport .
`

    abox := `
@prefix ex: <http://example.org/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .

ex:myCar rdf:type ex:Car .
`

    // Get all triples (including inferred)
    triples, err := reasoner.ForwardReason(abox, tbox)
    if err != nil {
        panic(err)
    }

    for _, t := range triples {
        fmt.Println(t)
    }
}
```

### Detailed Results

```go
result, err := reasoner.ForwardReasonWithDetails(abox, tbox)
if err != nil {
    panic(err)
}

fmt.Printf("Original: %d, Inferred: %d, Total: %d\n",
    result.OriginalCount, result.InferredCount, result.TotalCount)

fmt.Println("Inferred triples:")
for _, t := range result.InferredTriples {
    fmt.Println(t)
}
```

### Using the Reasoner Directly

```go
r := reasoner.NewReasoner()

// Load TBox
if err := r.LoadTurtle(tbox); err != nil {
    panic(err)
}

// Load ABox
if err := r.LoadTurtle(abox); err != nil {
    panic(err)
}

// Run reasoning
inferredCount := r.RunForwardReasoning()
fmt.Printf("Inferred %d new triples\n", inferredCount)

// Query for specific patterns (use "" as wildcard)
vehicles := r.Query("", reasoner.RDFType, "http://example.org/Vehicle")
for _, t := range vehicles {
    fmt.Printf("%s is a Vehicle\n", t.Subject)
}

// Get all types for a specific instance
types := r.GetInferredTypes("http://example.org/myCar")
fmt.Printf("Types of myCar: %v\n", types)
```

## API Reference

### Main Functions

| Function                                                                | Description                                                |
| ----------------------------------------------------------------------- | ---------------------------------------------------------- |
| `ForwardReason(abox, tbox string) ([]string, error)`                    | Main API - returns all triples as strings                  |
| `ForwardReasonWithDetails(abox, tbox string) (*ReasoningResult, error)` | Returns detailed results with original/inferred separation |

### Reasoner Methods

| Method                                              | Description                                                       |
| --------------------------------------------------- | ----------------------------------------------------------------- |
| `NewReasoner() *Reasoner`                           | Create a new reasoner with default rules                          |
| `LoadTurtle(content string) error`                  | Parse and load Turtle content                                     |
| `RunForwardReasoning() int`                         | Apply all rules until fixpoint, returns count of inferred triples |
| `GetAllTriples() []string`                          | Get all triples as N-Triples strings                              |
| `GetInferredTypes(subject string) []string`         | Get all rdf:type values for a subject                             |
| `Query(subject, predicate, object string) []Triple` | Pattern matching query (use "" as wildcard)                       |

### Supported Inference Rules

1. **rdfs:subClassOf transitivity**: If A ⊑ B and B ⊑ C, then A ⊑ C
2. **rdf:type inheritance**: If x:A and A ⊑ B, then x:B
3. **rdfs:domain inference**: If P domain C and x P y, then x:C
4. **rdfs:range inference**: If P range C and x P y, then y:C
5. **rdfs:subPropertyOf transitivity**: If P1 ⊑ P2 and P2 ⊑ P3, then P1 ⊑ P3
6. **rdfs:subPropertyOf inheritance**: If P1 ⊑ P2 and x P1 y, then x P2 y
7. **owl:equivalentClass symmetry/transitivity**
8. **owl:sameAs symmetry/transitivity**
9. **owl:inverseOf inference**: If P1 inverse P2 and x P1 y, then y P2 x
10. **owl:TransitiveProperty inference**
11. **owl:SymmetricProperty inference**

## File Structure

```
reasoner/
├── go.mod              # Module definition
├── reasoner.go         # Main API and Reasoner type
├── parser.go           # Turtle format parser
├── store.go            # In-memory triple store
├── rules.go            # Forward reasoning rules
├── reasoner_test.go    # Unit tests
└── example/
    └── main.go         # Usage example
```

## Running Tests

```bash
go test -v
```

## Running the Example

```bash
go run example/main.go
```

## License

MIT
