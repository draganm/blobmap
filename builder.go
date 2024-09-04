package blobmap

import (
	"fmt"
	"os"
	"sync"

	"github.com/tysonmote/gommap"
)

type Builder struct {
	f        *os.File
	mm       gommap.MMap
	mu       *sync.Mutex
	firstKey uint64
	nextKey  uint64
	lastKey  uint64
	capacity uint64
	layout   layout
}

func NewBuilder(fileName string, firstKey, numberOfKeys uint64) (*Builder, error) {

	if numberOfKeys == 0 {
		return nil, fmt.Errorf("numberOfKeys must be greater than 0")
	}

	initialSize := numberOfKeys*8 + 16

	f, err := os.OpenFile(fileName, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	err = f.Truncate(int64(initialSize))
	if err != nil {
		return nil, fmt.Errorf("failed to truncate file: %w", err)
	}

	mm, err := gommap.Map(f.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
	if err != nil {
		return nil, fmt.Errorf("failed to mmap file: %w", err)
	}

	l := layout{
		header:    header(mm[:16]),
		positions: positions(mm[16 : numberOfKeys*8+16]),
		data:      data(mm[numberOfKeys*8+16:]),
	}

	l.header.setFirstKey(firstKey)
	l.header.setNumberOfKeys(numberOfKeys)

	b := &Builder{
		f:        f,
		mm:       mm,
		mu:       &sync.Mutex{},
		firstKey: firstKey,
		nextKey:  firstKey,
		lastKey:  firstKey + numberOfKeys,
		capacity: numberOfKeys,
		layout:   l,
	}

	return b, nil
}

func (b *Builder) Build() error {

	b.mu.Lock()
	defer b.mu.Unlock()

	if b.nextKey != b.lastKey {
		return fmt.Errorf("not all keys have been added, expected last key %d, got %d", b.lastKey, b.nextKey)
	}

	totalSize := b.layout.totalSize()

	err := b.mm.UnsafeUnmap()
	if err != nil {
		return fmt.Errorf("failed to unmap file: %w", err)
	}

	err = b.f.Truncate(int64(totalSize))
	if err != nil {
		return fmt.Errorf("failed to truncate file: %w", err)
	}

	err = b.f.Close()
	if err != nil {
		return fmt.Errorf("failed to close file: %w", err)
	}

	return nil
}

func (b *Builder) Add(key uint64, value []byte) error {
	b.mu.Lock()
	defer b.mu.Unlock()

	if key >= b.lastKey {
		return fmt.Errorf("key out of range")
	}

	if key != b.nextKey {
		return fmt.Errorf("expected key %d, got %d", b.nextKey, key)
	}

	relKey := key - b.firstKey

	startOffset := b.layout.positions.startOf(relKey)
	endOffset := startOffset + uint64(len(value))

	b.layout.positions.setEndOf(relKey, endOffset)

	err := b.ensureSize(endOffset + b.layout.getPrefixSize())
	if err != nil {
		return fmt.Errorf("failed to grow file: %w", err)
	}

	copy(b.layout.getData(relKey), value)

	b.nextKey++

	return nil
}

func (b *Builder) ensureSize(size uint64) error {
	if size > uint64(len(b.mm)) {
		err := b.mm.UnsafeUnmap()
		if err != nil {
			return fmt.Errorf("failed to unmap file: %w", err)
		}

		// Grow the file by 50% over the required size
		err = b.f.Truncate(int64((size * 3) / 2))
		if err != nil {
			return fmt.Errorf("failed to truncate file: %w", err)
		}

		mm, err := gommap.Map(b.f.Fd(), gommap.PROT_READ|gommap.PROT_WRITE, gommap.MAP_SHARED)
		if err != nil {
			return fmt.Errorf("failed to mmap file: %w", err)
		}

		b.mm = mm

		h := header(b.mm[:16])

		l := layout{
			header:    h,
			positions: positions(mm[16 : h.size()*8+16]),
			data:      data(mm[h.size()*8+16:]),
		}

		b.layout = l

	}

	return nil
}
