package flags

import (
	"github.com/Meduzz/commando/flags"
	"github.com/Meduzz/commando/model"
	"github.com/spf13/cobra"
)

var ArrayKind model.FlagKind = "array"

func Array(name string) *model.Flag {
	return &model.Flag{
		Name: name,
		Kind: ArrayKind,
	}
}

func ArrayFlagVisitor() flags.FlagVisitor {
	return flags.NewVisitor(ArrayKind, func(f *model.Flag, c *cobra.Command) {
		zeArray, ok := f.Default.([]string)

		if !ok {
			zeArray = make([]string, 0)
		}

		c.Flags().StringSlice(f.Name, zeArray, f.Description)
	}, func(s string, c *cobra.Command) (any, error) {
		return c.Flags().GetStringSlice(s)
	})
}
