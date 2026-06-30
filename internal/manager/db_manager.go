package manager

import (
	"errors"
	"fmt"

	"github.com/harishtpj/indiansql/internal/btree"
	"github.com/harishtpj/indiansql/internal/consts"
	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/pager"
	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
)

type DBManager struct {
	Pager   *pager.Pager
	Catalog *schema.Catalog
	Trees   map[string]*btree.BTree
	dbName  string
}

func NewDBManager(dbFile string) (*DBManager, error) {
	pgr, err := pager.Open(dbFile)
	if err != nil {
		return nil, err
	}

	isNewDB := pgr.PageCount() == 0

	pg, err := pgr.GetPage(0)
	if err != nil {
		return nil, err
	}

	if isNewDB {
		if err := page.InitDBHeader(pg.Data, consts.PageSize); err != nil {
			return nil, err
		}
		cat, err := pgr.InitCatalog()
		if err != nil {
			return nil, err
		}
		pg.Dirty = true

		return &DBManager{
			Pager:   pgr,
			Catalog: cat,
			Trees:   make(map[string]*btree.BTree),
			dbName:  dbFile,
		}, nil
	}

	if _, err = page.DecodeDBHeader(pg.Data); err != nil {
		return nil, err
	}

	cat, err := pgr.LoadCatalog()
	if err != nil {
		return nil, err
	}

	tables := cat.ListTables()
	trees := make(map[string]*btree.BTree)
	for _, table := range tables {
		trees[table.Name] = btree.NewTree(pgr, table.RootPageID)
	}

	return &DBManager{
		Pager:   pgr,
		Catalog: cat,
		Trees:   trees,
		dbName:  dbFile,
	}, nil
}

func (man *DBManager) ExecuteCreateTable(tableName string, columns []*schema.Column) error {
	if man.Catalog.TableExists(tableName) {
		return errors.New("Table already exists: " + tableName)
	}

	table, err := schema.NewTable(tableName, columns)
	if err != nil {
		return err
	}

	if err := man.Catalog.CreateTable(table); err != nil {
		return err
	}

	pg, err := man.Pager.GetPage(table.RootPageID)
	if err != nil {
		return err
	}

	if err := page.InitPage(pg.Data, page.PageTypeLeaf); err != nil {
		return err
	}

	if err := man.Pager.SaveCatalog(man.Catalog); err != nil {
		return err
	}
	man.Trees[tableName] = btree.NewTree(man.Pager, table.RootPageID)
	pg.Dirty = true
	return nil
}

func (man *DBManager) ExecuteInsert(tableName string, values []row.Value) error {
	if !man.Catalog.TableExists(tableName) {
		return errors.New("Table doesn't exist: " + tableName)
	}

	table := man.Catalog.GetTable(tableName)

	// Basic Table Validation
	if table.ColumnCount() != len(values) {
		return fmt.Errorf("Invalid no. of columns supplied: want %d, have %d", table.ColumnCount(), len(values))
	}

	for i, value := range values {
		col := table.Columns[i]
		if value.Type() != col.Type {
			return fmt.Errorf(
				"Invalid datatype for column %s: want %s, have %s",
				col.Name,
				col.Type.String(),
				value.Type().String(),
			)
		}
	}

	insertRow, err := row.NewRow(values, table)
	if err != nil {
		return err
	}

	pk, rawRow, err := row.NewEncoder().EncodeRow(insertRow)
	if err != nil {
		return err
	}

	if err := man.Trees[tableName].Insert(pk, rawRow); err != nil {
		return err
	}

	pg, err := man.Pager.GetPage(table.RootPageID)
	if err != nil {
		return err
	}
	pg.Dirty = true
	return nil
}

func (man *DBManager) ExecuteSelect(tableName string, cols []string, cond string) ([]*row.Row, error) {
	if !man.Catalog.TableExists(tableName) {
		return nil, errors.New("Table doesn't exist: " + tableName)
	}

	table := man.Catalog.GetTable(tableName)
	cursor, err := btree.NewCursor(man.Trees[tableName])
	if err != nil {
		return nil, err
	}
	rows := make([]*row.Row, 0)
	for !cursor.IsFinished() {
		key, err := cursor.Key()
		if err != nil {
			return nil, err
		}

		value, err := cursor.Value()
		if err != nil {
			return nil, err
		}

		r, err := row.NewDecoder().DecodeRow(key, value, table)
		if err != nil {
			return nil, err
		}

		// Filter based on columns
		if cols != nil {
			vals := make([]row.Value, len(cols))
			for i, cname := range cols {
				vals[i], err = r.GetValueByName(cname)
				if err != nil {
					return nil, err
				}
			}

			if r, err = row.NewRow(vals, table); err != nil {
				return nil, err
			}
		}

		rows = append(rows, r)
		cursor.Next()
	}

	return rows, nil
}

func (man *DBManager) ExecuteSelectAll(tableName string) ([]*row.Row, error) {
	return man.ExecuteSelect(tableName, nil, "")
}

func (man *DBManager) Commit() error {
	if err := man.Pager.SaveCatalog(man.Catalog); err != nil {
		return err
	}

	return man.Pager.FlushToFile()
}

func (man *DBManager) Close() error {
	if err := man.Pager.SaveCatalog(man.Catalog); err != nil {
		return err
	}

	return man.Pager.Close()
}

func (man *DBManager) GetTableInfo(tableName string) (*schema.Table, error) {
	return man.Catalog.GetTable(tableName), nil
}

func (man *DBManager) GetDBName() string {
	return man.dbName
}
