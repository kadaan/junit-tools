package cmd

import (
	osErrors "errors"
	"github.com/kadaan/junit-tools/config"
	"github.com/kadaan/junit-tools/lib/command"
	"github.com/kadaan/junit-tools/lib/errors"
	"github.com/kadaan/junit-tools/lib/verifier"
	"github.com/spf13/cobra"
	"io/fs"
	"os"
)

func init() {
	command.NewCommand(
		Root,
		"verify",
		"Verify junit results were successful",
		"Verify the specified junit result XML files and ensure that all test cases were successful.",
		new(config.VerifyConfig),
		verifier.NewVerifier()).Configure(func(cb command.CommandBuilder, fb config.FlagBuilder, cfg *config.VerifyConfig) {
		cb.Args(func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("no junit result files were specified")
			}
			var errs []error
			for _, arg := range args {
				if _, err := os.Stat(arg); osErrors.Is(err, fs.ErrNotExist) {
					errs = append(errs, err)
				}
			}
			if len(errs) > 0 {
				return errors.NewMulti(errs, "one or more provided junit result files do not exist")
			}
			return nil
		})
	})
}
