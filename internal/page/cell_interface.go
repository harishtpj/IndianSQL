package page

type PageCell interface {
	Size() uint16
	Encode(dst []byte)
}
