package packet

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"io/ioutil"

	"github.com/ulikunitz/xz"
)

type Packet struct {
	Name              string
	Arch              string
	Version           string
	Data              []byte
	Compression       CompressionType
	PreInstallScript  string
	PostInstallScript string
}

type CompressionType int

const (
	NO_COMPRESSION CompressionType = iota
	GZIP
	BZIP2
	LZMA
)

func New(name, arch, version string, data io.Reader, compression CompressionType, preInstallScript, postInstallScript string) (*Packet, error) {
	bs, err := ioutil.ReadAll(data)
	if err != nil {
		return nil, err
	}
	return &Packet{
		Name:              name,
		Arch:              arch,
		Version:           version,
		Data:              bs,
		Compression:       compression,
		PreInstallScript:  preInstallScript,
		PostInstallScript: postInstallScript,
	}, nil
}

func (packet *Packet) GetTarReader() (*tar.Reader, error) {
	if packet.Data == nil || len(packet.Data) == 0 {
		return nil, errors.New("empty archive")
	}
	var compressedReader io.Reader
	byteReader := bytes.NewReader(packet.Data)
	switch packet.Compression {
	case NO_COMPRESSION:
		{
			compressedReader = byteReader
		}
	case GZIP:
		{
			gzipReader, err := gzip.NewReader(byteReader)
			if err != nil {
				return nil, err
			}
			compressedReader = gzipReader
		}
	case BZIP2:
		{
			compressedReader = bzip2.NewReader(byteReader)
		}
	case LZMA:
		{
			r, err := xz.NewReader(byteReader)
			if err != nil {
				return nil, err
			}
			compressedReader = r
		}
	}
	return tar.NewReader(compressedReader), nil
}
