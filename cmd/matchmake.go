package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tengen-io/server/matchmake"
)

func init() {
	rootCmd.AddCommand(matchmakeCmd)
}

var matchmakeCmd = &cobra.Command{
	Use:   "matchmake",
	Short: "runs the tengen.io matchmaker",
	Run: func(cmd *cobra.Command, args []string) {
		runMatchmake()
	},
}

func runMatchmake() {
	matchmake.Start()
}
