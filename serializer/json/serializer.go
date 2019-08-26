package serializer

import (
	"encoding/json"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/pkg/errors"
)

// Redirect is an implementation of shortener.Encoder
type Redirect struct{}

// Encode receives json message in bytes and convert to pointer of shortener.Redirect
func (r *Redirect) Encode(input []byte) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	if err := json.Unmarshal(input, redirect); err != nil {
		return redirect, errors.Wrap(err, "serializer.Redirect.Encode")
	}
	return redirect, nil
}

// Decode receives a pointer of shortener.Redirect and returns json message in bytes
func (r *Redirect) Decode(input *shortener.Redirect) ([]byte, error) {
	rawMsg, err := json.Marshal(input)
	if err != nil {
		return rawMsg, errors.Wrap(err, "serializer.Redirect.Decode")
	}
	return rawMsg, nil
}
