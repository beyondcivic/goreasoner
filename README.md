# goreasoner

[![Version](https://img.shields.io/badge/version-v0.3.0-blue)](https://github.com/beyondcivic/goreasoner/releases/tag/v0.3.0)
[![Go Version](https://img.shields.io/badge/Go-1.24+-00ADD8?logo=go)](https://golang.org/doc/devel/release.html)
[![Go Reference](https://pkg.go.dev/badge/github.com/beyondcivic/goreasoner.svg)](https://pkg.go.dev/github.com/beyondcivic/goreasoner)
[![License](https://img.shields.io/badge/license-MIT-green)](LICENSE)

A Go implementation of a forward reasoner for RDF/OWL ontologies. This library provides both command-line interface and Go library for semantic reasoning, parsing Turtle format inputs and applying RDFS/OWL inference rules to derive new facts from TBox (terminology/schema) and ABox (assertions/instances).

## Overview

The Semantic Web relies on ontologies and inference to derive implicit knowledge from explicit facts. This tool streamlines RDF/OWL reasoning by:

- **Parsing Turtle format RDF data** with full prefix and IRI support
- **Applying forward reasoning rules** based on RDFS and OWL specifications
- **Deriving new knowledge** through transitive, symmetric, and hierarchical inference
- **Providing efficient triple storage** with indexed lookups for fast querying
- **Offering both CLI and library interfaces** for different integration needs

This project provides both a command-line interface and a Go library for semantic reasoning on RDF/OWL data.

## Key Features

- ✅ **Turtle Parser**: Custom parser (no external dependencies) supporting prefixes, IRIs, blank nodes, and literals
- ✅ **Forward Reasoning**: Complete RDFS/OWL inference rule implementation
- ✅ **Class Hierarchies**: Transitive subclass relationships and type inheritance
- ✅ **Property Reasoning**: Domain/range inference and property hierarchies
- ✅ **OWL Support**: Equivalent classes, same-as reasoning, inverse and transitive properties
- ✅ **Datalog Reasoning**: Built-in Datalog parser and evaluator for rules, facts, and boolean queries
- ✅ **Multiple Output Formats**: N-Triples and Datalog output formats supported
- ✅ **CLI & Library**: Both command-line tool and Go library interfaces
- ✅ **Cross-platform**: Works on Linux, macOS, and Windows

## Getting Started

### Prerequisites

- Go 1.24 or later
- Nix 2.25.4 or later (optional but recommended)
- PowerShell v7.5.1 or later (for building)

### Installation

#### Option 1: Install from Source

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/goreasoner.git
cd goreasoner
```

2. Build the application:

```bash
go build -o goreasoner .
```

#### Option 2: Using Nix (Recommended)

1. Clone the repository:

```bash
git clone https://github.com/beyondcivic/goreasoner.git
cd goreasoner
```

2. Prepare the environment using Nix flakes:

```bash
nix develop
```

3. Build the application:

```bash
./build.ps1
```

#### Option 3: Go Install

```bash
go install github.com/beyondcivic/goreasoner@latest
```

## Quick Start

### Command Line Interface

The `goreasoner` tool provides commands for semantic reasoning on RDF data:

```bash
# Run forward reasoning on RDF data (N-Triples output)
goreasoner run instances.ttl schema.ttl -o results.nt

# Run forward reasoning with Datalog output
goreasoner run instances.ttl schema.ttl --outputType=datalog -o results.dl

# Query a Datalog program
goreasoner dlquery results.dl "?- type(myTesla, Vehicle)."

# Show version information
goreasoner version
```

### Go Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/beyondcivic/goreasoner/pkg/reasoner"
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
        log.Fatalf("Error running forward reasoning: %v", err)
    }

    fmt.Printf("Inferred %d total triples\n", len(triples))
    for _, t := range triples {
        fmt.Println(t)
    }
}
```

## Detailed Command Reference

### `run` - Execute Forward Reasoning

Run forward reasoning on RDF data using TBox (schema) and ABox (instances).

```bash
goreasoner run [ABOX_FILE] [TBOX_FILE] [OPTIONS]
```

**Options:**

- `-o, --output`: Output file path (default: `[abox_filename]_inferred.nt`)
- `--outputType`: Output format - `ntriple` or `datalog` (default: `ntriple`)

**Examples:**

```bash
# Basic reasoning (N-Triples output)
goreasoner run instances.ttl schema.ttl

# With custom output path
goreasoner run instances.ttl schema.ttl -o my-results.nt

# Datalog output format
goreasoner run instances.ttl schema.ttl --outputType=datalog

# Datalog output with custom file
goreasoner run instances.ttl schema.ttl --outputType=datalog -o results.dl
```

### `dlquery` - Query a Datalog Program

Evaluate a boolean query against a Datalog program (facts and rules).

```bash
goreasoner dlquery [DATALOG_FILE] [QUERY]
```

**Arguments:**

- `DATALOG_FILE`: Path to a `.dl` file containing Datalog facts and rules
- `QUERY`: A Datalog query string in `?- predicate(args).` format

**Examples:**

```bash
# Query a ground fact
goreasoner dlquery data.dl "?- type(myTesla, Car)."

# Query a derived fact
goreasoner dlquery data.dl "?- Ancestor(john, jane)."

# Query with variable (returns true if any binding exists)
goreasoner dlquery data.dl "?- type(X, Vehicle)."
```

### `version` - Show Version Information

Display version, build information, and system details.

```bash
goreasoner version
```

## Supported Inference Rules

The reasoner implements comprehensive RDFS/OWL inference rules:

| Rule Category                    | Description                                | Example                   |
| -------------------------------- | ------------------------------------------ | ------------------------- |
| **rdfs:subClassOf transitivity** | If A ⊑ B and B ⊑ C, then A ⊑ C             | Car ⊑ Vehicle ⊑ Transport |
| **rdf:type inheritance**         | If x:A and A ⊑ B, then x:B                 | myCar:Car → myCar:Vehicle |
| **rdfs:domain inference**        | If P domain C and x P y, then x:C          | hasOwner domain Person    |
| **rdfs:range inference**         | If P range C and x P y, then y:C           | hasAge range Integer      |
| **rdfs:subPropertyOf**           | Property hierarchy reasoning               | drives ⊑ operates         |
| **owl:equivalentClass**          | Class equivalence (symmetric/transitive)   | Vehicle ≡ Automobile      |
| **owl:sameAs**                   | Individual identity (symmetric/transitive) | person1 ≡ person2         |
| **owl:inverseOf**                | Inverse property inference                 | owns ⟷ isOwnedBy          |
| **owl:TransitiveProperty**       | Transitive property chains                 | locatedIn transitivity    |
| **owl:SymmetricProperty**        | Symmetric property inference               | marriedTo symmetry        |

## Output Formats

The tool supports two output formats for reasoning results:

### N-Triples Format (Default)

Standard N-Triples format with full IRIs:

```
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Car> .
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Vehicle> .
<http://example.org/myCar> <http://www.w3.org/1999/02/22-rdf-syntax-ns#type> <http://example.org/Transport> .
<http://example.org/Car> <http://www.w3.org/2000/01/rdf-schema#subClassOf> <http://example.org/Vehicle> .
<http://example.org/Vehicle> <http://www.w3.org/2000/01/rdf-schema#subClassOf> <http://example.org/Transport> .
```

### Datalog Format

Datalog facts format with simplified names, suitable for Datalog reasoning systems:

```
type(myCar, Car).
type(myCar, Vehicle).
type(myCar, Transport).
subClassOf(Car, Vehicle).
subClassOf(Vehicle, Transport).
```

The Datalog format converts RDF triples `<subject> <predicate> <object>` to facts `predicate(subject, object).` and simplifies IRIs by extracting local names.

## Datalog Reasoning

In addition to RDF/OWL forward reasoning, goreasoner includes a built-in Datalog evaluator that supports facts, rules with variables, recursive rules, and boolean queries. The evaluator uses **naive bottom-up (forward-chaining) evaluation** with fixed-point computation to derive all possible facts before answering queries.

### Datalog Syntax

#### Facts

Ground atoms terminated by a period:

```
Parent(john, mary).
Type(myTesla, Car).
```

#### Rules

Horn clauses with a head and one or more body atoms separated by `:-`:

```
Ancestor(X, Y) :- Parent(X, Y).
Ancestor(X, Z) :- Parent(X, Y), Ancestor(Y, Z).
```

#### Variables

Variables are recognized by these conventions:

- Single uppercase letter: `X`, `Y`, `Z`
- All-uppercase identifiers: `VAR_X`, `PERSON`, `NODE1`
- `?`-prefixed identifiers: `?x`, `?person`

Everything else is treated as a constant.

#### Queries

Queries use the `?-` prefix and return a boolean (true if any matching fact exists):

```
?- Ancestor(john, jane).
?- type(X, Vehicle).
```

#### Comments

Both Prolog-style (`%`) and C-style (`//`) line comments are supported:

```
% This is a comment
Parent(john, mary).  // This is also a comment
```

### Datalog Go Library Usage

```go
package main

import (
    "fmt"
    "log"

    "github.com/beyondcivic/goreasoner/pkg/reasoner"
)

func main() {
    program := `
Parent(john, mary).
Parent(mary, jane).
Ancestor(X, Y) :- Parent(X, Y).
Ancestor(X, Z) :- Parent(X, Y), Ancestor(Y, Z).
`

    result, err := reasoner.DLQuery(program, "?- Ancestor(john, jane).")
    if err != nil {
        log.Fatalf("Error: %v", err)
    }
    fmt.Println(result) // true
}
```

For more control, you can use the lower-level API directly:

```go
program, err := reasoner.ParseDatalog(input)
if err != nil {
    log.Fatal(err)
}

// Derive all facts via fixed-point evaluation
derivedFacts := program.Reason()

// Parse and evaluate a query
query, err := reasoner.ParseQuery("?- Ancestor(john, jane).")
if err != nil {
    log.Fatal(err)
}

satisfied := program.EvaluateQuery(query, derivedFacts)
```

### Datalog API Reference

#### `DLQuery(datalogContent, queryStr string) (bool, error)`

Main API function for Datalog querying. Parses the program, runs reasoning to fixed point, and evaluates the query.

#### `ParseDatalog(input string) (*DatalogProgram, error)`

Parses a Datalog program string into a `DatalogProgram` containing facts and rules.

#### `ParseQuery(s string) (DLAtom, error)`

Parses a query string (with or without `?-` prefix) into a `DLAtom`.

#### `(*DatalogProgram) Reason() []DLAtom`

Runs forward-chaining evaluation until no new facts are derived. Returns all ground facts (original and inferred).

#### `(*DatalogProgram) EvaluateQuery(query DLAtom, derivedFacts []DLAtom) bool`

Checks whether a query matches any derived fact. Variables in the query act as wildcards.

### Datalog Limitations

The Datalog evaluator is designed for simple, positive Datalog programs. Be aware of the following limitations:

- **No negation**: There is no support for negation-as-failure (`not`, `\+`) in rule bodies. Only positive atoms can appear in rules.
- **No built-in comparisons or arithmetic**: Operators such as `!=`, `<`, `>`, and arithmetic expressions are not supported. All terms are symbolic constants or variables.
- **Boolean queries only**: `DLQuery` returns `true`/`false`. It does not return variable bindings. For example, querying `?- Ancestor(john, X).` will tell you whether any ancestor exists, but will not enumerate them.
- **No safety checks on rules**: The parser accepts rules where the head contains variables that do not appear in the body (e.g., `Foo(X) :- Bar(Y).`). Such rules will not produce incorrect results (ungrounded heads are silently discarded), but no warning is emitted.
- **No indexing on facts**: The evaluator performs a linear scan over all facts when matching rule body atoms. This is adequate for small to medium programs but may become slow with thousands of facts.
- **No aggregation or constraints**: Features like `count`, `min`, `max`, or integrity constraints found in extended Datalog systems are not supported.

## Examples

### Example 1: Basic Class Hierarchy Reasoning

```bash
# Create schema file (tbox.ttl)
echo '@prefix ex: <http://example.org/> .
@prefix rdfs: <http://www.w3.org/2000/01/rdf-schema#> .

ex:SportsCar rdfs:subClassOf ex:Car .
ex:Car rdfs:subClassOf ex:Vehicle .' > tbox.ttl

# Create instance file (abox.ttl)
echo '@prefix ex: <http://example.org/> .
@prefix rdf: <http://www.w3.org/1999/02/22-rdf-syntax-ns#> .

ex:myFerrari rdf:type ex:SportsCar .' > abox.ttl

# Run reasoning (N-Triples output)
goreasoner run abox.ttl tbox.ttl -o results.nt

# Run reasoning (Datalog output)
goreasoner run abox.ttl tbox.ttl --outputType=datalog -o results.dl

# View results
cat results.nt
cat results.dl
```

### Example 2: Complex Ontology Reasoning

Given an ontology with domain/range restrictions, property hierarchies, and equivalent classes, the reasoner will derive all implicit knowledge according to RDFS/OWL semantics.

## API Reference

### Core Functions

#### `ForwardReason(abox, tbox string) ([]string, error)`

Main API function that performs forward reasoning on RDF data.

**Parameters:**

- `abox`: Turtle string containing instance data (assertions)
- `tbox`: Turtle string containing schema/ontology definitions

**Returns:**

- `[]string`: List of all triples (including inferred) in N-Triples format
- `error`: Any parsing or processing errors

#### `ForwardReasonWithDetails(abox, tbox string) (*ReasoningResult, error)`

Returns detailed reasoning results with original/inferred triple separation.

**Parameters:**

- `abox`: Turtle string containing instance data
- `tbox`: Turtle string containing schema definitions

**Returns:**

- `*ReasoningResult`: Detailed results structure
- `error`: Any parsing or processing errors

### Using the Reasoner Directly

```go
r := reasoner.NewReasoner()

// Load TBox (schema)
if err := r.LoadTurtle(tbox); err != nil {
    log.Fatalf("Error loading TBox: %v", err)
}

// Load ABox (instances)
if err := r.LoadTurtle(abox); err != nil {
    log.Fatalf("Error loading ABox: %v", err)
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

### Data Structures

#### `ReasoningResult`

Represents detailed reasoning results:

```go
type ReasoningResult struct {
    OriginalTriples []string // Triples from input
    InferredTriples []string // Newly inferred triples
    AllTriples      []string // All triples combined
    OriginalCount   int      // Number of original triples
    InferredCount   int      // Number of inferred triples
    TotalCount      int      // Total number of triples
}
```

#### `Triple`

Represents an RDF triple:

```go
type Triple struct {
    Subject   string
    Predicate string
    Object    string
}
```

#### `Reasoner`

Main reasoner structure with methods:

```go
type Reasoner struct {
    store  *TripleStore
    rules  []Rule
    parser *TurtleParser
}
```

### Reasoner Methods

| Method                                              | Description                                                       |
| --------------------------------------------------- | ----------------------------------------------------------------- |
| `NewReasoner() *Reasoner`                           | Create a new reasoner with default rules                          |
| `NewReasonerWithRules(rules []Rule) *Reasoner`      | Create a reasoner with custom rules                               |
| `LoadTurtle(content string) error`                  | Parse and load Turtle content                                     |
| `RunForwardReasoning() int`                         | Apply all rules until fixpoint, returns count of inferred triples |
| `GetAllTriples() []string`                          | Get all triples as N-Triples strings                              |
| `GetInferredTypes(subject string) []string`         | Get all rdf:type values for a subject                             |
| `Query(subject, predicate, object string) []Triple` | Pattern matching query (use "" as wildcard)                       |
| `GetStore() *TripleStore`                           | Access the underlying triple store                                |

## Architecture

The library is organized into several key components:

### Core Package (`pkg/reasoner`)

- **Forward Reasoning Engine**: Rule-based inference with fixpoint computation
- **Turtle Parser**: Complete Turtle format parser with prefix support
- **Triple Store**: Indexed in-memory storage for efficient querying
- **Rule System**: Modular RDFS/OWL inference rules
- **Datalog Evaluator**: Parser and naive bottom-up reasoner for Datalog programs
- **Query Interface**: Pattern matching, type inference, and Datalog queries

### Command Line Interface (`cmd/goreasoner`)

- **Cobra-based CLI** with subcommands for reasoning operations
- **File I/O handling** for Turtle input and N-Triples output
- **Comprehensive help system** with detailed usage examples
- **Flexible output options** and error handling

### Contributing

1. Fork the repository
2. Create a feature branch: `git checkout -b feature/new-feature`
3. Make your changes and add tests
4. Ensure all tests pass: `go test ./...`
5. Commit your changes: `git commit -am 'Add new feature'`
6. Push to the branch: `git push origin feature/new-feature`
7. Submit a pull request

### Testing

Run the test suite:

```bash
go test ./...
```

Run tests with coverage:

```bash
go test -cover ./...
```

## Build Environment

### Using Nix (Recommended)

Use Nix flakes to set up the build environment:

```bash
nix develop
```

### Manual Build

Check the build arguments in `build.ps1`:

```bash
# Build static binary with version information
$env:CGO_ENABLED = "1"
$env:GOOS = "linux"
$env:GOARCH = "amd64"
```

Then run:

```bash
./build.ps1
```

Or build manually:

```bash
go build -o goreasoner .
```

## File Structure

```
goreasoner/
├── go.mod                    # Module definition
├── main.go                   # Main entry point
├── build.ps1                 # Build script
├── flake.nix                 # Nix flake configuration
├── cmd/
│   └── goreasoner/
│       ├── main.go           # CLI interface
│       └── commands.go       # Command definitions
├── pkg/
│   ├── reasoner/
│   │   ├── core.go           # Main API and Reasoner type
│   │   ├── parser.go         # Turtle format parser
│   │   ├── store.go          # In-memory triple store
│   │   ├── rules.go          # Forward reasoning rules
│   │   ├── datalog.go        # Datalog parser and reasoner
│   │   ├── utils.go          # Utility functions
│   │   └── error.go          # Error handling
│   └── version/
│       └── version.go        # Version information
└── docs/
    └── docs-gen.go           # Documentation generator
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
