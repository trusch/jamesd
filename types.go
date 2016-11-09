package jamesd

import "encoding/gob"

func init() {
	gob.Register(&Request{})
	gob.Register(&Response{})
}

type CompressionType int

const (
	NONE CompressionType = iota
	GZIP
	BZIP2
	LZMA
)

type Request struct {
	Package string
	Version string
	Arch    string
}

type Response struct {
	*Package
	Error string
}
