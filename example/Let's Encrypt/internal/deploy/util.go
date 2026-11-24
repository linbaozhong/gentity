package deploy

import (
	"bytes"
	"encoding/json"
	"io"
	"time"
)

func unixMilli() int64 {
	return time.Now().UnixMilli()
}

func bytesReader(data []byte) io.Reader {
	return bytes.NewReader(data)
}

func jsonUnmarshalECS(data []byte, v any) error {
	if err := json.Unmarshal(data, v); err != nil {
		return err
	}
	return nil
}
