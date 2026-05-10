package cli

import (
	"fmt"

	"github.com/harishtpj/havensql/internal/consts"
	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Prints the version of HavenSQL",
	Long:  "Prints the version info of HavenSQL",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Printf("%s v%s\n", consts.AppName, consts.Version)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
