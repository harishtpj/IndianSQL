package page

import (
	"bytes"
	"errors"
	"testing"

	"github.com/harishtpj/indiansql/internal/apperrors"
	"github.com/harishtpj/indiansql/internal/consts"
)

func TestInitDBHeader(t *testing.T) {
	page := make([]byte, 100)

	err := InitDBHeader(page, 4096)
	if err != nil {
		t.Fatal(err)
	}

	hdr, err := DecodeDBHeader(page)
	if err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(hdr.Magic[:], []byte(consts.MagicStr)) {
		t.Fatalf("magic mismatch")
	}

	if hdr.Version != consts.VersionNum {
		t.Fatalf("version mismatch")
	}

	if hdr.PageSize != 4096 {
		t.Fatalf("page size mismatch")
	}

	if hdr.PageCount != 1 {
		t.Fatalf("page count mismatch")
	}

	if hdr.FreeListHead != 0 {
		t.Fatalf("freelist head not zero")
	}

	if hdr.FreeListCount != 0 {
		t.Fatalf("freelist count not zero")
	}

	if hdr.RootPage != 0 {
		t.Fatalf("root page not zero")
	}

	if !bytes.Equal(hdr.Reserved[:], make([]byte, reservedSize)) {
		t.Fatalf("reserved not zeroed")
	}
}

func TestDBHeaderEncodeDecodeRoundTrip(t *testing.T) {
	var hdr DBHeader

	copy(hdr.Magic[:], []byte(consts.MagicStr))
	hdr.Version = consts.VersionNum
	hdr.PageSize = 8192
	hdr.PageCount = 42
	hdr.FreeListHead = 3
	hdr.FreeListCount = 7
	hdr.RootPage = 2

	for i := range hdr.Reserved {
		hdr.Reserved[i] = byte(i)
	}

	page := make([]byte, 100)

	if err := hdr.Encode(page); err != nil {
		t.Fatal(err)
	}

	decoded, err := DecodeDBHeader(page)
	if err != nil {
		t.Fatal(err)
	}

	if *decoded != hdr {
		t.Fatalf("roundtrip mismatch")
	}
}

func TestDBHeaderRejectsBadMagic(t *testing.T) {
	page := make([]byte, 100)

	copy(page, []byte("INVALID HEADER"))

	_, err := DecodeDBHeader(page)
	if err == nil {
		t.Fatalf("expected error for bad magic")
	}

	if !errors.Is(err, apperrors.ErrInvalidHeader) {
		t.Fatalf("expected ErrInvalidHeader, got %v", err)
	}
}

func TestDBHeaderRejectsWrongVersion(t *testing.T) {
	page := make([]byte, 100)

	err := InitDBHeader(page, 4096)
	if err != nil {
		t.Fatal(err)
	}

	page[versionOffset] = 0xFF
	page[versionOffset+1] = 0xFF

	_, err = DecodeDBHeader(page)
	if err == nil {
		t.Fatalf("expected version error")
	}

	if !errors.Is(err, apperrors.ErrInvalidHeader) {
		t.Fatalf("expected ErrInvalidHeader, got %v", err)
	}
}

func TestDBHeaderRejectsShortBuffer(t *testing.T) {
	_, err := DecodeDBHeader(make([]byte, 10))
	if err == nil {
		t.Fatalf("expected short buffer error")
	}

	if !errors.Is(err, apperrors.ErrDBHeaderSmall) {
		t.Fatalf("expected ErrDBHeaderSmall, got %v", err)
	}
}

func TestDBHeaderEncodeDeterministic(t *testing.T) {
	var hdr DBHeader

	copy(hdr.Magic[:], []byte(consts.MagicStr))
	hdr.Version = consts.VersionNum
	hdr.PageSize = 4096
	hdr.PageCount = 10

	p1 := make([]byte, 100)
	p2 := make([]byte, 100)

	if err := hdr.Encode(p1); err != nil {
		t.Fatal(err)
	}

	if err := hdr.Encode(p2); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(p1, p2) {
		t.Fatalf("encoding not deterministic")
	}
}
