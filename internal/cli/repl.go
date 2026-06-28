package cli

import (
	"bufio"
	"fmt"
	"os"

	"github.com/harishtpj/indiansql/internal/consts"
	"github.com/harishtpj/indiansql/internal/engine"
)

func printBanner() {
	fmt.Printf("IndianSQL v%s\n", consts.Version)
	fmt.Println("Enter 'help' for more information.")
}

func printPrompt() {
	fmt.Print("indsql >>> ")
}

func getInput() string {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	return scanner.Text()
}

func repl(dbFile string) {
	printBanner()

	sqlEngine, err := engine.NewSQLEngine(dbFile)
	if err != nil {
		fmt.Println(err)
	}

	for {
		printPrompt()

		result, err := sqlEngine.Execute(getInput())
		if err != nil {
			fmt.Println(err)
			continue
		}

		switch r := result.(type) {

		case *engine.EmptyResult:

		case *engine.MessageResult:
			fmt.Println(r.Message)

		case *engine.TableResult:
			engine.PrintRows(r.Table, r.Rows)

		case *engine.SchemaResult:
			fmt.Println(r.Table)

		case *engine.TablesResult:
			if len(r.Tables) == 0 {
				fmt.Println("(no tables)")
				continue
			}

			for _, t := range r.Tables {
				fmt.Println(t.Name)
			}

		case *engine.ExitResult:
			fmt.Println("Bye.")
			return
		}
	}
}
