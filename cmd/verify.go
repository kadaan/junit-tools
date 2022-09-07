package cmd

import (
	"github.com/kadaan/junit-tools/config"
	"github.com/kadaan/junit-tools/lib/command"
	"github.com/kadaan/junit-tools/lib/verifier"
)

func init() {
	command.NewCommand(
		Root,
		"verify",
		"Verify junit results were successful",
		"Verify the specified junit result XML files and ensure that all test cases were successful.",
		new(config.VerifyConfig),
		verifier.NewVerifier()).Configure(func(cb command.CommandBuilder, fb config.FlagBuilder, cfg *config.VerifyConfig) {})
}
