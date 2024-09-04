package blobmap

import "encoding/binary"

type header []byte

func (h header) size() uint64 {
	return binary.BigEndian.Uint64(h)
}

func (h header) firstKey() uint64 {
	return binary.BigEndian.Uint64(h[8:])
}

func (h header) setFirstKey(k uint64) {
	binary.BigEndian.PutUint64(h[8:], k)
}

func (h header) setNumberOfKeys(k uint64) {
	binary.BigEndian.PutUint64(h, k)

}

type positions []byte

func (p positions) startOf(i uint64) uint64 {
	if i == 0 {
		return 0
	}
	return binary.BigEndian.Uint64(p[(i-1)*8:])
}

func (p positions) endOf(i uint64) uint64 {
	return binary.BigEndian.Uint64(p[i*8:])
}

func (p positions) setEndOf(i, v uint64) {
	binary.BigEndian.PutUint64(p[i*8:], v)
}

type data []byte

type layout struct {
	header    header
	positions positions
	data      data
}

func (l layout) getPrefixSize() uint64 {
	return uint64(len(l.header) + len(l.positions))
}

func (l layout) getData(key uint64) []byte {
	start := l.positions.startOf(key)
	end := l.positions.endOf(key)
	return l.data[start:end]
}

func (l layout) totalSize() uint64 {
	return l.positions.endOf(l.header.size()-1) + l.getPrefixSize()
}
