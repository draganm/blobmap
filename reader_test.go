package blobmap_test

import (
	"testing"

	"github.com/draganm/blobmap"
	"github.com/stretchr/testify/require"
)

func TestReader(t *testing.T) {
	t.Run("read single blob map", func(t *testing.T) {
		r, err := blobmap.Open("testdata/test1.blobmap")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, uint64(42), r.FirstKey())
		require.Equal(t, uint64(42), r.LastKey())

		data, err := r.Read(42)
		require.NoError(t, err)

		require.Equal(t, []byte("hello"), data)
	})

	t.Run("read two blobs map", func(t *testing.T) {
		r, err := blobmap.Open("testdata/test2.blobmap")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, uint64(42), r.FirstKey())
		require.Equal(t, uint64(43), r.LastKey())

		data, err := r.Read(42)
		require.NoError(t, err)

		require.Equal(t, []byte("hello"), data)

		data, err = r.Read(43)
		require.NoError(t, err)

		require.Equal(t, []byte("world"), data)

	})

	t.Run("read four blobs map", func(t *testing.T) {
		r, err := blobmap.Open("testdata/test4.blobmap")
		require.NoError(t, err)
		defer r.Close()

		require.Equal(t, uint64(42), r.FirstKey())
		require.Equal(t, uint64(45), r.LastKey())

		data, err := r.Read(42)
		require.NoError(t, err)

		require.Equal(t, []byte("hello"), data)

		data, err = r.Read(43)
		require.NoError(t, err)

		require.Equal(t, []byte("world"), data)

		data, err = r.Read(44)
		require.NoError(t, err)

		require.Equal(t, []byte("foo"), data)

		data, err = r.Read(45)
		require.NoError(t, err)

		require.Equal(t, []byte("bar"), data)

	})

}
