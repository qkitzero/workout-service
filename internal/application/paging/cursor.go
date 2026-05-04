package paging

import (
	"encoding/base64"
	"encoding/json"
)

func EncodeCursor[T any](c T) string {
	b, _ := json.Marshal(c)
	return base64.RawURLEncoding.EncodeToString(b)
}

func DecodeCursor[T any](s string) (T, error) {
	var c T
	b, err := base64.RawURLEncoding.DecodeString(s)
	if err != nil {
		return c, ErrInvalidPageToken
	}
	if err := json.Unmarshal(b, &c); err != nil {
		return c, ErrInvalidPageToken
	}
	return c, nil
}
