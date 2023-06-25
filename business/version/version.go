/*
Copyright athenasoft Corp. All Rights Reserved.

*/

package version

import (
	"fmt"
	"github.com/spf13/cobra"
	"runtime"
	"time"
)

var ProgramName = "business"

var Version string = "0.0.2"
var CommitSHA string = ""
var BaseVersion string = "0.0.1"
var BaseDockerLabel string = ProgramName
var DockerNamespace string = ProgramName
var BaseDockerNamespace string = ProgramName
var BuildDate string = time.Now().String()

// Program name

func Cmd() *cobra.Command {
	return cobraCommand
}

var cobraCommand = &cobra.Command{
	Use:   "version",
	Short: "Print version.",
	Long:  `Print current version.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 0 {
			return fmt.Errorf("trailing args detected")
		}
		// Parsing of the command line is done so silence cmd usage
		cmd.SilenceUsage = true
		fmt.Print(GetInfo())
		return nil
	},
}

// GetInfo returns version information for the peer
func GetInfo() string {
	if Version == "" {
		Version = "development build"
	}

	if CommitSHA == "" {
		CommitSHA = "development build"
	}

	return fmt.Sprintf("%s:\n Version: %s\n Commit SHA: %s\n Go version: %s\n"+
		" OS/Arch: %s\n BuildDate: %s\n",
		ProgramName,
		Version,
		CommitSHA,
		runtime.Version(),
		fmt.Sprintf("%s/%s\n", runtime.GOOS, runtime.GOARCH),
		BuildDate)
}
