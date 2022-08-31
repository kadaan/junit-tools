package cmd

import (
	"github.com/kadaan/junit-tools/lib/command"
	"github.com/kadaan/junit-tools/version"
)

var (
	Root = command.NewRootCommand(
		"junit tools",
		version.Name+` provides tools to help integrating with junit`)
)
