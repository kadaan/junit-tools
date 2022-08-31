package config

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

func NewFlagBuilder(cmd *cobra.Command) FlagBuilder {
	return &flagBuilder{
		cmd: cmd,
	}
}

type Flag interface {
	Required() Flag
}

type FileFlag interface {
	Flag
	Extensions(extensions ...string) FileFlag
}

type flag struct {
	builder *flagBuilder
	flag    *pflag.Flag
}

func (f *flag) Required() Flag {
	_ = f.builder.cmd.MarkFlagRequired(f.flag.Name)
	return f
}

func (f *flag) Extensions(extensions ...string) FileFlag {
	_ = f.builder.cmd.MarkFlagFilename(f.flag.Name, extensions...)
	return f
}

type FlagBuilder interface {
}

type flagBuilder struct {
	cmd *cobra.Command
}

func (fb *flagBuilder) newFlag(name string, creator func(flagSet *pflag.FlagSet)) *flag {
	creator(fb.cmd.Flags())
	f := fb.cmd.Flags().Lookup(name)
	_ = viper.BindPFlag(name, f)
	return &flag{
		builder: fb,
		flag:    f,
	}
}

func (fb *flagBuilder) addValidation(validation func(cmd *cobra.Command, args []string) error) {
	if fb.cmd.PreRunE != nil {
		existingValidation := fb.cmd.PreRunE
		fb.cmd.PreRunE = func(cmd *cobra.Command, args []string) error {
			if err := validation(cmd, args); err != nil {
				return err
			}
			return existingValidation(cmd, args)
		}
	} else {
		fb.cmd.PreRunE = validation
	}
}
