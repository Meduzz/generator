package flags

import (
	"github.com/Meduzz/commando/flags"
	"github.com/Meduzz/commando/model"
	"github.com/spf13/cobra"
)

var MapKind model.FlagKind = "map"

func MapFlag(name string) *model.Flag {
	return &model.Flag{
		Kind: MapKind,
		Name: name,
	}
}

func MapFlagVisitor() flags.FlagVisitor {
	return flags.NewVisitor(MapKind, func(f *model.Flag, c *cobra.Command) {
		zeMap, ok := f.Default.(map[string]string)

		if !ok {
			zeMap = make(map[string]string)
		}

		c.Flags().StringToString(f.Name, zeMap, f.Description)
	}, func(s string, c *cobra.Command) (any, error) {
		return c.Flags().GetStringToString(s)
	})
}
