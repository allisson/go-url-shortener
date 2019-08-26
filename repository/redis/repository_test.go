package repository

import (
	"os"
	"testing"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	redirect := &shortener.Redirect{
		Code:      "github-allisson",
		URL:       "https://github.com/allisson",
		CreatedAt: 949407194000,
	}
	redisURL := os.Getenv("REDIS_URL")
	repo, err := NewRedirect(redisURL)
	assert.Nil(t, err)

	t.Run("Store", func(t *testing.T) {
		err := repo.Store(redirect)
		assert.Nil(t, err)
	})

	t.Run("Find", func(t *testing.T) {
		redirectResult, err := repo.Find("github-allisson")
		assert.Nil(t, err)
		assert.Equal(t, redirectResult, redirect)
	})

	t.Run("Find invalid code", func(t *testing.T) {
		_, err := repo.Find("invalid-code")
		assert.Equal(t, "repository.Redirect.Find: Redirect Not Found", err.Error())
	})
}
