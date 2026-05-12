package page

type PageType uint8

const (
	PageTypeFree PageType = iota
	PageTypeMeta
	PageTypeLeaf
	PageTypeInternal
	PageTypeOverflow
	MaxPageType
)
