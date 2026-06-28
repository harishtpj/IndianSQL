package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/harishtpj/indiansql/internal/consts"
)

func printBanner() {
	fmt.Printf("IndianSQL v%s\n", consts.Version)
	fmt.Println("Enter 'help' for more information.")
}

func printPrompt() {
	fmt.Print("indsql >>> ")
}

func getInput() (cmd string, args string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cmd, args, _ = strings.Cut(scanner.Text(), " ")
	return
}

func repl(dbFile string) {
	printBanner()

	db, err := NewREPLContext(dbFile)

	for {
		printPrompt()

		cmd, args := getInput()
		cmd = strings.ToLower(strings.TrimSpace(cmd))
		args = strings.TrimSpace(args)

		switch cmd {

		case "":
			continue

		case "help":
			fmt.Println("Commands:")
			fmt.Println("  open <database>")
			fmt.Println("  create <table> <col:type[:pk]> ...")
			fmt.Println("  insert <table> <values...>")
			fmt.Println("  select <table>")
			fmt.Println("  schema <table>")
			fmt.Println("  tables")
			fmt.Println("  exit")

		case "open":
			if db != nil {
				_ = db.Close()
			}

			db, err = NewREPLContext(args)
			if err == nil {
				fmt.Println("Database opened.")
			}

		case "create":
			if db == nil {
				fmt.Println("No database opened.")
				continue
			}

			fields := strings.Fields(args)
			if len(fields) < 2 {
				fmt.Println("Usage: create <table> <col:type[:pk]> ...")
				continue
			}

			tableName := fields[0]

			columns, parseErr := parseCreateColumns(fields[1:])
			if parseErr != nil {
				fmt.Println(parseErr)
				continue
			}

			err = db.ExecuteCreateTable(tableName, columns)
			if err == nil {
				fmt.Println("Table created.")
			}

		case "insert":
			if db == nil {
				fmt.Println("No database opened.")
				continue
			}

			fields := strings.Fields(args)
			if len(fields) < 2 {
				fmt.Println("Usage: insert <table> <values...>")
				continue
			}

			tableName := fields[0]

			table, err := db.GetTableInfo(tableName)
			if err != nil {
				break
			}

			values, parseErr := parseInsertValues(table, fields[1:])
			if parseErr != nil {
				fmt.Println(parseErr)
				continue
			}

			err = db.ExecuteInsert(tableName, values)
			if err == nil {
				fmt.Println("1 row inserted.")
			}

		case "select":
			if db == nil {
				fmt.Println("No database opened.")
				continue
			}

			rows, selectErr := db.ExecuteSelectAll(args)
			if selectErr != nil {
				err = selectErr
				break
			}

			table, _ := db.GetTableInfo(args)
			printRows(table, rows)

		case "schema":
			if db == nil {
				fmt.Println("No database opened.")
				continue
			}

			table, schemaErr := db.GetTableInfo(args)
			if schemaErr != nil {
				err = schemaErr
				break
			}

			fmt.Println(table)

		case "tables":
			if db == nil {
				fmt.Println("No database opened.")
				continue
			}

			if db.Catalog.TableCount() == 0 {
				fmt.Println("(no tables)")
				continue
			}

			for _, t := range db.Catalog.ListTables() {
				fmt.Println(t.Name)
			}

		case "exit", "quit":
			if db != nil {
				if err := db.Close(); err != nil {
					fmt.Fprintln(os.Stderr, err)
				}
			}
			fmt.Println("Bye.")
			return

		default:
			fmt.Printf("Unknown command: %q\n", cmd)
			fmt.Println("Type 'help' for available commands.")
		}

		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			err = nil
		}
	}
}
