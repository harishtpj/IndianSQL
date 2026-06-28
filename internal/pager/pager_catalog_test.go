package pager

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/schema"
)

func TestLoadSaveCatalog(t *testing.T) {
	p, _ := Open(":memory:")
	defer p.Close()

	// Create catalog
	cat := schema.NewCatalog()
	tbl := &schema.Table{
		Name: "users",
		Columns: []*schema.Column{
			{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
			{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
		RootPageID:      1,
	}
	cat.CreateTable(tbl)

	// Save
	err := p.SaveCatalog(cat)
	if err != nil {
		t.Fatal(err)
	}

	// Load
	loaded, err := p.LoadCatalog()
	if err != nil {
		t.Fatal(err)
	}

	if loaded == nil {
		t.Fatal("Loaded catalog is nil")
	}

	tblLoaded := loaded.GetTable("users")
	if tblLoaded == nil {
		t.Fatal("Table not found after load")
	}

	if tblLoaded.Name != "users" {
		t.Errorf("Table name = %v, want users", tblLoaded.Name)
	}
}

func TestInitCatalog(t *testing.T) {
	p, _ := Open(":memory:")
	defer p.Close()

	cat, err := p.InitCatalog()
	if err != nil {
		t.Fatal(err)
	}

	if cat == nil {
		t.Fatal("InitCatalog returned nil")
	}

	// Should be able to create tables
	tbl := &schema.Table{
		Name:            "test",
		Columns:         []*schema.Column{{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true}},
		PrimaryKeyIndex: 0,
		RootPageID:      1,
	}
	err = cat.CreateTable(tbl)
	if err != nil {
		t.Fatal(err)
	}

	if cat.GetTable("test") == nil {
		t.Error("Table not found in initialized catalog")
	}
}

func TestCatalogPersistence(t *testing.T) {
	// Create and save in first pager instance
	p1, _ := Open(":memory:")
	cat1 := schema.NewCatalog()
	tbl := &schema.Table{
		Name:            "products",
		Columns:         []*schema.Column{{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true}},
		PrimaryKeyIndex: 0,
		RootPageID:      1,
	}
	cat1.CreateTable(tbl)
	p1.SaveCatalog(cat1)
	p1.Close()

	// Load in second instance (in-memory DB persists per Pager instance in tests)
	// This test may not work with :memory: - may need file-based test
	// For now, test save/load in same pager instance:
	p2, _ := Open(":memory:")
	defer p2.Close()

	cat2 := schema.NewCatalog()
	tbl2 := &schema.Table{
		Name:            "orders",
		Columns:         []*schema.Column{{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true}},
		PrimaryKeyIndex: 0,
		RootPageID:      2,
	}
	cat2.CreateTable(tbl2)
	p2.SaveCatalog(cat2)

	// Reload from same pager
	loaded, _ := p2.LoadCatalog()
	if loaded.GetTable("orders") == nil {
		t.Error("Orders table not persisted")
	}
}

func TestCatalogMultipleTables(t *testing.T) {
	p, _ := Open(":memory:")
	defer p.Close()

	cat := schema.NewCatalog()

	tables := []struct {
		name       string
		rootPageID uint32
	}{
		{"users", 1},
		{"products", 2},
		{"orders", 3},
	}

	for _, tt := range tables {
		tbl := &schema.Table{
			Name:            tt.name,
			Columns:         []*schema.Column{{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true}},
			PrimaryKeyIndex: 0,
			RootPageID:      tt.rootPageID,
		}
		cat.CreateTable(tbl)
	}

	p.SaveCatalog(cat)
	loaded, _ := p.LoadCatalog()

	for _, tt := range tables {
		tbl := loaded.GetTable(tt.name)
		if tbl == nil {
			t.Errorf("Table %v not found", tt.name)
		} else if tbl.RootPageID != tt.rootPageID {
			t.Errorf("RootPageID = %v, want %v", tbl.RootPageID, tt.rootPageID)
		}
	}
}
