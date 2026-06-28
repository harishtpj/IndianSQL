package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/harishtpj/indiansql/internal/consts"
)

var rootCmd = &cobra.Command{
	Use:   consts.AppName + " [dbfile]",
	Short: consts.Desc,
	Long:  `A minimal file-based DBMS engine with built-in server support.`,
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dbFile := ":memory:"
		if len(args) == 1 {
			dbFile = args[0]
		}

		repl(dbFile)
	},
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.InitDefaultHelpCmd()
	rootCmd.Flags().MarkHidden("help")
	rootCmd.SetUsageTemplate(`Usage:
{{.UseLine}} [arguments]

Available Commands:{{range .Commands}}{{if or .IsAvailableCommand (eq .Name "help")}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}

Use "{{.CommandPath}} help <command>" for more information about a command.
`)
}
