package server

import (
	"testing"

	"github.com/harishtpj/indiansql/internal/engine"
)

func TestHandleQueryShowTables(t *testing.T) {
	e, err := engine.NewSQLEngine(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = e.Execute("exit") }()

	if _, err := e.Execute("create table users (id int primary key, name varchar)"); err != nil {
		t.Fatal(err)
	}

	h := NewHandler(e)
	res, err := h.HandleQuery("show tables")
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasResultset() {
		t.Fatal("expected resultset")
	}

	if got := res.Resultset.ColumnNumber(); got != 1 {
		t.Fatalf("column count = %d, want 1", got)
	}
	if got := len(res.Resultset.RowDatas); got != 1 {
		t.Fatalf("row count = %d, want 1", got)
	}
	if got := string(res.Resultset.Fields[0].Name); got != "Tables_in_:memory:" {
		t.Fatalf("column name = %v, want Tables_in_:memory:", got)
	}
}

func TestHandleQueryDesc(t *testing.T) {
	e, err := engine.NewSQLEngine(":memory:")
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _, _ = e.Execute("exit") }()

	if _, err := e.Execute("create table users (id int primary key, name varchar)"); err != nil {
		t.Fatal(err)
	}

	h := NewHandler(e)
	res, err := h.HandleQuery("desc users")
	if err != nil {
		t.Fatal(err)
	}
	if !res.HasResultset() {
		t.Fatal("expected resultset")
	}

	if got := res.Resultset.ColumnNumber(); got != 6 {
		t.Fatalf("column count = %d, want 6", got)
	}
	if got := len(res.Resultset.RowDatas); got != 2 {
		t.Fatalf("row count = %d, want 2", got)
	}
	if got := string(res.Resultset.Fields[0].Name); got != "Field" {
		t.Fatalf("first column name = %v, want Field", got)
	}
	if got := string(res.Resultset.Fields[3].Name); got != "Key" {
		t.Fatalf("fourth column name = %v, want Key", got)
	}
}
