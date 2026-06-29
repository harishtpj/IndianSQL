package engine

import (
	"errors"
	"fmt"
	"strings"

	"github.com/harishtpj/indiansql/internal/manager"
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
  create <table> <col:type[:pk]> ...
  insert <table> <values...>
  select <table>
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

	case "create":
		fields := strings.Fields(args)
		if len(fields) < 2 {
			return nil, errors.New("usage: create <table> <col:type[:pk]> ...")
		}

		tableName := fields[0]

		columns, err := parseCreateColumns(fields[1:])
		if err != nil {
			return nil, err
		}

		if err := engine.db.ExecuteCreateTable(tableName, columns); err != nil {
			return nil, err
		}

		return &MessageResult{
			Message: "Table created.",
		}, nil

	case "insert":
		fields := strings.Fields(args)
		if len(fields) < 2 {
			return nil, errors.New("usage: insert <table> <values...>")
		}

		tableName := fields[0]

		table, err := engine.db.GetTableInfo(tableName)
		if err != nil {
			return nil, err
		}

		values, err := parseInsertValues(table, fields[1:])
		if err != nil {
			return nil, err
		}

		if err := engine.db.ExecuteInsert(tableName, values); err != nil {
			return nil, err
		}

		return &MessageResult{
			Message: "1 row inserted.",
		}, nil

	case "select":
		rows, err := engine.db.ExecuteSelectAll(args)
		if err != nil {
			return nil, err
		}

		table, err := engine.db.GetTableInfo(args)
		if err != nil {
			return nil, err
		}

		return &TableResult{
			Table: table,
			Rows:  rows,
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

	default:
		return nil, fmt.Errorf("unknown command: %q", cmd)
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
