package engine

import "testing"

func TestExecuteMySQLShowTables(t *testing.T) {
	e, err := NewSQLEngine(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer e.db.Close()

	if _, err := e.Execute("create table users (id int primary key, name varchar)"); err != nil {
		t.Fatal(err)
	}

	res, err := e.Execute("show tables")
	if err != nil {
		t.Fatal(err)
	}

	tables, ok := res.(*TablesResult)
	if !ok {
		t.Fatalf("got %T, want *TablesResult", res)
	}
	if len(tables.Tables) != 1 || tables.Tables[0].Name != "users" {
		t.Fatalf("unexpected tables result: %+v", tables.Tables)
	}
}

func TestExecuteMySQLDescAndShowColumns(t *testing.T) {
	e, err := NewSQLEngine(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer e.db.Close()

	if _, err := e.Execute("create table users (id int primary key, name varchar)"); err != nil {
		t.Fatal(err)
	}

	for _, query := range []string{"desc users"} {
		res, err := e.Execute(query)
		if err != nil {
			t.Fatalf("%s: %v", query, err)
		}

		schemaRes, ok := res.(*SchemaResult)
		if !ok {
			t.Fatalf("%s: got %T, want *SchemaResult", query, res)
		}
		if schemaRes.Table == nil || schemaRes.Table.Name != "users" {
			t.Fatalf("%s: unexpected table: %+v", query, schemaRes.Table)
		}
		if len(schemaRes.Table.Columns) != 2 {
			t.Fatalf("%s: unexpected column count %d", query, len(schemaRes.Table.Columns))
		}
	}
}
