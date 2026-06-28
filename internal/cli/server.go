package cli

import (
	"github.com/harishtpj/indiansql/internal/server"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "serve",
	Short: "start the IndianSQL server",
	Long:  "starts the MySQL-compatiable server for IndianSQL",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbFile := ":memory:"
		if len(args) == 1 {
			dbFile = args[0]
		}
		s := server.Server{Addr: "0.0.0.0:3306"}
		s.Serve(dbFile)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
