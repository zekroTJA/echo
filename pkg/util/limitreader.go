package util

import (
	"errors"
	"fmt"
	"io"
)

var ErrLimitReached = errors.New("payload size limit exceeded")

type LimitReader struct {
	Reader io.Reader
	Limit  int
}

func (t *LimitReader) Read(p []byte) (int, error) {
	if t.Limit < 0 {
		return 0, ErrLimitReached
	}

	fmt.Println(t.Limit)

	if len(p) > t.Limit+1 {
		p = p[:t.Limit+1]
	}

	i, err := t.Reader.Read(p)
	if err != nil {
		return i, err
	}

	t.Limit -= i

	return i, nil
}
