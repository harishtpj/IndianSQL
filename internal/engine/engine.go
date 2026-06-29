package engine

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"github.com/harishtpj/indiansql/internal/manager"
	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
	"github.com/xwb1989/sqlparser"
)

type SQLEngine struct {
	db     *manager.DBManager
	dbFile string
}

func NewSQLEngine(dbFile string) (*SQLEngine, error) {
	db, err := manager.NewDBManager(dbFile)
	return &SQLEngine{db, dbFile}, err
}

func (engine *SQLEngine) Execute(query string) (Result, error) {
	query = strings.TrimSpace(query)
	query = strings.TrimSuffix(query, ";")

	cmd, args, _ := strings.Cut(query, " ")
	cmd = strings.ToLower(strings.TrimSpace(cmd))
	args = strings.TrimSpace(args)

	switch cmd {
	case "":
		return &EmptyResult{}, nil

	case "help":
		return &MessageResult{
			Message: `Commands:
  open <database>
  create table <table> (col type [primary key], ...)
  insert into <table> values (...)
  select * from <table>
  schema <table>
  tables
  exit`,
		}, nil

	case "open":
		if engine.db != nil {
			_ = engine.db.Close()
		}

		db, err := manager.NewDBManager(args)
		if err != nil {
			return nil, err
		}

		engine.db = db

		return &MessageResult{
			Message: "Database opened.",
		}, nil

	case "schema":
		table, err := engine.db.GetTableInfo(args)
		if err != nil {
			return nil, err
		}

		return &SchemaResult{
			Table: table,
		}, nil

	case "tables":
		return &TablesResult{
			Tables: engine.db.Catalog.ListTables(),
		}, nil

	case "desc", "describe":
		fields := strings.Fields(args)
		if len(fields) < 1 {
			return nil, errors.New("usage: desc <table>")
		}

		table, err := engine.db.GetTableInfo(fields[0])
		if err != nil {
			return nil, err
		}

		return &SchemaResult{
			Table: table,
		}, nil

	case "show":
		fields := strings.Fields(args)
		if len(fields) == 0 {
			return nil, errors.New("usage: show tables|columns from <table>")
		}

		switch strings.ToLower(fields[0]) {
		case "tables":
			return &TablesResult{
				Tables: engine.db.Catalog.ListTables(),
			}, nil

		case "columns", "fields":
			if len(fields) < 3 {
				return nil, errors.New("usage: show columns from <table>")
			}
			if !strings.EqualFold(fields[1], "from") && !strings.EqualFold(fields[1], "in") {
				return nil, errors.New("usage: show columns from <table>")
			}
			table, err := engine.db.GetTableInfo(fields[2])
			if err != nil {
				return nil, err
			}
			return &SchemaResult{Table: table}, nil
		default:
			return nil, fmt.Errorf("unknown show command: %q", fields[0])
		}

	case "exit", "quit":
		if engine.db != nil {
			_ = engine.db.Close()
		}

		return &ExitResult{}, nil
	}

	stmt, err := sqlparser.Parse(query)
	if err != nil {
		return nil, err
	}

	switch stmt := stmt.(type) {
	case *sqlparser.DDL:
		if stmt.Action == "create" {
			tableName := stmt.NewName.Name.String()
			columns := make([]*schema.Column, 0, len(stmt.TableSpec.Columns))

			for _, colDef := range stmt.TableSpec.Columns {
				col := &schema.Column{
					Name: colDef.Name.String(),
				}

				switch strings.ToLower(colDef.Type.Type) {
				case "int", "integer":
					col.Type = schema.ColumnTypeInteger
				case "varchar", "text", "string":
					col.Type = schema.ColumnTypeVarchar
				case "boolean", "bool":
					col.Type = schema.ColumnTypeBoolean
				case "numeric", "num", "float":
					col.Type = schema.ColumnTypeNumeric
				default:
					return nil, fmt.Errorf("unknown type %q", colDef.Type.Type)
				}

				if colDef.Type.KeyOpt == 1 {
					col.IsPrimaryKey = true
				}

				columns = append(columns, col)
			}

			if err := engine.db.ExecuteCreateTable(tableName, columns); err != nil {
				return nil, err
			}

			return &MessageResult{
				Message: "Table created.",
			}, nil
		}
		return nil, fmt.Errorf("unsupported DDL action")

	case *sqlparser.Insert:
		tableName := stmt.Table.Name.String()
		table, err := engine.db.GetTableInfo(tableName)
		if err != nil {
			return nil, err
		}

		rows, ok := stmt.Rows.(sqlparser.Values)
		if !ok || len(rows) != 1 {
			return nil, errors.New("insert must have exactly one row of values")
		}

		vals := rows[0]
		if len(vals) != table.ColumnCount() {
			return nil, fmt.Errorf(
				"expected %d values, got %d",
				table.ColumnCount(),
				len(vals),
			)
		}

		values := make([]row.Value, 0, len(vals))

		for i, valExpr := range vals {
			switch val := valExpr.(type) {
			case *sqlparser.SQLVal:
				strVal := string(val.Val)
				switch table.Columns[i].Type {
				case schema.ColumnTypeInteger:
					n, err := strconv.ParseInt(strVal, 10, 64)
					if err != nil {
						return nil, fmt.Errorf(
							"column %q expects INTEGER",
							table.Columns[i].Name,
						)
					}
					values = append(values, row.NewIntegerValue(n))
				case schema.ColumnTypeVarchar:
					values = append(values, row.NewVarcharValue(strVal))
				case schema.ColumnTypeNumeric:
					n, err := strconv.ParseFloat(strVal, 64)
					if err != nil {
						return nil, fmt.Errorf(
							"column %q expects NUMERIC",
							table.Columns[i].Name,
						)
					}
					values = append(values, row.NewNumericValue(n))
				case schema.ColumnTypeBoolean:
					values = append(values, row.NewBoolValue(strings.EqualFold(strVal, "TRUE") || strVal == "1"))
				default:
					return nil, fmt.Errorf(
						"unsupported datatype for column %q",
						table.Columns[i].Name,
					)
				}
			case sqlparser.BoolVal:
				if table.Columns[i].Type != schema.ColumnTypeBoolean {
					return nil, fmt.Errorf(
						"column %q expects BOOLEAN",
						table.Columns[i].Name,
					)
				}
				values = append(values, row.NewBoolValue(bool(val)))
			default:
				return nil, fmt.Errorf(
					"unsupported value expression for column %q",
					table.Columns[i].Name,
				)
			}
		}

		if err := engine.db.ExecuteInsert(tableName, values); err != nil {
			return nil, err
		}

		return &MessageResult{
			Message: "1 row inserted.",
		}, nil

	case *sqlparser.Select:
		if len(stmt.From) != 1 {
			return nil, errors.New("select must have exactly one table")
		}

		var tableName string
		if aliased, ok := stmt.From[0].(*sqlparser.AliasedTableExpr); ok {
			if table, ok := aliased.Expr.(sqlparser.TableName); ok {
				tableName = table.Name.String()
			}
		}

		if tableName == "" {
			return nil, errors.New("unsupported select from clause")
		}

		dbRows, err := engine.db.ExecuteSelectAll(tableName)
		if err != nil {
			return nil, err
		}

		table, err := engine.db.GetTableInfo(tableName)
		if err != nil {
			return nil, err
		}

		return &TableResult{
			Table: table,
			Rows:  dbRows,
		}, nil

	default:
		return nil, fmt.Errorf("unknown statement type")
	}
}

func (engine *SQLEngine) CommitDB() error {
	err := engine.db.Close()
	if err != nil {
		return err
	}
	engine.db, err = manager.NewDBManager(engine.dbFile)
	return err
}

func (engine *SQLEngine) GetDBName() string {
	return engine.dbFile
}

