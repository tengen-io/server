package cmd

import (
	"github.com/spf13/cobra"
	"github.com/tengen-io/server/server"
)

func init() {
	rootCmd.AddCommand(serverCmd)
}

var serverCmd = &cobra.Command{
	Use: "server",
	Short: "runs the tengen.io server",
	Run: func(cmd *cobra.Command, args []string) {
		serve()
	},
}

func serve() {
	server.Serve()
}
