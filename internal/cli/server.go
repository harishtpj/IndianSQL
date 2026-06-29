package cli

import (
	"fmt"
	"os"

	"github.com/harishtpj/indiansql/internal/server"
	"github.com/spf13/cobra"
	"golang.org/x/term"
	"gopkg.in/yaml.v3"
)

var (
	host       string
	port       int
	username   string
	askPass    bool
	configFile string
)

var serverCmd = &cobra.Command{
	Use:   "server [database]",
	Short: "Start the IndianSQL server",
	Long:  "Starts the MySQL-compatiable server for IndianSQL",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		s := server.Server{
			DBFile:   ":memory:",
			Host:     host,
			Port:     port,
			Username: username,
			Password: "",
		}

		if configFile != "" {
			data, err := os.ReadFile(configFile)
			if err != nil {
				return err
			}

			if err := yaml.Unmarshal(data, &s); err != nil {
				return err
			}
		}

		if len(args) == 1 {
			if args[0] != "" {
				s.DBFile = args[0]
			}
		}

		if askPass {
			fmt.Print("Enter password for server: ")

			pswd, err := term.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				return err
			}

			fmt.Println()
			s.Password = string(pswd)
		}

		return s.Serve()
	},
}

func init() {
	serverCmd.Flags().StringVarP(
		&host,
		"host",
		"H",
		"127.0.0.1",
		"Server host",
	)

	serverCmd.Flags().IntVarP(
		&port,
		"port",
		"P",
		4405,
		"Server port",
	)

	serverCmd.Flags().StringVarP(
		&username,
		"username",
		"u",
		"root",
		"Server username",
	)

	serverCmd.Flags().BoolVarP(
		&askPass,
		"password",
		"p",
		false,
		"Prompt for server password",
	)

	serverCmd.Flags().StringVarP(
		&configFile,
		"config",
		"c",
		"",
		"Load server configuration from YAML file",
	)

	rootCmd.AddCommand(serverCmd)
}
