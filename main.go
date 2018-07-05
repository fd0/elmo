package main

import (
	"fmt"
	"os"

	"github.com/fd0/elmo/cmd/clone"
	"github.com/fd0/elmo/cmd/list"
	"github.com/fd0/elmo/options"
	"github.com/spf13/cobra"
)

var cmdRoot = &cobra.Command{
	Use:           "elmo COMMAND [options]",
	SilenceErrors: true,
	SilenceUsage:  true,
}

var opts options.Options

func init() {
	// configure cobra help texts
	setupHelp(cmdRoot)
	// add version command
	cmdRoot.AddCommand(cmdVersion)

	// add the other commands
	list.AddCommand(cmdRoot, &opts)
	clone.AddCommand(cmdRoot, &opts)

	fs := cmdRoot.PersistentFlags()
	fs.StringVarP(&opts.Server, "server", "s", os.Getenv("ELMO_SERVER"), "gitolite server name (default: $ELMO_SERVER)")
	fs.StringVarP(&opts.Target, "target", "t", os.Getenv("ELMO_TARGET"), "local directory where repos are checked out (default: $ELMO_TARGET)")
}

func main() {
	err := cmdRoot.Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
