package list

import (
	"fmt"

	"github.com/fd0/elmo/gitolite"
	"github.com/fd0/elmo/options"
	"github.com/spf13/cobra"
)

var gopts *options.Options
var opts Options

// Options collect configuration for a single run of the command.
type Options struct {
	OnlyReadOnly  bool
	OnlyReadWrite bool
}

// AddCommand adds the command to c.
func AddCommand(c *cobra.Command, globalOptions *options.Options) {
	c.AddCommand(cmd)
	gopts = globalOptions

	fs := cmd.Flags()
	fs.BoolVar(&opts.OnlyReadOnly, "read-only", false, "only display repositories with read-only access")
	fs.BoolVar(&opts.OnlyReadWrite, "read-write-only", false, "only display repositories with read-write access")
}

var cmd = &cobra.Command{
	Use: "list [options] [pattern]",
	DisableFlagsInUseLine: true,

	Short:   helpShort,
	Long:    helpLong,
	Example: helpExamples,

	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) > 1 {
			return fmt.Errorf("more than one argument passed to 'list': %v", args)
		}

		var pattern string
		if len(args) == 1 {
			pattern = args[0]

			if err := gitolite.ValidatePattern(pattern); err != nil {
				return err
			}
		}

		srv := gitolite.Server{
			Hostname: gopts.Server,
		}

		info, err := srv.Info()
		if err != nil {
			return err
		}

		repos := gitolite.Filter(info.Repos, pattern)
		for _, name := range repos.Names() {
			repo := info.Repos[name]
			fmt.Printf(" % -6s %s\n", repo.Perms, name)
		}

		return nil
	},
}
