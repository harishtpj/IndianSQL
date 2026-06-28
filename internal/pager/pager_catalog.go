package pager

import (
	"github.com/harishtpj/indiansql/internal/page"
	"github.com/harishtpj/indiansql/internal/schema"
)

func (p *Pager) LoadCatalog() (*schema.Catalog, error) {
	pg, err := p.GetPage(0)
	if err != nil {
		return nil, err
	}

	cat := schema.NewCatalog()
	err = cat.Deserialize(pg.Data[page.DBHeaderSize:])
	return cat, err
}

func (p *Pager) SaveCatalog(cat *schema.Catalog) error {
	pg, err := p.GetPage(0)
	if err != nil {
		return err
	}
	serialized, err := cat.Serialize()
	if err != nil {
		return err
	}
	copy(pg.Data[page.DBHeaderSize:], serialized)
	pg.Dirty = true
	return nil
}

func (p *Pager) InitCatalog() (cat *schema.Catalog, err error) {
	cat = schema.NewCatalog()
	err = p.SaveCatalog(cat)
	return
}
