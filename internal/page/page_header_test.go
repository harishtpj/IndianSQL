package page

import (
	"errors"
	"testing"

	"github.com/harishtpj/indiansql/internal/apperrors"
)

const testPageSize = 4096

func TestPageHeaderEncodeDecode(t *testing.T) {
	buf := make([]byte, testPageSize)

	h := PageHeader{
		Type:      PageTypeLeaf,
		FreeStart: 32,
		FreeEnd:   testPageSize,
		CellCount: 3,
	}

	if err := h.Encode(buf); err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodePageHeader(buf)
	if err != nil {
		t.Fatal(err)
	}

	if *decoded != h {
		t.Fatalf("page header mismatch\n%+v\n%+v", *decoded, h)
	}
}

func TestInitPage(t *testing.T) {
	buf := make([]byte, testPageSize)

	if err := InitPage(buf, PageTypeLeaf); err != nil {
		t.Fatal(err)
	}

	h, err := DecodePageHeader(buf)
	if err != nil {
		t.Fatal(err)
	}

	if h.Type != PageTypeLeaf {
		t.Fatalf("wrong page type")
	}

	if h.CellCount != 0 {
		t.Fatalf("cell count must be zero")
	}

	if h.FreeStart == 0 {
		t.Fatalf("free start not initialized")
	}

	if h.FreeEnd != testPageSize {
		t.Fatalf("free end incorrect")
	}

	if h.FreeStart >= h.FreeEnd {
		t.Fatalf("invalid free space bounds")
	}
}

func TestDecodePageHeaderInvalidSize(t *testing.T) {
	buf := make([]byte, 4)

	_, err := DecodePageHeader(buf)
	if err == nil {
		t.Fatalf("expected error for short buffer")
	}

	if !errors.Is(err, apperrors.ErrInvalidHeader) {
		t.Fatalf("expected ErrInvalidHeader, got %v", err)
	}
}

func TestDecodePageHeaderInvalidType(t *testing.T) {
	buf := make([]byte, testPageSize)

	buf[0] = 255

	_, err := DecodePageHeader(buf)
	if err == nil {
		t.Fatalf("expected invalid page type error")
	}

	if !errors.Is(err, apperrors.ErrInvalidHeader) {
		t.Fatalf("expected ErrInvalidHeader, got %v", err)
	}
}
