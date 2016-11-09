package jamesd

import (
	"archive/tar"
	"bytes"
	"compress/bzip2"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/ulikunitz/xz"
)

func GetPackage(request *Request, db *DB) *Response {
	response := &Response{}
	pack, err := db.GetPackage(request.Package, request.Arch, request.Version)
	if err != nil {
		response.Error = err.Error()
		return response
	}
	response.Package = pack
	return response
}

func InstallPackage(response *Response, installRoot string) error {
	if err := untarPackage(response.Data, response.Compression, installRoot); err != nil {
		return err
	}
	if response.Script != nil && len(response.Script) > 0 {
		cmd := exec.Command("sh", "-c", string(response.Script))
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return err
		}
	}
	return nil
}

func untarPackage(data []byte, compression CompressionType, installRoot string) error {
	if data == nil || len(data) == 0 {
		return errors.New("empty archive")
	}
	var compressedReader io.Reader
	byteReader := bytes.NewReader(data)
	switch compression {
	case NONE:
		{
			compressedReader = byteReader
		}
	case GZIP:
		{
			gzipReader, err := gzip.NewReader(byteReader)
			if err != nil {
				return err
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
				return err
			}
			compressedReader = r
		}
	}
	tarReader := tar.NewReader(compressedReader)
	for {
		hdr, err := tarReader.Next()
		if err == io.EOF {
			break
		}
		if hdr.FileInfo().IsDir() {
			err = os.MkdirAll(filepath.Join(installRoot, hdr.Name), os.FileMode(hdr.Mode))
			if err != nil {
				return err
			}
		} else {
			path := filepath.Join(installRoot, hdr.Name)
			f, e := os.Create(path)
			if e != nil {
				return e
			}
			if _, e = io.Copy(f, tarReader); e != nil {
				f.Close()
				return e
			}
			e = f.Chmod(os.FileMode(hdr.Mode))
			if e != nil {
				f.Close()
				return e
			}
			f.Close()
		}
	}
	return nil
}
