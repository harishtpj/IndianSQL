package engine

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/harishtpj/indiansql/internal/manager"
	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
	"github.com/xwb1989/sqlparser"
)

type SQLEngine struct {
	db *manager.DBManager
}

func NewSQLEngine(dbFile string) (*SQLEngine, error) {
	db, err := manager.NewDBManager(dbFile)
	return &SQLEngine{db}, err
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
  use <database> /* in current directory */
  create table <table> (col type [primary key], ...)
  insert into <table> values (...)
  select * from <table>
  desc[ribe] <table>
  show tables
  commit
  exit/quit`,
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
				case "boolean", "bool", "tinyint", "bit":
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

		var predicate manager.Predicate
		if stmt.Where != nil {
			predicate, err = CompilePredicate(stmt.Where.Expr)
			if err != nil {
				return nil, fmt.Errorf("failed compilation of where predicate: %v", err)
			}
		}

		// TODO: Assuming that user could provide only one StarExpr
		// In normal cases, user could also give table.*, which is NOT handled
		var selectCols []string
		for _, expr := range stmt.SelectExprs {
			switch expr := expr.(type) {
			case *sqlparser.StarExpr:
				break
			case *sqlparser.AliasedExpr:
				switch col := expr.Expr.(type) {

				case *sqlparser.ColName:
					selectCols = append(selectCols, col.Name.String())

				default:
					return nil, fmt.Errorf("unsupported select expression %T", col)
				}
			default:
				return nil, fmt.Errorf("invalid select expression: %T", expr)
			}
		}

		dbRows, err := engine.db.ExecuteSelect(tableName, selectCols, predicate)
		if err != nil {
			return nil, err
		}

		table, err := engine.db.GetTableInfo(tableName)
		if err != nil {
			return nil, err
		}

		if selectCols != nil {
			table = table.Filter(selectCols)
		}

		return &TableResult{
			Table: table,
			Rows:  dbRows,
		}, nil

	case *sqlparser.Show:
		switch stmt.Type {
		case "tables":
			return &TablesResult{
				Tables: engine.db.Catalog.ListTables(),
			}, nil

		default:
			return nil, fmt.Errorf("unknown show command: %q", stmt.Type)
		}

	case *sqlparser.Use:
		if engine.db != nil {
			_ = engine.db.Close()
		}

		db, err := manager.NewDBManager(stmt.DBName.String() + ".idb")
		if err != nil {
			return nil, err
		}

		engine.db = db

		return &MessageResult{
			Message: "Database opened.",
		}, nil

	case *sqlparser.Commit:
		engine.CommitDB()
		return &MessageResult{
			Message: "DB committed",
		}, nil

	default:
		return nil, fmt.Errorf("unknown statement type")
	}
}

func (engine *SQLEngine) CommitDB() error {
	return engine.db.Commit()
}

func (engine *SQLEngine) Close() error {
	return engine.db.Close()
}

func (engine *SQLEngine) GetDBName() string {
	base := filepath.Base(engine.db.GetDBName())
	return strings.TrimSuffix(base, filepath.Ext(base))
}
