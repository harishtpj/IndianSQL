package cli

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/harishtpj/havensql/internal/consts"
)

var rootCmd = &cobra.Command{
	Use:   consts.AppName,
	Short: consts.Desc,
	Long:  `A minimal file-based DBMS engine with built-in server support.`,
	Run: func(cmd *cobra.Command, args []string) {
		repl()
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
