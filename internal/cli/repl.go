package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/harishtpj/havensql/internal/consts"
	"github.com/harishtpj/havensql/internal/pager"
)

func printBanner() {
	fmt.Printf("HavenSQL v%s\n", consts.Version)
	fmt.Println("Enter 'help' for more information.")
}

func printPrompt() {
	fmt.Print("hsdb >>> ")
}

func getInput() (cmd string, args string) {
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	cmd, args, _ = strings.Cut(scanner.Text(), " ")
	return
}

func repl() {
	printBanner()
	var pgr *pager.Pager
	var err error

	for {
		printPrompt()
		cmd, args := getInput()

		switch cmd {
		case "open":
			pgr, err = pager.Open(args)
		case "write":
			if pgr == nil {
				fmt.Println("No database opened.")
				continue
			}
			n, str, _ := strings.Cut(args, " ")
			nPg, _ := strconv.Atoi(n)
			pg, _ := pgr.GetPage(uint32(nPg))
			copy(pg.Data, []byte(str))
			pg.Dirty = true
		case "read":
			nPg, _ := strconv.Atoi(args)
			pg, _ := pgr.GetPage(uint32(nPg))
			fmt.Println(string(pg.Data[:20]))
		case "exit":
			pgr.Close()
			return
		default:
			fmt.Fprintf(os.Stderr, "Unknown command: '%s'\n", cmd)
		}

		if err != nil {
			panic(err)
		}
	}
}
