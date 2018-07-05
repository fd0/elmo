package clone

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"github.com/fd0/elmo/gitolite"
	"github.com/fd0/elmo/options"
	"github.com/spf13/cobra"
)

var gopts *options.Options
var opts Options

// Options collect configuration for a single run of the command.
type Options struct {
}

// AddCommand adds the command to c.
func AddCommand(c *cobra.Command, globalOptions *options.Options) {
	c.AddCommand(cmd)
	gopts = globalOptions
}

var cmd = &cobra.Command{
	Use: "clone [options] [pattern]",
	DisableFlagsInUseLine: true,

	Short:   helpShort,
	Long:    helpLong,
	Example: helpExamples,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 {
			return errors.New("no repository to clone")
		}

		if gopts.Target == "" {
			return errors.New("target directory is empty, please pass --target (or set $ELMO_TARGET)")
		}

		pattern := args[0]
		if err := gitolite.ValidatePattern(pattern); err != nil {
			return fmt.Errorf("pattern %q is invalid: %v", pattern, err)
		}

		srv := gitolite.Server{
			Hostname: gopts.Server,
		}

		info, err := srv.Info()
		if err != nil {
			return err
		}

		for name := range gitolite.Filter(info.Repos, pattern) {
			target := filepath.Join(gopts.Target, name)
			_, err := os.Lstat(target)
			if err == nil {
				fmt.Printf("skip %v: already cloned\n", name)
				continue
			}

			fmt.Printf("clone %v\n", name)

			err = srv.Clone(name, target)
			if err != nil {
				fmt.Fprintf(os.Stderr, "error cloning %v: %v\n", name, err)
			}
		}

		return nil
	},
}
