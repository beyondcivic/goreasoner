// commands.go
// Contains cobra command definitions
//
//nolint:funlen,mnd
package cmd

import (
	"fmt"

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
