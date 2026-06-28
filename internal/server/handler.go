package server

import (
	"fmt"
	"strings"

	"github.com/go-mysql-org/go-mysql/mysql"
	"github.com/harishtpj/indiansql/internal/handler"
	rowPackage "github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
)

type Handler struct {
	db *handler.REPLContext
}

// HandleFieldList implements [server.Handler].
func (h Handler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	panic("unimplemented")
}

// HandleOtherCommand implements [server.Handler].
func (h Handler) HandleOtherCommand(cmd byte, data []byte) error {
	panic("unimplemented")
}

// HandleStmtClose implements [server.Handler].
func (h Handler) HandleStmtClose(context any) error {
	panic("unimplemented")
}

// HandleStmtExecute implements [server.Handler].
func (h Handler) HandleStmtExecute(context any, query string, args []any) (*mysql.Result, error) {
	panic("unimplemented")
}

// HandleStmtPrepare implements [server.Handler].
func (h Handler) HandleStmtPrepare(query string) (params int, columns int, context any, err error) {
	panic("unimplemented")
}

// UseDB implements [server.Handler].
func (h Handler) UseDB(dbName string) error {
	panic("unimplemented")
}

func (h *Handler) HandleQuery(query string) (*mysql.Result, error) {
	cmd, args, _ := strings.Cut(query, " ")
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	args = strings.TrimSpace(args)

	switch cmd {
	case "select":
		rows, err := h.db.ExecuteSelectAll(args)
		if err != nil {
			return nil, err
		}

		table, err := h.db.GetTableInfo(args)
		if err != nil {
			return nil, err
		}

		// Column names
		columns := make([]string, len(table.Columns))
		for i, col := range table.Columns {
			columns[i] = col.Name
		}

		// Row values
		values := make([][]any, 0, len(rows))
		for _, r := range rows {
			row := make([]any, len(r.Values))

			for i, v := range r.Values {
				switch (*v).Type() {
				case schema.ColumnTypeNumeric:
					row[i] = (*v).(*rowPackage.NumericValue).GetInt64()

				case schema.ColumnTypeVarchar:
					row[i] = (*v).String()

				default:
					row[i] = (*v).String()
				}
			}

			values = append(values, row)
		}

		rs, err := mysql.BuildSimpleResultset(columns, values, false)
		if err != nil {
			return nil, err
		}

		return &mysql.Result{
			Resultset: rs,
		}, nil

	default:
		return nil, fmt.Errorf("unknown command: %q", cmd)
	}
}
