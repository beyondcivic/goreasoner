// Package cmd provides the command-line interface for goreasoner.
//
// goreasoner is a Go implementation of a forward reasoner for RDF/OWL ontologies.
// It parses Turtle format inputs and applies RDFS/OWL inference rules to derive
// new facts from the given TBox (terminology/schema) and ABox (assertions/instances).
//
// The command-line tool provides functionality to:
//   - Parse and load Turtle format RDF data
//   - Apply forward reasoning with RDFS/OWL inference rules
//   - Query the resulting knowledge base
//   - Export inferred triples in N-Triples format
//   - Display reasoning statistics and version information
//
// # Command Reference
//
// Perform forward reasoning on RDF data:
//
//	goreasoner reason ontology.ttl data.ttl
//
// Load TBox and ABox separately:
//
//	goreasoner reason --tbox schema.ttl --abox instances.ttl
//
// Query the knowledge base:
//
//	goreasoner query ontology.ttl data.ttl --subject "ex:Person"
//
// Export all triples including inferred:
//
//	goreasoner reason ontology.ttl data.ttl --output results.nt
//
// Show reasoning statistics:
//
//	goreasoner reason ontology.ttl data.ttl --verbose
//
// Show version information:
//
//	goreasoner version
//
// # Features
//
// Semantic Reasoning:
//   - RDFS subclass transitivity (rdfs:subClassOf)
//   - Type inheritance through class hierarchies
//   - Domain and range inference (rdfs:domain, rdfs:range)
//   - Property hierarchy reasoning (rdfs:subPropertyOf)
//   - OWL class equivalence (owl:equivalentClass)
//   - Individual identity (owl:sameAs)
//   - Inverse properties (owl:inverseOf)
//   - Transitive and symmetric properties
//
// RDF Processing:
//   - Turtle format parsing with prefix support
//   - Triple store with efficient indexing
//   - N-Triples output format
//   - Query interface for knowledge base exploration
//
// Performance:
//
//   - Forward chaining inference engine
//
//   - Optimized rule application until fixpoint
//
//   - Memory-efficient triple storage
//
//   - Detailed reasoning statistics and progress reporting
//
//     goreasoner version
package cmd

import (
	"fmt"
	"os"

	"github.com/beyondcivic/goreasoner/pkg/version"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// Root cobra command.
// Call Init() once to initialize child commands.
// Global so it can be picked up by docs/doc-gen.go.
// nolint:gochecknoglobals
var RootCmd = &cobra.Command{
	Use:   "goreasoner",
	Short: "Forward reasoner for RDF/OWL ontologies",
	Long: `A Go implementation of a forward reasoner for RDF/OWL ontologies.
goreasoner parses Turtle format inputs and applies RDFS/OWL inference rules 
to derive new facts from TBox (terminology/schema) and ABox (assertions/instances).`,
	Version: version.Version,
}

// Call Once.
func Init() {
	// Initialize viper for configuration
	viper.SetEnvPrefix("GOREASONER")
	viper.AutomaticEnv()

	// Add child commands
	RootCmd.AddCommand(versionCmd())
	RootCmd.AddCommand(runCmd())
	RootCmd.AddCommand(dlQueryCmd())
}

func Execute() {
	// Execute the command
	if err := RootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
