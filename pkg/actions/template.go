package actions

import "github.com/spf13/cobra"

func AsTemplate(template string, model any) func(*cobra.Command, []string) error {
	return nil
}
