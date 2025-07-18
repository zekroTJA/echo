package util

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func getZeroBytes(n int) []byte {
	return bytes.Repeat([]byte{0}, n)[:n]
}

func TestLimitReader_Read(t *testing.T) {
	t.Run("under-limit", func(t *testing.T) {
		buf := bytes.NewReader(getZeroBytes(600))

		r := &LimitReader{
			Reader: buf,
			Limit:  1000,
		}

		read, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, getZeroBytes(600), read)
	})

	t.Run("over-limit", func(t *testing.T) {
		buf := bytes.NewReader(getZeroBytes(600))

		r := &LimitReader{
			Reader: buf,
			Limit:  300,
		}

		_, err := io.ReadAll(r)
		assert.ErrorIs(t, err, ErrLimitReached)
		// assert.Len(t, read, 300)
	})

	t.Run("exact-limit", func(t *testing.T) {
		buf := bytes.NewReader(getZeroBytes(300))

		r := &LimitReader{
			Reader: buf,
			Limit:  300,
		}

		read, err := io.ReadAll(r)
		assert.Nil(t, err)
		assert.Equal(t, getZeroBytes(300), read)
	})

	t.Run("small-over-limit", func(t *testing.T) {
		buf := bytes.NewReader(getZeroBytes(301))

		r := &LimitReader{
			Reader: buf,
			Limit:  300,
		}

		_, err := io.ReadAll(r)
		assert.ErrorIs(t, err, ErrLimitReached)
	})
}
