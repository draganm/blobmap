package blobmap_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/draganm/blobmap"
	"github.com/stretchr/testify/require"
)

func TestBuilder(t *testing.T) {

	t.Run("size 1 builder", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test1")

		b, err := blobmap.NewBuilder(tempFile, 42, 1)
		require.NoError(t, err)

		err = b.Add(42, []byte("hello"))
		require.NoError(t, err)

		err = b.Build()
		require.NoError(t, err)

		d, err := os.ReadFile(tempFile)
		require.NoError(t, err)

		require.Equal(
			t,
			[]byte{
				0, 0, 0, 0, 0, 0, 0, 1,
				0, 0, 0, 0, 0, 0, 0, 42,
				0, 0, 0, 0, 0, 0, 0, 5,
				104, 101, 108, 108, 111,
			},
			d,
		)
	})

	t.Run("size 2 builder", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test1")

		b, err := blobmap.NewBuilder(tempFile, 42, 2)
		require.NoError(t, err)

		err = b.Add(42, []byte("hello"))
		require.NoError(t, err)

		err = b.Add(43, []byte("world"))
		require.NoError(t, err)

		err = b.Build()
		require.NoError(t, err)

		d, err := os.ReadFile(tempFile)
		require.NoError(t, err)

		require.Equal(
			t,
			[]byte{
				0, 0, 0, 0, 0, 0, 0, 2,
				0, 0, 0, 0, 0, 0, 0, 42,
				0, 0, 0, 0, 0, 0, 0, 5,
				0, 0, 0, 0, 0, 0, 0, 10,
				104, 101, 108, 108, 111,
				119, 111, 114, 108, 100,
			},
			d,
		)
	})

	t.Run("size 4 builder", func(t *testing.T) {
		tempFile := filepath.Join(t.TempDir(), "test1")

		b, err := blobmap.NewBuilder(tempFile, 42, 4)
		require.NoError(t, err)

		err = b.Add(42, []byte("hello"))
		require.NoError(t, err)

		err = b.Add(43, []byte("world"))
		require.NoError(t, err)

		err = b.Add(44, []byte("foo"))
		require.NoError(t, err)

		err = b.Add(45, []byte("bar"))
		require.NoError(t, err)

		err = b.Build()
		require.NoError(t, err)

		d, err := os.ReadFile(tempFile)
		require.NoError(t, err)

		require.Equal(
			t,
			[]byte{
				0, 0, 0, 0, 0, 0, 0, 4,
				0, 0, 0, 0, 0, 0, 0, 42,
				0, 0, 0, 0, 0, 0, 0, 5,
				0, 0, 0, 0, 0, 0, 0, 10,
				0, 0, 0, 0, 0, 0, 0, 13,
				0, 0, 0, 0, 0, 0, 0, 16,
				104, 101, 108, 108, 111,
				119, 111, 114, 108, 100,
				102, 111, 111,
				98, 97, 114,
			},
			d,
		)
	})

}
