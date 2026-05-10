package pager

import (
	"io"

	"github.com/harishtpj/havensql/internal/storage"
)

type Pager struct {
	file   storage.Storage
	pages  map[uint32]*Page
	nPages uint32
}

func Open(path string) (*Pager, error) {
	var file storage.Storage
	var nPages uint32

	if path == ":memory:" {
		file = storage.NewMemoryStorage()
		nPages = 0
	} else {
		var err error
		file, err = storage.NewFileStorage(path)
		if err != nil {
			return nil, err
		}
		nBytes, err := file.Size()
		if err != nil {
			return nil, err
		}
		nPages = uint32((nBytes + PageSize - 1) / PageSize)
	}
	pages := make(map[uint32]*Page)
	return &Pager{file, pages, nPages}, nil
}

func (p *Pager) GetPage(pNum uint32) (*Page, error) {
	if pg, ok := p.pages[pNum]; ok {
		return pg, nil
	}

	pg := &Page{
		Data: make([]byte, PageSize),
	}

	if pNum < p.nPages {
		offset := int64(pNum) * PageSize
		_, err := p.file.ReadAt(pg.Data, offset)
		if err != nil && err != io.EOF {
			return nil, err
		}
	} else {
		p.nPages = pNum + 1
	}

	p.pages[pNum] = pg
	return pg, nil
}

func (p *Pager) FlushPage(pNum uint32) error {
	pg, ok := p.pages[pNum]
	if !ok || pg == nil {
		return nil
	}
	offset := int64(pNum) * PageSize
	_, err := p.file.WriteAt(pg.Data, offset)
	if err != nil {
		return err
	}
	pg.Dirty = false
	return nil
}

func (p *Pager) FlushAll() error {
	for i, pg := range p.pages {
		if pg != nil && pg.Dirty {
			if err := p.FlushPage(i); err != nil {
				return err
			}
		}
	}
	return nil
}

func (p *Pager) Close() error {
	if err := p.FlushAll(); err != nil {
		return err
	}
	if err := p.file.Close(); err != nil {
		return err
	}
	return nil
}
