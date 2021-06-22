package main

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"strings"

	"github.com/spf13/cobra"

	"github.com/confluentinc/cli/internal/cmd"
	"github.com/confluentinc/cli/internal/pkg/doc"
	"github.com/confluentinc/cli/internal/pkg/version"
)

// Injected from linker flags like `go build -ldflags "-X main.cliName=$NAME"`
var cliName = "confluent"

// See https://github.com/spf13/cobra/blob/master/doc/rest_docs.md
func main() {
	// Prevent printing the user's HOME in docs when generating confluent local services kafka
	if err := os.Setenv("HOME", "$HOME"); err != nil {
		panic(err)
	}

	cli := cmd.NewConfluentCommand(cliName, true, &version.Version{})

	root := path.Join(".", "docs")

	if err := doc.GenReSTTree(cli.Command, root, doc.SphinxRef, 0); err != nil {
		panic(err)
	}

	// This overwrites the root index generated by the call to GenReSTTree above.
	if err := doc.GenReSTIndex(cli.Command, path.Join(root, cliName, "index.rst"), rootIndexHeader, doc.SphinxRef); err != nil {
		panic(err)
	}
}

func rootIndexHeader(_ *cobra.Command) string {
	buf := new(bytes.Buffer)

	buf.WriteString(fmt.Sprintf(".. _%s-ref:\n\n", cliName))
	title := fmt.Sprintf("|%s| CLI Command Reference\n", cliName)
	buf.WriteString(title)
	buf.WriteString(strings.Repeat("=", len(title)-1) + "\n\n")
	buf.WriteString(fmt.Sprintf("The available |%s| CLI commands are documented here.\n\n", cliName))

	return buf.String()
}
