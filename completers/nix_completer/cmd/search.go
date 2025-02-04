package cmd

import (
	"github.com/rsteube/carapace"
	"github.com/spf13/cobra"
)

var searchCmd = &cobra.Command{
	Use:     "search",
	Short:   "search for packages",
	GroupID: "main",
	Run:     func(cmd *cobra.Command, args []string) {},
}

func init() {
	carapace.Gen(searchCmd).Standalone()

	searchCmd.Flags().BoolP("exclude", "e", false, "Hide packages whose attribute path, name or description contain regex")
	searchCmd.Flags().Bool("json", false, "Produce output in JSON format, suitable for consumption by another program")
	rootCmd.AddCommand(searchCmd)

	addEvaluationFlags(searchCmd)
	addFlakeFlags(searchCmd)
	addLoggingFlags(searchCmd)
}
