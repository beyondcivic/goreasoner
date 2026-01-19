// commands.go
// Contains cobra command definitions
//
//nolint:funlen,mnd
package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/beyondcivic/goreasoner/pkg/reasoner"
	"github.com/beyondcivic/goreasoner/pkg/version"
	"github.com/spf13/cobra"
)

// Version Command.
// Displays tool version and build information.
func versionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		Long:  `Print the version, git hash, and build time information of the goreasoner tool.`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("%s version %s\n", version.AppName, version.Version)
			stamp := version.RetrieveStamp()
			fmt.Printf("  Built with %s on %s\n", stamp.InfoGoCompiler, stamp.InfoBuildTime)
			fmt.Printf("  Git ref: %s\n", stamp.VCSRevision)
			fmt.Printf("  Go version %s, GOOS %s, GOARCH %s\n", stamp.InfoGoVersion, stamp.InfoGOOS, stamp.InfoGOARCH)
		},
	}
}

// Run command
func runCmd() *cobra.Command {
	var runCmd = &cobra.Command{
		Use:   "run [aboxPath] [tboxPath]",
		Short: "Run forward reasoning on RDF data",
		Long:  `Run forward reasoning on RDF data, applying RDFS/OWL inference rules to derive new facts from TBox and ABox.`,
		Args:  cobra.ExactArgs(2),
		Run: func(cmd *cobra.Command, args []string) {
			aboxPath := args[0]
			tboxPath := args[1]
			flagOutputPath, _ := cmd.Flags().GetString("output")
			flagOutputType, _ := cmd.Flags().GetString("outputType")

			// Validate input files
			if !fileExists(aboxPath) {
				fmt.Printf("Error: ABox file '%s' does not exist.\n", aboxPath)
				os.Exit(1)
			}

			if !fileExists(tboxPath) {
				fmt.Printf("Error: TBox file '%s' does not exist.\n", tboxPath)
				os.Exit(1)
			}

			if !isTurtleFile(aboxPath) {
				fmt.Printf("Error: File '%s' does not appear to be a Turtle file.\n", aboxPath)
				os.Exit(1)
			}

			if !isTurtleFile(tboxPath) {
				fmt.Printf("Error: File '%s' does not appear to be a Turtle file.\n", tboxPath)
				os.Exit(1)
			}

			// Determine output path
			outputPath := determineOutputPath(flagOutputPath, aboxPath)

			// Validate output type
			if flagOutputType != "ntriple" && flagOutputType != "datalog" {
				fmt.Printf("Error: Invalid output type '%s'. Must be 'ntriple' or 'datalog'.\n", flagOutputType)
				os.Exit(1)
			}

			// Read input files
			aboxContent, err := readFile(aboxPath)
			if err != nil {
				fmt.Printf("Error reading ABox file: %v\n", err)
				os.Exit(1)
			}

			tboxContent, err := readFile(tboxPath)
			if err != nil {
				fmt.Printf("Error reading TBox file: %v\n", err)
				os.Exit(1)
			}

			// Run forward reasoning
			fmt.Printf("Running forward reasoning on '%s' and '%s'...\n", aboxPath, tboxPath)
			inferredTriples, err := reasoner.ForwardReason(aboxContent, tboxContent)
			if err != nil {
				fmt.Printf("Error running forward reasoning: %v\n", err)
				os.Exit(1)
			}

			// Convert output format if needed
			var outputTriples []string
			if flagOutputType == "datalog" {
				outputTriples = reasoner.ConvertTriplesToDatalog(inferredTriples)
			} else {
				outputTriples = inferredTriples
			}

			// Write results to output file
			if outputPath != "" {
				err = writeTriplesToFile(outputTriples, outputPath)
				if err != nil {
					fmt.Printf("Error writing output file: %v\n", err)
					os.Exit(1)
				}
				fmt.Printf("✓ Forward reasoning completed successfully and saved to: %s\n", outputPath)
				fmt.Printf("  Total triples: %d (format: %s)\n", len(outputTriples), flagOutputType)
			} else {
				// Print to stdout if no output file specified
				for _, triple := range outputTriples {
					fmt.Println(triple)
				}
			}
		},
	}
	runCmd.Flags().StringP("output", "o", "", "Output path for the N-Triples file")
	runCmd.Flags().String("outputType", "ntriple", "Output format: 'ntriple' or 'datalog' (default: ntriple)")

	return runCmd
}

// Helper function to check if file exists
func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}

// Helper function to check if file is a Turtle file
func isTurtleFile(filename string) bool {
	ext := strings.ToLower(filename[strings.LastIndex(filename, ".")+1:])
	return ext == "ttl" || ext == "turtle" || ext == "n3"
}

// Helper function to read file contents
func readFile(filename string) (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	defer file.Close()

	content, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(content), nil
}

// Helper function to write triples to file
func writeTriplesToFile(triples []string, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	for _, triple := range triples {
		_, err = file.WriteString(triple + "\n")
		if err != nil {
			return err
		}
	}

	return nil
}

// Helper function to determine output path
func determineOutputPath(providedPath, inputPath string) string {
	if providedPath != "" {
		return providedPath
	}

	// Check environment variable
	envOutputPath := os.Getenv("GOREASONER_OUTPUT_PATH")
	if envOutputPath != "" {
		return envOutputPath
	}

	// Generate default path based on input filename
	baseName := strings.TrimSuffix(inputPath, ".ttl")
	baseName = strings.TrimSuffix(baseName, ".turtle")
	baseName = strings.TrimSuffix(baseName, ".n3")
	return baseName + "_inferred.nt"
}
