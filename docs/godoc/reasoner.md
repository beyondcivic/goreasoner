# reasoner

```go
import "github.com/beyondcivic/goreasoner/pkg/reasoner"
```

Package reasoner provides forward reasoning capabilities for RDF/OWL ontologies. It parses Turtle format inputs and applies RDFS/OWL inference rules to derive new facts from the given TBox \(terminology\) and ABox \(assertions\).

## Index

- [Constants](<#constants>)
- [func ConvertTriplesToDatalog\(triples \[\]string\) \[\]string](<#ConvertTriplesToDatalog>)
- [func DLQuery\(datalogContent, queryStr string\) \(bool, error\)](<#DLQuery>)
- [func ForwardReason\(abox, tbox string\) \(\[\]string, error\)](<#ForwardReason>)
- [type AppError](<#AppError>)
  - [func \(e AppError\) Error\(\) string](<#AppError.Error>)
- [type DLAtom](<#DLAtom>)
  - [func ParseQuery\(s string\) \(DLAtom, error\)](<#ParseQuery>)
  - [func \(a DLAtom\) String\(\) string](<#DLAtom.String>)
- [type DLRule](<#DLRule>)
- [type DLTerm](<#DLTerm>)
- [type DatalogProgram](<#DatalogProgram>)
  - [func ParseDatalog\(input string\) \(\*DatalogProgram, error\)](<#ParseDatalog>)
  - [func \(p \*DatalogProgram\) EvaluateQuery\(query DLAtom, derivedFacts \[\]DLAtom\) bool](<#DatalogProgram.EvaluateQuery>)
  - [func \(p \*DatalogProgram\) Reason\(\) \[\]DLAtom](<#DatalogProgram.Reason>)
- [type DomainInference](<#DomainInference>)
  - [func \(r \*DomainInference\) Apply\(store \*TripleStore\) \[\]Triple](<#DomainInference.Apply>)
  - [func \(r \*DomainInference\) Name\(\) string](<#DomainInference.Name>)
- [type EquivalentClassSymmetry](<#EquivalentClassSymmetry>)
  - [func \(r \*EquivalentClassSymmetry\) Apply\(store \*TripleStore\) \[\]Triple](<#EquivalentClassSymmetry.Apply>)
  - [func \(r \*EquivalentClassSymmetry\) Name\(\) string](<#EquivalentClassSymmetry.Name>)
- [type EquivalentClassTransitivity](<#EquivalentClassTransitivity>)
  - [func \(r \*EquivalentClassTransitivity\) Apply\(store \*TripleStore\) \[\]Triple](<#EquivalentClassTransitivity.Apply>)
  - [func \(r \*EquivalentClassTransitivity\) Name\(\) string](<#EquivalentClassTransitivity.Name>)
- [type InversePropertyInference](<#InversePropertyInference>)
  - [func \(r \*InversePropertyInference\) Apply\(store \*TripleStore\) \[\]Triple](<#InversePropertyInference.Apply>)
  - [func \(r \*InversePropertyInference\) Name\(\) string](<#InversePropertyInference.Name>)
- [type RangeInference](<#RangeInference>)
  - [func \(r \*RangeInference\) Apply\(store \*TripleStore\) \[\]Triple](<#RangeInference.Apply>)
  - [func \(r \*RangeInference\) Name\(\) string](<#RangeInference.Name>)
- [type Reasoner](<#Reasoner>)
  - [func NewReasoner\(\) \*Reasoner](<#NewReasoner>)
  - [func NewReasonerWithRules\(rules \[\]Rule\) \*Reasoner](<#NewReasonerWithRules>)
  - [func \(r \*Reasoner\) GetAllTriples\(\) \[\]string](<#Reasoner.GetAllTriples>)
  - [func \(r \*Reasoner\) GetInferredTypes\(subject string\) \[\]string](<#Reasoner.GetInferredTypes>)
  - [func \(r \*Reasoner\) GetStore\(\) \*TripleStore](<#Reasoner.GetStore>)
  - [func \(r \*Reasoner\) LoadTurtle\(content string\) error](<#Reasoner.LoadTurtle>)
  - [func \(r \*Reasoner\) Query\(subject, predicate, object string\) \[\]Triple](<#Reasoner.Query>)
  - [func \(r \*Reasoner\) RunForwardReasoning\(\) int](<#Reasoner.RunForwardReasoning>)
- [type ReasoningResult](<#ReasoningResult>)
  - [func ForwardReasonWithDetails\(abox, tbox string\) \(\*ReasoningResult, error\)](<#ForwardReasonWithDetails>)
  - [func \(r \*ReasoningResult\) String\(\) string](<#ReasoningResult.String>)
- [type Rule](<#Rule>)
  - [func DefaultRules\(\) \[\]Rule](<#DefaultRules>)
- [type SameAsSymmetry](<#SameAsSymmetry>)
  - [func \(r \*SameAsSymmetry\) Apply\(store \*TripleStore\) \[\]Triple](<#SameAsSymmetry.Apply>)
  - [func \(r \*SameAsSymmetry\) Name\(\) string](<#SameAsSymmetry.Name>)
- [type SameAsTransitivity](<#SameAsTransitivity>)
  - [func \(r \*SameAsTransitivity\) Apply\(store \*TripleStore\) \[\]Triple](<#SameAsTransitivity.Apply>)
  - [func \(r \*SameAsTransitivity\) Name\(\) string](<#SameAsTransitivity.Name>)
- [type SubClassTransitivity](<#SubClassTransitivity>)
  - [func \(r \*SubClassTransitivity\) Apply\(store \*TripleStore\) \[\]Triple](<#SubClassTransitivity.Apply>)
  - [func \(r \*SubClassTransitivity\) Name\(\) string](<#SubClassTransitivity.Name>)
- [type SubPropertyInheritance](<#SubPropertyInheritance>)
  - [func \(r \*SubPropertyInheritance\) Apply\(store \*TripleStore\) \[\]Triple](<#SubPropertyInheritance.Apply>)
  - [func \(r \*SubPropertyInheritance\) Name\(\) string](<#SubPropertyInheritance.Name>)
- [type SubPropertyTransitivity](<#SubPropertyTransitivity>)
  - [func \(r \*SubPropertyTransitivity\) Apply\(store \*TripleStore\) \[\]Triple](<#SubPropertyTransitivity.Apply>)
  - [func \(r \*SubPropertyTransitivity\) Name\(\) string](<#SubPropertyTransitivity.Name>)
- [type SymmetricPropertyInference](<#SymmetricPropertyInference>)
  - [func \(r \*SymmetricPropertyInference\) Apply\(store \*TripleStore\) \[\]Triple](<#SymmetricPropertyInference.Apply>)
  - [func \(r \*SymmetricPropertyInference\) Name\(\) string](<#SymmetricPropertyInference.Name>)
- [type TransitivePropertyInference](<#TransitivePropertyInference>)
  - [func \(r \*TransitivePropertyInference\) Apply\(store \*TripleStore\) \[\]Triple](<#TransitivePropertyInference.Apply>)
  - [func \(r \*TransitivePropertyInference\) Name\(\) string](<#TransitivePropertyInference.Name>)
- [type Triple](<#Triple>)
  - [func \(t Triple\) String\(\) string](<#Triple.String>)
- [type TripleStore](<#TripleStore>)
  - [func NewTripleStore\(\) \*TripleStore](<#NewTripleStore>)
  - [func \(ts \*TripleStore\) Add\(t Triple\) bool](<#TripleStore.Add>)
  - [func \(ts \*TripleStore\) All\(\) \[\]Triple](<#TripleStore.All>)
  - [func \(ts \*TripleStore\) Contains\(t Triple\) bool](<#TripleStore.Contains>)
  - [func \(ts \*TripleStore\) FindByObject\(object string\) \[\]Triple](<#TripleStore.FindByObject>)
  - [func \(ts \*TripleStore\) FindByPredicate\(predicate string\) \[\]Triple](<#TripleStore.FindByPredicate>)
  - [func \(ts \*TripleStore\) FindByPredicateObject\(predicate, object string\) \[\]Triple](<#TripleStore.FindByPredicateObject>)
  - [func \(ts \*TripleStore\) FindBySubject\(subject string\) \[\]Triple](<#TripleStore.FindBySubject>)
  - [func \(ts \*TripleStore\) FindBySubjectPredicate\(subject, predicate string\) \[\]Triple](<#TripleStore.FindBySubjectPredicate>)
  - [func \(ts \*TripleStore\) Size\(\) int](<#TripleStore.Size>)
- [type TurtleParser](<#TurtleParser>)
  - [func NewTurtleParser\(\) \*TurtleParser](<#NewTurtleParser>)
  - [func \(p \*TurtleParser\) FallbackParse\(content string\) \(\[\]Triple, error\)](<#TurtleParser.FallbackParse>)
  - [func \(p \*TurtleParser\) Parse\(content string\) \(\[\]Triple, error\)](<#TurtleParser.Parse>)
- [type TypeInheritance](<#TypeInheritance>)
  - [func \(r \*TypeInheritance\) Apply\(store \*TripleStore\) \[\]Triple](<#TypeInheritance.Apply>)
  - [func \(r \*TypeInheritance\) Name\(\) string](<#TypeInheritance.Name>)


## Constants

<a name="RDFType"></a>Common RDF/RDFS/OWL URIs

```go
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
```

<a name="ConvertTriplesToDatalog"></a>
## func [ConvertTriplesToDatalog](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/utils.go#L11>)

```go
func ConvertTriplesToDatalog(triples []string) []string
```

ConvertTriplesToDatalog converts a list of N\-Triple strings to Datalog format Each RDF triple \(subject, predicate, object\) becomes a Datalog fact: predicate\(subject, object\) IRIs are converted to simplified names by extracting the local part after \# or /

<a name="DLQuery"></a>
## func [DLQuery](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L313>)

```go
func DLQuery(datalogContent, queryStr string) (bool, error)
```

DLQuery is the main public API function for Datalog querying. It accepts a Datalog program and a query, performs reasoning, and returns true if the query is satisfied.

<a name="ForwardReason"></a>
## func [ForwardReason](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L149>)

```go
func ForwardReason(abox, tbox string) ([]string, error)
```

ForwardReason is the main public API function. It accepts TBox \(terminology/schema\) and ABox \(assertions/instances\) in Turtle format, performs forward reasoning, and returns all inferred triples as strings.

Parameters:

- abox: Turtle string containing instance data \(assertions\)
- tbox: Turtle string containing schema/ontology definitions

Returns:

- \[\]string: List of all triples \(including inferred\) in N\-Triples format
- error: Any parsing or processing errors

<a name="AppError"></a>
## type [AppError](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/error.go#L5-L10>)



```go
type AppError struct {
    // Message to show the user.
    Message string
    // Value to include with message
    Value any
}
```

<a name="AppError.Error"></a>
### func \(AppError\) [Error](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/error.go#L12>)

```go
func (e AppError) Error() string
```



<a name="DLAtom"></a>
## type [DLAtom](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L16-L19>)

DLAtom represents a Datalog atom: Predicate\(DLTerm1, DLTerm2, ...\)

```go
type DLAtom struct {
    Predicate string
    Terms     []DLTerm
}
```

<a name="ParseQuery"></a>
### func [ParseQuery](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L249>)

```go
func ParseQuery(s string) (DLAtom, error)
```

ParseQuery parses a Datalog query

<a name="DLAtom.String"></a>
### func \(DLAtom\) [String](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L33>)

```go
func (a DLAtom) String() string
```



<a name="DLRule"></a>
## type [DLRule](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L22-L25>)

DLRule represents a Datalog rule: Head :\- Body1, Body2, ...

```go
type DLRule struct {
    Head DLAtom
    Body []DLAtom
}
```

<a name="DLTerm"></a>
## type [DLTerm](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L10-L13>)

DLTerm represents a Datalog term \(either a constant or a variable\)

```go
type DLTerm struct {
    Value      string
    IsVariable bool
}
```

<a name="DatalogProgram"></a>
## type [DatalogProgram](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L28-L31>)

DatalogProgram represents a collection of facts and rules

```go
type DatalogProgram struct {
    Facts []DLAtom
    Rules []DLRule
}
```

<a name="ParseDatalog"></a>
### func [ParseDatalog](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L42>)

```go
func ParseDatalog(input string) (*DatalogProgram, error)
```

ParseDatalog parses a Datalog program from a string

<a name="DatalogProgram.EvaluateQuery"></a>
### func \(\*DatalogProgram\) [EvaluateQuery](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L393>)

```go
func (p *DatalogProgram) EvaluateQuery(query DLAtom, derivedFacts []DLAtom) bool
```

EvaluateQuery checks if a query matches any derived facts

<a name="DatalogProgram.Reason"></a>
### func \(\*DatalogProgram\) [Reason](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/datalog.go#L259>)

```go
func (p *DatalogProgram) Reason() []DLAtom
```

Reason evaluates the Datalog program and returns all derived facts

<a name="DomainInference"></a>
## type [DomainInference](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L93>)

DomainInference implements rdfs:domain inference If P rdfs:domain C and X P Y, then X rdf:type C

```go
type DomainInference struct{}
```

<a name="DomainInference.Apply"></a>
### func \(\*DomainInference\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L99>)

```go
func (r *DomainInference) Apply(store *TripleStore) []Triple
```



<a name="DomainInference.Name"></a>
### func \(\*DomainInference\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L95>)

```go
func (r *DomainInference) Name() string
```



<a name="EquivalentClassSymmetry"></a>
## type [EquivalentClassSymmetry](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L218>)

EquivalentClassSymmetry implements owl:equivalentClass symmetry If A owl:equivalentClass B, then B owl:equivalentClass A

```go
type EquivalentClassSymmetry struct{}
```

<a name="EquivalentClassSymmetry.Apply"></a>
### func \(\*EquivalentClassSymmetry\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L224>)

```go
func (r *EquivalentClassSymmetry) Apply(store *TripleStore) []Triple
```



<a name="EquivalentClassSymmetry.Name"></a>
### func \(\*EquivalentClassSymmetry\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L220>)

```go
func (r *EquivalentClassSymmetry) Name() string
```



<a name="EquivalentClassTransitivity"></a>
## type [EquivalentClassTransitivity](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L241>)

EquivalentClassTransitivity implements owl:equivalentClass transitivity If A owl:equivalentClass B and B owl:equivalentClass C, then A owl:equivalentClass C

```go
type EquivalentClassTransitivity struct{}
```

<a name="EquivalentClassTransitivity.Apply"></a>
### func \(\*EquivalentClassTransitivity\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L247>)

```go
func (r *EquivalentClassTransitivity) Apply(store *TripleStore) []Triple
```



<a name="EquivalentClassTransitivity.Name"></a>
### func \(\*EquivalentClassTransitivity\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L243>)

```go
func (r *EquivalentClassTransitivity) Name() string
```



<a name="InversePropertyInference"></a>
## type [InversePropertyInference](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L320>)

InversePropertyInference implements owl:inverseOf If P1 owl:inverseOf P2 and X P1 Y, then Y P2 X

```go
type InversePropertyInference struct{}
```

<a name="InversePropertyInference.Apply"></a>
### func \(\*InversePropertyInference\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L326>)

```go
func (r *InversePropertyInference) Apply(store *TripleStore) []Triple
```



<a name="InversePropertyInference.Name"></a>
### func \(\*InversePropertyInference\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L322>)

```go
func (r *InversePropertyInference) Name() string
```



<a name="RangeInference"></a>
## type [RangeInference](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L125>)

RangeInference implements rdfs:range inference If P rdfs:range C and X P Y, then Y rdf:type C

```go
type RangeInference struct{}
```

<a name="RangeInference.Apply"></a>
### func \(\*RangeInference\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L131>)

```go
func (r *RangeInference) Apply(store *TripleStore) []Triple
```



<a name="RangeInference.Name"></a>
### func \(\*RangeInference\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L127>)

```go
func (r *RangeInference) Name() string
```



<a name="Reasoner"></a>
## type [Reasoner](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L12-L16>)

Reasoner performs forward reasoning on RDF data

```go
type Reasoner struct {
    // contains filtered or unexported fields
}
```

<a name="NewReasoner"></a>
### func [NewReasoner](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L19>)

```go
func NewReasoner() *Reasoner
```

NewReasoner creates a new reasoner with default rules

<a name="NewReasonerWithRules"></a>
### func [NewReasonerWithRules](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L28>)

```go
func NewReasonerWithRules(rules []Rule) *Reasoner
```

NewReasonerWithRules creates a new reasoner with custom rules

<a name="Reasoner.GetAllTriples"></a>
### func \(\*Reasoner\) [GetAllTriples](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L78>)

```go
func (r *Reasoner) GetAllTriples() []string
```

GetAllTriples returns all triples in the store as strings

<a name="Reasoner.GetInferredTypes"></a>
### func \(\*Reasoner\) [GetInferredTypes](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L89>)

```go
func (r *Reasoner) GetInferredTypes(subject string) []string
```

GetInferredTypes returns all rdf:type assertions for a given subject

<a name="Reasoner.GetStore"></a>
### func \(\*Reasoner\) [GetStore](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L134>)

```go
func (r *Reasoner) GetStore() *TripleStore
```

GetStore returns the underlying triple store

<a name="Reasoner.LoadTurtle"></a>
### func \(\*Reasoner\) [LoadTurtle](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L37>)

```go
func (r *Reasoner) LoadTurtle(content string) error
```

LoadTurtle parses and loads Turtle content into the store

<a name="Reasoner.Query"></a>
### func \(\*Reasoner\) [Query](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L100>)

```go
func (r *Reasoner) Query(subject, predicate, object string) []Triple
```

Query returns all triples matching the given pattern Use empty string "" as wildcard

<a name="Reasoner.RunForwardReasoning"></a>
### func \(\*Reasoner\) [RunForwardReasoning](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L52>)

```go
func (r *Reasoner) RunForwardReasoning() int
```

RunForwardReasoning applies all rules until no new facts are derived Returns the number of new triples inferred

<a name="ReasoningResult"></a>
## type [ReasoningResult](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L224-L231>)

ReasoningResult contains detailed results from forward reasoning

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

<a name="ForwardReasonWithDetails"></a>
### func [ForwardReasonWithDetails](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/core.go#L173>)

```go
func ForwardReasonWithDetails(abox, tbox string) (*ReasoningResult, error)
```

ForwardReasonWithDetails returns both original and inferred triples separately

<a name="ReasoningResult.String"></a>
### func \(\*ReasoningResult\) [String](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/utils.go#L43>)

```go
func (r *ReasoningResult) String() string
```

String returns a human\-readable summary of the reasoning result

<a name="Rule"></a>
## type [Rule](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L20-L25>)

Rule represents a forward reasoning rule

```go
type Rule interface {
    // Name returns the rule name
    Name() string
    // Apply applies the rule to the store and returns new inferred triples
    Apply(store *TripleStore) []Triple
}
```

<a name="DefaultRules"></a>
### func [DefaultRules](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L422>)

```go
func DefaultRules() []Rule
```

DefaultRules returns the default set of reasoning rules

<a name="SameAsSymmetry"></a>
## type [SameAsSymmetry](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L269>)

SameAsSymmetry implements owl:sameAs symmetry

```go
type SameAsSymmetry struct{}
```

<a name="SameAsSymmetry.Apply"></a>
### func \(\*SameAsSymmetry\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L275>)

```go
func (r *SameAsSymmetry) Apply(store *TripleStore) []Triple
```



<a name="SameAsSymmetry.Name"></a>
### func \(\*SameAsSymmetry\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L271>)

```go
func (r *SameAsSymmetry) Name() string
```



<a name="SameAsTransitivity"></a>
## type [SameAsTransitivity](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L291>)

SameAsTransitivity implements owl:sameAs transitivity

```go
type SameAsTransitivity struct{}
```

<a name="SameAsTransitivity.Apply"></a>
### func \(\*SameAsTransitivity\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L297>)

```go
func (r *SameAsTransitivity) Apply(store *TripleStore) []Triple
```



<a name="SameAsTransitivity.Name"></a>
### func \(\*SameAsTransitivity\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L293>)

```go
func (r *SameAsTransitivity) Name() string
```



<a name="SubClassTransitivity"></a>
## type [SubClassTransitivity](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L29>)

SubClassTransitivity implements rdfs:subClassOf transitivity If A rdfs:subClassOf B and B rdfs:subClassOf C, then A rdfs:subClassOf C

```go
type SubClassTransitivity struct{}
```

<a name="SubClassTransitivity.Apply"></a>
### func \(\*SubClassTransitivity\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L35>)

```go
func (r *SubClassTransitivity) Apply(store *TripleStore) []Triple
```



<a name="SubClassTransitivity.Name"></a>
### func \(\*SubClassTransitivity\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L31>)

```go
func (r *SubClassTransitivity) Name() string
```



<a name="SubPropertyInheritance"></a>
## type [SubPropertyInheritance](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L190>)

SubPropertyInheritance implements property inheritance If P1 rdfs:subPropertyOf P2 and X P1 Y, then X P2 Y

```go
type SubPropertyInheritance struct{}
```

<a name="SubPropertyInheritance.Apply"></a>
### func \(\*SubPropertyInheritance\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L196>)

```go
func (r *SubPropertyInheritance) Apply(store *TripleStore) []Triple
```



<a name="SubPropertyInheritance.Name"></a>
### func \(\*SubPropertyInheritance\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L192>)

```go
func (r *SubPropertyInheritance) Name() string
```



<a name="SubPropertyTransitivity"></a>
## type [SubPropertyTransitivity](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L161>)

SubPropertyTransitivity implements rdfs:subPropertyOf transitivity If P1 rdfs:subPropertyOf P2 and P2 rdfs:subPropertyOf P3, then P1 rdfs:subPropertyOf P3

```go
type SubPropertyTransitivity struct{}
```

<a name="SubPropertyTransitivity.Apply"></a>
### func \(\*SubPropertyTransitivity\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L167>)

```go
func (r *SubPropertyTransitivity) Apply(store *TripleStore) []Triple
```



<a name="SubPropertyTransitivity.Name"></a>
### func \(\*SubPropertyTransitivity\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L163>)

```go
func (r *SubPropertyTransitivity) Name() string
```



<a name="SymmetricPropertyInference"></a>
## type [SymmetricPropertyInference](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L394>)

SymmetricPropertyInference implements owl:SymmetricProperty If P is symmetric and X P Y, then Y P X

```go
type SymmetricPropertyInference struct{}
```

<a name="SymmetricPropertyInference.Apply"></a>
### func \(\*SymmetricPropertyInference\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L400>)

```go
func (r *SymmetricPropertyInference) Apply(store *TripleStore) []Triple
```



<a name="SymmetricPropertyInference.Name"></a>
### func \(\*SymmetricPropertyInference\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L396>)

```go
func (r *SymmetricPropertyInference) Name() string
```



<a name="TransitivePropertyInference"></a>
## type [TransitivePropertyInference](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L357>)

TransitivePropertyInference implements owl:TransitiveProperty If P is transitive and X P Y and Y P Z, then X P Z

```go
type TransitivePropertyInference struct{}
```

<a name="TransitivePropertyInference.Apply"></a>
### func \(\*TransitivePropertyInference\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L363>)

```go
func (r *TransitivePropertyInference) Apply(store *TripleStore) []Triple
```



<a name="TransitivePropertyInference.Name"></a>
### func \(\*TransitivePropertyInference\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L359>)

```go
func (r *TransitivePropertyInference) Name() string
```



<a name="Triple"></a>
## type [Triple](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L9-L13>)

Triple represents an RDF triple \(subject, predicate, object\)

```go
type Triple struct {
    Subject   string
    Predicate string
    Object    string
}
```

<a name="Triple.String"></a>
### func \(Triple\) [String](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L16>)

```go
func (t Triple) String() string
```

String returns the triple in N\-Triples format

<a name="TripleStore"></a>
## type [TripleStore](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L41-L49>)

TripleStore is an in\-memory store for RDF triples

```go
type TripleStore struct {
    // contains filtered or unexported fields
}
```

<a name="NewTripleStore"></a>
### func [NewTripleStore](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L52>)

```go
func NewTripleStore() *TripleStore
```

NewTripleStore creates a new empty triple store

<a name="TripleStore.Add"></a>
### func \(\*TripleStore\) [Add](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L68>)

```go
func (ts *TripleStore) Add(t Triple) bool
```

Add adds a triple to the store, returns true if it was new

<a name="TripleStore.All"></a>
### func \(\*TripleStore\) [All](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L142>)

```go
func (ts *TripleStore) All() []Triple
```

All returns all triples in the store

<a name="TripleStore.Contains"></a>
### func \(\*TripleStore\) [Contains](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L86>)

```go
func (ts *TripleStore) Contains(t Triple) bool
```

Contains checks if a triple exists in the store

<a name="TripleStore.FindByObject"></a>
### func \(\*TripleStore\) [FindByObject](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L109>)

```go
func (ts *TripleStore) FindByObject(object string) []Triple
```

FindByObject returns all triples with the given object

<a name="TripleStore.FindByPredicate"></a>
### func \(\*TripleStore\) [FindByPredicate](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L100>)

```go
func (ts *TripleStore) FindByPredicate(predicate string) []Triple
```

FindByPredicate returns all triples with the given predicate

<a name="TripleStore.FindByPredicateObject"></a>
### func \(\*TripleStore\) [FindByPredicateObject](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L130>)

```go
func (ts *TripleStore) FindByPredicateObject(predicate, object string) []Triple
```

FindByPredicateObject returns all triples matching predicate and object

<a name="TripleStore.FindBySubject"></a>
### func \(\*TripleStore\) [FindBySubject](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L91>)

```go
func (ts *TripleStore) FindBySubject(subject string) []Triple
```

FindBySubject returns all triples with the given subject

<a name="TripleStore.FindBySubjectPredicate"></a>
### func \(\*TripleStore\) [FindBySubjectPredicate](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L118>)

```go
func (ts *TripleStore) FindBySubjectPredicate(subject, predicate string) []Triple
```

FindBySubjectPredicate returns all triples matching subject and predicate

<a name="TripleStore.Size"></a>
### func \(\*TripleStore\) [Size](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/store.go#L149>)

```go
func (ts *TripleStore) Size() int
```

Size returns the number of triples in the store

<a name="TurtleParser"></a>
## type [TurtleParser](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/parser.go#L11-L16>)

TurtleParser parses Turtle format RDF

```go
type TurtleParser struct {
    // contains filtered or unexported fields
}
```

<a name="NewTurtleParser"></a>
### func [NewTurtleParser](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/parser.go#L19>)

```go
func NewTurtleParser() *TurtleParser
```

NewTurtleParser creates a new Turtle parser

<a name="TurtleParser.FallbackParse"></a>
### func \(\*TurtleParser\) [FallbackParse](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/parser.go#L527>)

```go
func (p *TurtleParser) FallbackParse(content string) ([]Triple, error)
```

Alternative fallback parser using regex \(for complex cases\)

<a name="TurtleParser.Parse"></a>
### func \(\*TurtleParser\) [Parse](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/parser.go#L27>)

```go
func (p *TurtleParser) Parse(content string) ([]Triple, error)
```

Parse parses Turtle content and returns triples

<a name="TypeInheritance"></a>
## type [TypeInheritance](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L61>)

TypeInheritance implements rdf:type inheritance through subClassOf If X rdf:type A and A rdfs:subClassOf B, then X rdf:type B

```go
type TypeInheritance struct{}
```

<a name="TypeInheritance.Apply"></a>
### func \(\*TypeInheritance\) [Apply](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L67>)

```go
func (r *TypeInheritance) Apply(store *TripleStore) []Triple
```



<a name="TypeInheritance.Name"></a>
### func \(\*TypeInheritance\) [Name](<https://github.com:beyondcivic/goreasoner/blob/main/pkg/reasoner/rules.go#L63>)

```go
func (r *TypeInheritance) Name() string
```

