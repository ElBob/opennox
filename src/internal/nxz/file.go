package nxz

import (
	"errors"
	"io"

	"github.com/opennox/libs/ifs"
)

func readFile(path string) ([]byte, error) {
	f, err := ifs.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	return io.ReadAll(f)
}

func writeFile(path string, data []byte) error {
	f, err := ifs.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()
	if _, err = f.Write(data); err != nil {
		return err
	}
	return f.Close()
}

// CompressFile compresses src file to dst in nxz format.
func CompressFile(src, dst string) error {
	if src == "" {
		return errors.New("empty source path")
	}
	if dst == "" {
		return errors.New("empty destination path")
	}
	data, err := readFile(src)
	if err != nil {
		return err
	}
	return writeFile(dst, Compress(data))
}

// DecompressFile decompresses nxz file src to dst.
func DecompressFile(src, dst string) error {
	if src == "" {
		return errors.New("empty source path")
	}
	if dst == "" {
		return errors.New("empty destination path")
	}
	data, err := readFile(src)
	if err != nil {
		return err
	}
	data, err = Decompress(data)
	if err != nil {
		return err
	}
	return writeFile(dst, data)
}
