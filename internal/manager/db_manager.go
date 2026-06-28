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
	}, nil
}

func (rc *DBManager) ExecuteCreateTable(tableName string, columns []*schema.Column) error {
	if rc.Catalog.TableExists(tableName) {
		return errors.New("Table already exists: " + tableName)
	}

	table, err := schema.NewTable(tableName, columns)
	if err != nil {
		return err
	}

	if err := rc.Catalog.CreateTable(table); err != nil {
		return err
	}

	pg, err := rc.Pager.GetPage(table.RootPageID)
	if err != nil {
		return err
	}

	if err := page.InitPage(pg.Data, page.PageTypeLeaf); err != nil {
		return err
	}

	if err := rc.Pager.SaveCatalog(rc.Catalog); err != nil {
		return err
	}
	rc.Trees[tableName] = btree.NewTree(rc.Pager, table.RootPageID)
	pg.Dirty = true
	return nil
}

func (rc *DBManager) ExecuteInsert(tableName string, values []row.Value) error {
	if !rc.Catalog.TableExists(tableName) {
		return errors.New("Table doesn't exist: " + tableName)
	}

	table := rc.Catalog.GetTable(tableName)

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

	return rc.Trees[tableName].Insert(pk, rawRow)
}

func (rc *DBManager) ExecuteSelectAll(tableName string) ([]*row.Row, error) {
	if !rc.Catalog.TableExists(tableName) {
		return nil, errors.New("Table doesn't exist: " + tableName)
	}

	table := rc.Catalog.GetTable(tableName)
	cursor, err := btree.NewCursor(rc.Trees[tableName])
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
		rows = append(rows, r)
		cursor.Next()
	}

	return rows, nil
}

func (rc *DBManager) Close() error {
	if err := rc.Pager.SaveCatalog(rc.Catalog); err != nil {
		return err
	}

	if err := rc.Pager.Close(); err != nil {
		return err
	}

	return nil
}

func (rc *DBManager) GetTableInfo(tableName string) (*schema.Table, error) {
	return rc.Catalog.GetTable(tableName), nil
}
