package buffer

import (
	"bytes"
	"io"
)

func Read(reader io.ReadCloser) (*bytes.Buffer, error) {
	defer reader.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return nil, err
	}

	return buf, nil
}
