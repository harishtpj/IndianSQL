package pager

const PageSize = 4096

type Page struct {
	Data  []byte
	Dirty bool
}
