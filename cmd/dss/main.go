// dss is the Design System Spec CLI tool.
// It generates code artifacts and validates implementations against DSS specifications.
package main

import (
	"os"

	"github.com/plexusone/design-system-spec/cmd/dss/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}
