package tsutsu

import (
	"encoding/json"
	"io"
)

type httpBodyDecoder struct {
	body    io.ReadCloser
	decoder *json.Decoder
}

func newHttpBodyDecoder(body io.ReadCloser) *httpBodyDecoder {
	decoder := json.NewDecoder(body)
	return &httpBodyDecoder{
		body:    body,
		decoder: decoder,
	}
}

func (h *httpBodyDecoder) Close() error {
	return h.body.Close()
}

func (h *httpBodyDecoder) Decode(v interface{}) error {
	return h.decoder.Decode(v)
}
