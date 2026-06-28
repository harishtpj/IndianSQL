package cli

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/row"
	"github.com/harishtpj/indiansql/internal/schema"
)

func TestNewREPLContext(t *testing.T) {
	ctx, err := NewREPLContext(":memory:")
	if err != nil {
		t.Fatal(err)
	}

	if ctx == nil {
		t.Fatal("Context is nil")
	}

	if ctx.Pager == nil {
		t.Error("Pager is nil")
	}

	if ctx.Catalog == nil {
		t.Error("Catalog is nil")
	}

	ctx.Close()
}

func TestExecuteCreateTable(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
	}

	err := ctx.ExecuteCreateTable("users", columns)
	if err != nil {
		t.Fatal(err)
	}

	tbl := ctx.Catalog.GetTable("users")
	if tbl == nil {
		t.Fatal("Table not created")
	}

	if tbl.Name != "users" {
		t.Errorf("Table name = %v, want users", tbl.Name)
	}

	if len(tbl.Columns) != 2 {
		t.Errorf("Column count = %v, want 2", len(tbl.Columns))
	}
}

func TestExecuteCreateTableDuplicate(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
	}

	ctx.ExecuteCreateTable("users", columns)

	// Create same table again
	err := ctx.ExecuteCreateTable("users", columns)
	if err == nil {
		t.Error("Should error on duplicate table")
	}
}

func TestExecuteInsert(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
	}
	ctx.ExecuteCreateTable("users", columns)

	values := []row.Value{
		row.NewNumericValue(1),
		row.NewVarcharValue("Alice"),
	}

	err := ctx.ExecuteInsert("users", values)
	if err != nil {
		t.Fatal(err)
	}
}

func TestExecuteInsertInvalidTable(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	values := []row.Value{
		row.NewNumericValue(1),
	}

	err := ctx.ExecuteInsert("nonexistent", values)
	if err == nil {
		t.Error("Should error on missing table")
	}
}

func TestExecuteInsertWrongColumnCount(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
	}
	ctx.ExecuteCreateTable("users", columns)

	// Insert with wrong column count
	values := []row.Value{
		row.NewNumericValue(1),
		// Missing name column
	}

	err := ctx.ExecuteInsert("users", values)
	if err == nil {
		t.Error("Should error on column count mismatch")
	}
}

func TestExecuteSelectAll(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		{Name: "value", Type: schema.ColumnTypeNumeric, IsPrimaryKey: false},
	}
	ctx.ExecuteCreateTable("test", columns)

	// Insert rows
	for i := 0; i < 5; i++ {
		values := []row.Value{
			row.NewNumericValue(int64(i)),
			row.NewNumericValue(int64(i * 10)),
		}
		ctx.ExecuteInsert("test", values)
	}

	// Select all
	rows, err := ctx.ExecuteSelectAll("test")
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 5 {
		t.Errorf("Row count = %v, want 5", len(rows))
	}
}

func TestExecuteSelectAllEmpty(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
	}
	ctx.ExecuteCreateTable("empty", columns)

	rows, err := ctx.ExecuteSelectAll("empty")
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != 0 {
		t.Errorf("Should have 0 rows, got %v", len(rows))
	}
}

func TestExecuteSelectAllInvalidTable(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	_, err := ctx.ExecuteSelectAll("nonexistent")
	if err == nil {
		t.Error("Should error on missing table")
	}
}

func TestGetTableInfo(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
	}
	ctx.ExecuteCreateTable("users", columns)

	tbl, err := ctx.GetTableInfo("users")
	if err != nil {
		t.Fatal(err)
	}

	if tbl == nil {
		t.Fatal("Table info is nil")
	}

	if tbl.Name != "users" {
		t.Errorf("Table name = %v, want users", tbl.Name)
	}
}

func TestCreateInsertSelectFlow(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")
	defer ctx.Close()

	// CREATE TABLE
	columns := []*schema.Column{
		{Name: "id", Type: schema.ColumnTypeNumeric, IsPrimaryKey: true},
		{Name: "name", Type: schema.ColumnTypeVarchar, IsPrimaryKey: false},
		{Name: "age", Type: schema.ColumnTypeNumeric, IsPrimaryKey: false},
	}
	ctx.ExecuteCreateTable("persons", columns)

	// INSERT multiple rows
	inserts := []struct {
		id   int64
		name string
		age  int64
	}{
		{1, "Alice", 30},
		{2, "Bob", 25},
		{3, "Charlie", 35},
	}

	for _, ins := range inserts {
		values := []row.Value{
			row.NewNumericValue(ins.id),
			row.NewVarcharValue(ins.name),
			row.NewNumericValue(ins.age),
		}
		ctx.ExecuteInsert("persons", values)
	}

	// SELECT *
	rows, err := ctx.ExecuteSelectAll("persons")
	if err != nil {
		t.Fatal(err)
	}

	if len(rows) != len(inserts) {
		t.Errorf("Got %v rows, want %v", len(rows), len(inserts))
	}

	// Verify data
	for i, r := range rows {
		v0, _ := r.GetValue(0)
		id := v0.(*row.NumericValue).GetInt64()

		if id != inserts[i].id {
			t.Errorf("Row %v ID = %v, want %v", i, id, inserts[i].id)
		}
	}
}

func TestContextClose(t *testing.T) {
	ctx, _ := NewREPLContext(":memory:")

	err := ctx.Close()
	if err != nil {
		t.Fatal(err)
	}

	// After close, pager should be closed (can't test directly without exposing state)
	// But context should be in valid closed state
	if ctx.Pager == nil {
		t.Error("Pager should still be set after close")
	}
}
