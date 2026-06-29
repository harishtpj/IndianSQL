package schema

import (
	"testing"
)

// Test NewCatalog()
func TestNewCatalog(t *testing.T) {
	cat := NewCatalog()

	if cat == nil {
		t.Fatal("NewCatalog() returned nil")
	}

	if cat.TableCount() != 0 {
		t.Errorf("NewCatalog() initial table count = %v, want 0", cat.TableCount())
	}

	if cat.GetNextRootPageID() != 1 {
		t.Errorf("NewCatalog() initial nextRootPageID = %v, want 1 (page 0 reserved)", cat.GetNextRootPageID())
	}
}

// Test Catalog.CreateTable()
func TestCatalogCreateTable(t *testing.T) {
	cat := NewCatalog()

	tests := []struct {
		name      string
		table     *Table
		shouldErr bool
		errMsg    string
	}{
		{
			name: "Create single table",
			table: &Table{
				Name: "users",
				Columns: []*Column{
					{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
					{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
				},
				PrimaryKeyIndex: 0,
			},
			shouldErr: false,
		},
		{
			name: "Duplicate table name should error",
			table: &Table{
				Name: "users",
				Columns: []*Column{
					{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
				},
				PrimaryKeyIndex: 0,
			},
			shouldErr: true,
			errMsg:    "exists",
		},
		{
			name: "Invalid table should error",
			table: &Table{
				Name:            "",
				Columns:         []*Column{},
				PrimaryKeyIndex: 0,
			},
			shouldErr: true,
			errMsg:    "invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cat.CreateTable(tt.table)

			if (err != nil) != tt.shouldErr {
				t.Errorf("CreateTable() error = %v, shouldErr = %v", err, tt.shouldErr)
			}

			// If no error, table should be retrievable
			if !tt.shouldErr && tt.table.Name != "" {
				retrieved := cat.GetTable(tt.table.Name)
				if retrieved == nil {
					t.Errorf("CreateTable() table not found after creation")
				}
				// RootPageID should be assigned
				if retrieved.RootPageID == 0 {
					t.Errorf("CreateTable() did not assign RootPageID")
				}
			}
		})
	}
}

// Test Catalog.GetTable()
func TestCatalogGetTable(t *testing.T) {
	cat := NewCatalog()

	// Create a table
	table := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
		},
		PrimaryKeyIndex: 0,
	}
	cat.CreateTable(table)

	tests := []struct {
		name          string
		tableName     string
		expectedFound bool
	}{
		{
			name:          "Get existing table",
			tableName:     "users",
			expectedFound: true,
		},
		{
			name:          "Get non-existent table",
			tableName:     "nonexistent",
			expectedFound: false,
		},
		{
			name:          "Case sensitive lookup",
			tableName:     "Users",
			expectedFound: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cat.GetTable(tt.tableName)

			if (result != nil) != tt.expectedFound {
				t.Errorf("GetTable() found = %v, want %v", result != nil, tt.expectedFound)
			}

			if tt.expectedFound && result != nil && result.Name != tt.tableName {
				t.Errorf("GetTable() returned wrong table")
			}
		})
	}
}

