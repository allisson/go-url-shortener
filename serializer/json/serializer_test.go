package serializer

import (
	"encoding/json"
	"testing"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	serializer := Redirect{}

	t.Run("Encode with invalid input", func(t *testing.T) {
		invalidRawMsg := []byte("msg=strogonoficamente-sensivel-message")
		_, err := serializer.Encode(invalidRawMsg)
		assert.Equal(t, "serializer.Redirect.Encode: invalid character 'm' looking for beginning of value", err.Error())
	})

	t.Run("Encode with valid input", func(t *testing.T) {
		redirect := &shortener.Redirect{
			Code: "github-allisson",
			URL:  "https://github.com/allisson",
		}
		rawMsg, err := json.Marshal(redirect)
		assert.Nil(t, err)
		redirectResult, err := serializer.Encode(rawMsg)
		assert.Nil(t, err)
		assert.Equal(t, redirect, redirectResult)
	})

	t.Run("Decode with valid input", func(t *testing.T) {
		redirect := &shortener.Redirect{
			Code:      "github-allisson",
			URL:       "https://github.com/allisson",
			CreatedAt: 949407194000,
		}
		rawMsg, err := json.Marshal(redirect)
		assert.Nil(t, err)
		rawMsgResult, err := serializer.Decode(redirect)
		assert.Nil(t, err)
		assert.Equal(t, rawMsgResult, rawMsg)
	})
}
