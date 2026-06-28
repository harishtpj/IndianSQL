package server

import (
	"fmt"

	"github.com/go-mysql-org/go-mysql/mysql"

	"github.com/harishtpj/indiansql/internal/engine"
	rowPackage "github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
)

type Handler struct {
	engine *engine.SQLEngine
}

func NewHandler(engine *engine.SQLEngine) *Handler {
	return &Handler{
		engine: engine,
	}
}

func (h *Handler) HandleQuery(query string) (*mysql.Result, error) {
	result, err := h.engine.Execute(query)
	if err != nil {
		return nil, err
	}

	switch r := result.(type) {

	case *engine.EmptyResult:
		return &mysql.Result{}, nil

	case *engine.MessageResult:
		// CREATE, INSERT, OPEN, etc.
		return &mysql.Result{}, nil

	case *engine.TableResult:
		columns := make([]string, len(r.Table.Columns))
		for i, col := range r.Table.Columns {
			columns[i] = col.Name
		}

		values := make([][]any, 0, len(r.Rows))

		for _, row := range r.Rows {
			vals := make([]any, len(row.Values))

			for i, v := range row.Values {
				switch (*v).Type() {

				case schema.ColumnTypeNumeric:
					vals[i] = (*v).(*rowPackage.NumericValue).GetInt64()

				case schema.ColumnTypeVarchar:
					vals[i] = (*v).String()

				default:
					vals[i] = (*v).String()
				}
			}

			values = append(values, vals)
		}

		rs, err := mysql.BuildSimpleResultset(columns, values, false)
		if err != nil {
			return nil, err
		}

		return &mysql.Result{
			Resultset: rs,
		}, nil

	case *engine.SchemaResult:
		return nil, fmt.Errorf("SCHEMA is not supported over the MySQL protocol")

	case *engine.TablesResult:
		return nil, fmt.Errorf("TABLES is not supported over the MySQL protocol")

	case *engine.ExitResult:
		return &mysql.Result{}, nil

	default:
		return nil, fmt.Errorf("unknown engine result %T", result)
	}
}

func (h *Handler) HandleFieldList(table string, fieldWildcard string) ([]*mysql.Field, error) {
	res, err := h.engine.Execute("schema " + table)
	if err != nil {
		return nil, err
	}

	schemaResult, ok := res.(*engine.SchemaResult)
	if !ok {
		return nil, fmt.Errorf("expected SchemaResult")
	}

	fields := make([]*mysql.Field, len(schemaResult.Table.Columns))

	for i, col := range schemaResult.Table.Columns {
		field := &mysql.Field{
			Name: []byte(col.Name),
		}

		switch col.Type {
		case schema.ColumnTypeNumeric:
			field.Type = mysql.MYSQL_TYPE_LONG

		case schema.ColumnTypeVarchar:
			field.Type = mysql.MYSQL_TYPE_VARCHAR

		default:
			field.Type = mysql.MYSQL_TYPE_VAR_STRING
		}

		fields[i] = field
	}

	return fields, nil
}

func (h *Handler) HandleStmtPrepare(query string) (params int, columns int, context any, err error) {
	return 0, 0, nil, fmt.Errorf("prepared statements are not supported")
}

func (h *Handler) HandleStmtExecute(context any, query string, args []any) (*mysql.Result, error) {
	return nil, fmt.Errorf("prepared statements are not supported")
}

func (h *Handler) HandleStmtClose(context any) error {
	return nil
}

func (h *Handler) HandleOtherCommand(cmd byte, data []byte) error {
	return fmt.Errorf("unsupported command: %d", cmd)
}

func (h *Handler) UseDB(dbName string) error {
	return nil
}