// Test Catalog.TableExists()
func TestCatalogTableExists(t *testing.T) {
	cat := NewCatalog()

	table := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
		},
		PrimaryKeyIndex: 0,
	}
	cat.CreateTable(table)

	tests := []struct {
		name     string
		table    string
		expected bool
	}{
		{
			name:     "Existing table",
			table:    "users",
			expected: true,
		},
		{
			name:     "Non-existent table",
			table:    "products",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := cat.TableExists(tt.table)

			if result != tt.expected {
				t.Errorf("TableExists() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// Test Catalog.ListTables()
func TestCatalogListTables(t *testing.T) {
	cat := NewCatalog()

	// Create multiple tables
	tables := []*Table{
		{
			Name: "users",
			Columns: []*Column{
				{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
			},
			PrimaryKeyIndex: 0,
		},
		{
			Name: "products",
			Columns: []*Column{
				{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
			},
			PrimaryKeyIndex: 0,
		},
	}

	for _, table := range tables {
		cat.CreateTable(table)
	}

	list := cat.ListTables()

	if len(list) != len(tables) {
		t.Errorf("ListTables() returned %v tables, want %v", len(list), len(tables))
	}

	// Check that all created tables are in the list
	for _, expected := range tables {
		found := false
		for _, actual := range list {
			if actual.Name == expected.Name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("ListTables() missing table '%s'", expected.Name)
		}
	}
}

// Test Catalog.TableCount()
func TestCatalogTableCount(t *testing.T) {
	cat := NewCatalog()

	if cat.TableCount() != 0 {
		t.Errorf("Initial TableCount() = %v, want 0", cat.TableCount())
	}

	// Add tables and check count
	for i := 0; i < 3; i++ {
		table := &Table{
			Name: "table" + string(rune('0'+i)),
			Columns: []*Column{
				{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
			},
			PrimaryKeyIndex: 0,
		}
		cat.CreateTable(table)

		if cat.TableCount() != i+1 {
			t.Errorf("TableCount() = %v, want %v after adding table", cat.TableCount(), i+1)
		}
	}
}

// Test Catalog.GetNextRootPageID() and page ID assignment
func TestCatalogRootPageIDAssignment(t *testing.T) {
	cat := NewCatalog()

	initialID := cat.GetNextRootPageID()
	if initialID != 1 {
		t.Errorf("Initial GetNextRootPageID() = %v, want 1", initialID)
	}

	// Create first table
	table1 := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
		},
		PrimaryKeyIndex: 0,
	}
	cat.CreateTable(table1)

	if table1.RootPageID != 1 {
		t.Errorf("First table RootPageID = %v, want 1", table1.RootPageID)
	}

	// Create second table
	table2 := &Table{
		Name: "products",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
		},
		PrimaryKeyIndex: 0,
	}
	cat.CreateTable(table2)

	if table2.RootPageID != 2 {
		t.Errorf("Second table RootPageID = %v, want 2", table2.RootPageID)
	}

	if cat.GetNextRootPageID() != 3 {
		t.Errorf("GetNextRootPageID() after two tables = %v, want 3", cat.GetNextRootPageID())
	}
}

// Test Catalog.Serialize() and Deserialize()
func TestCatalogSerializeDeserialize(t *testing.T) {
	cat1 := NewCatalog()

	// Create test tables
	table1 := &Table{
		Name: "users",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
			{Name: "name", Type: ColumnTypeVarchar, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}
	cat1.CreateTable(table1)

	table2 := &Table{
		Name: "products",
		Columns: []*Column{
			{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
			{Name: "title", Type: ColumnTypeVarchar, IsPrimaryKey: false},
			{Name: "price", Type: ColumnTypeInteger, IsPrimaryKey: false},
		},
		PrimaryKeyIndex: 0,
	}
	cat1.CreateTable(table2)

	// Serialize
	data, err := cat1.Serialize()
	if err != nil {
		t.Fatalf("Serialize() error = %v", err)
	}

	if len(data) == 0 {
		t.Fatal("Serialize() returned empty data")
	}

	// Deserialize into new catalog
	cat2 := NewCatalog()
	cat2.Clear() // Start fresh
	err = cat2.Deserialize(data)
	if err != nil {
		t.Fatalf("Deserialize() error = %v", err)
	}

	// Verify deserialized catalog matches original
	if cat2.TableCount() != cat1.TableCount() {
		t.Errorf("Deserialized table count = %v, want %v", cat2.TableCount(), cat1.TableCount())
	}

	// Verify table1
	restored1 := cat2.GetTable("users")
	if restored1 == nil {
		t.Fatal("Deserialized catalog missing table 'users'")
	}
	if restored1.ColumnCount() != 2 {
		t.Errorf("Restored users table column count = %v, want 2", restored1.ColumnCount())
	}
	if restored1.GetColumn("name") == nil {
		t.Error("Restored users table missing 'name' column")
	}

	// Verify table2
	restored2 := cat2.GetTable("products")
	if restored2 == nil {
		t.Fatal("Deserialized catalog missing table 'products'")
	}
	if restored2.ColumnCount() != 3 {
		t.Errorf("Restored products table column count = %v, want 3", restored2.ColumnCount())
	}
	if restored2.GetColumn("price") == nil {
		t.Error("Restored products table missing 'price' column")
	}

	// Verify RootPageID is restored correctly
	if restored1.RootPageID != table1.RootPageID {
		t.Errorf("Restored users RootPageID = %v, want %v", restored1.RootPageID, table1.RootPageID)
	}
	if restored2.RootPageID != table2.RootPageID {
		t.Errorf("Restored products RootPageID = %v, want %v", restored2.RootPageID, table2.RootPageID)
	}
}

// Test Catalog.Deserialize() with corrupted data
func TestCatalogDeserializeCorrupted(t *testing.T) {
	cat := NewCatalog()

	tests := []struct {
		name      string
		data      []byte
		shouldErr bool
	}{
		{
			name:      "Empty data",
			data:      []byte{},
			shouldErr: true,
		},
		{
			name:      "Invalid magic string",
			data:      []byte("INVALID" + string(make([]byte, 100))),
			shouldErr: true,
		},
		{
			name:      "Truncated data",
			data:      []byte("SHORT"),
			shouldErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cat.Deserialize(tt.data)

			if (err != nil) != tt.shouldErr {
				t.Errorf("Deserialize() error = %v, shouldErr = %v", err, tt.shouldErr)
			}
		})
	}
}

// Test Catalog.Clear()
func TestCatalogClear(t *testing.T) {
	cat := NewCatalog()

	// Add some tables
	for i := 0; i < 3; i++ {
		table := &Table{
			Name: "table" + string(rune('0'+i)),
			Columns: []*Column{
				{Name: "id", Type: ColumnTypeInteger, IsPrimaryKey: true},
			},
			PrimaryKeyIndex: 0,
		}
		cat.CreateTable(table)
	}

	if cat.TableCount() == 0 {
		t.Fatal("Test setup failed: no tables created")
	}

	// Clear catalog
	cat.Clear()

	if cat.TableCount() != 0 {
		t.Errorf("After Clear(), TableCount() = %v, want 0", cat.TableCount())
	}

	if cat.GetNextRootPageID() != 1 {
		t.Errorf("After Clear(), GetNextRootPageID() = %v, want 1", cat.GetNextRootPageID())
	}

	if len(cat.ListTables()) != 0 {
		t.Errorf("After Clear(), ListTables() returned %v tables, want 0", len(cat.ListTables()))
	}
}
