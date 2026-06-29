package server

import (
	"fmt"
	"strings"

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

	case *engine.TablesResult:
		return buildTablesResult(r.Tables, h.engine.GetDBName())

	case *engine.SchemaResult:
		return buildSchemaResult(r.Table)

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

				case schema.ColumnTypeInteger:
					vals[i] = (*v).(*rowPackage.IntegerValue).GetInt64()

				case schema.ColumnTypeVarchar:
					vals[i] = (*v).String()

				case schema.ColumnTypeBoolean:
					vals[i] = (*v).(*rowPackage.BoolValue).GetInt()

				case schema.ColumnTypeNumeric:
					vals[i] = (*v).(*rowPackage.NumericValue).GetFloat()

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
		case schema.ColumnTypeInteger:
			field.Type = mysql.MYSQL_TYPE_LONG

		case schema.ColumnTypeVarchar:
			field.Type = mysql.MYSQL_TYPE_VARCHAR

		case schema.ColumnTypeBoolean:
			field.Type = mysql.MYSQL_TYPE_TINY

		case schema.ColumnTypeNumeric:
			field.Type = mysql.MYSQL_TYPE_DOUBLE

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
	return h.engine.CommitDB()
}

func (h *Handler) HandleOtherCommand(cmd byte, data []byte) error {
	return fmt.Errorf("unsupported command: %d", cmd)
}

func (h *Handler) UseDB(dbName string) error {
	//if dbName != "" {
	//	h.dbName = dbName
	//}
	return nil
}

func buildTablesResult(tables []*schema.Table, dbName string) (*mysql.Result, error) {
	columnName := "Tables"
	if dbName != "" {
		columnName = "Tables_in_" + dbName
	}

	values := make([][]any, 0, len(tables))
	for _, table := range tables {
		values = append(values, []any{table.Name})
	}

	rs, err := mysql.BuildSimpleResultset([]string{columnName}, values, false)
	if err != nil {
		return nil, err
	}

	return &mysql.Result{Resultset: rs}, nil
}

func buildSchemaResult(table *schema.Table) (*mysql.Result, error) {
	columns := []string{"Field", "Type", "Null", "Key", "Default", "Extra"}
	values := make([][]any, 0, len(table.Columns))

	for _, col := range table.Columns {
		nullability := "YES"
		key := ""
		if col.IsPrimaryKey {
			nullability = "NO"
			key = "PRI"
		}

		values = append(values, []any{
			col.Name,
			mysqlTypeName(col.Type),
			nullability,
			key,
			nil,
			"",
		})
	}

	rs, err := mysql.BuildSimpleResultset(columns, values, false)
	if err != nil {
		return nil, err
	}

	return &mysql.Result{Resultset: rs}, nil
}

func mysqlTypeName(colType schema.ColumnType) string {
	switch colType {
	case schema.ColumnTypeInteger:
		return "BIGINT"
	case schema.ColumnTypeVarchar:
		return "VARCHAR"
	case schema.ColumnTypeBoolean:
		return "BOOLEAN"
	case schema.ColumnTypeNumeric:
		return "NUMERIC"
	default:
		return strings.ToUpper(colType.String())
	}
}
