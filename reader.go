package blobmap

import (
	"fmt"
	"os"

	"github.com/tysonmote/gommap"
)

type Reader struct {
	f        *os.File
	mm       gommap.MMap
	layout   layout
	firstKey uint64
	lastKey  uint64
}

func Open(fileName string) (*Reader, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	mm, err := gommap.Map(f.Fd(), gommap.PROT_READ, gommap.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("failed to mmap file: %w", err)
	}

	header := header(mm[:16])

	reader := &Reader{
		f:  f,
		mm: mm,
		layout: layout{
			header:    header,
			positions: positions(mm[16 : 16+header.size()*8]),
			data:      data(mm[16+header.size()*8 : len(mm)-8]),
		},
		firstKey: header.firstKey(),
		lastKey:  header.firstKey() + header.size() - 1,
	}

	return reader, nil
}

func (r *Reader) Close() error {
	err := r.mm.UnsafeUnmap()
	if err != nil {
		return fmt.Errorf("failed to unmap file: %w", err)
	}

	err = r.f.Close()

	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}

func (r *Reader) Read(key uint64) ([]byte, error) {

	if key < r.firstKey || key > r.lastKey {
		return nil, fmt.Errorf("key %d out of range [%d, %d]", key, r.firstKey, r.lastKey)
	}

	d := r.layout.getData(key - r.firstKey)

	return d, nil
}

func (r *Reader) FirstKey() uint64 {
	return r.firstKey
}

func (r *Reader) LastKey() uint64 {
	return r.lastKey
}
