package service

import (
	"testing"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/allisson/go-url-shortener/mocks"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

func TestRedirect(t *testing.T) {
	t.Run("Store", func(t *testing.T) {
		redirect := &shortener.Redirect{URL: "https://github.com/allisson"}
		redirectRepository := mocks.RedirectRepository{}
		redirectRepository.On("Store", redirect).Return(nil)
		service := NewRedirect(&redirectRepository)
		err := service.Store(redirect)
		assert.Nil(t, err)
		assert.NotEmpty(t, redirect.Code)
		assert.NotEmpty(t, redirect.CreatedAt)
	})

	t.Run("Store with validation error", func(t *testing.T) {
		redirect := &shortener.Redirect{URL: "invalid-url"}
		redirectRepository := mocks.RedirectRepository{}
		service := NewRedirect(&redirectRepository)
		err := service.Store(redirect)
		assert.Equal(t, "service.Redirect.Store: Redirect Invalid", err.Error())
		assert.Equal(t, shortener.ErrRedirectInvalid, errors.Cause(err))
	})

	t.Run("Store with repository error", func(t *testing.T) {
		redirect := &shortener.Redirect{URL: "https://github.com/allisson"}
		redirectRepository := mocks.RedirectRepository{}
		redirectRepository.On("Store", redirect).Return(errors.New("Repository Error"))
		service := NewRedirect(&redirectRepository)
		err := service.Store(redirect)
		assert.Equal(t, "Repository Error", err.Error())
	})

	t.Run("Find", func(t *testing.T) {
		redirect := &shortener.Redirect{
			Code:      "github-allisson",
			URL:       "https://github.com/allisson",
			CreatedAt: 949407194000,
		}
		redirectRepository := mocks.RedirectRepository{}
		redirectRepository.On("Find", "github-allisson").Return(redirect, nil)
		service := NewRedirect(&redirectRepository)
		redirectResult, err := service.Find("github-allisson")
		assert.Nil(t, err)
		assert.Equal(t, redirect, redirectResult)
	})

	t.Run("Find with invalid code", func(t *testing.T) {
		redirectRepository := mocks.RedirectRepository{}
		redirectRepository.On("Find", "invalid-code").Return(nil, errors.Wrap(shortener.ErrRedirectNotFound, "repository.MockedRedirectRepository.Find"))
		service := NewRedirect(&redirectRepository)
		_, err := service.Find("invalid-code")
		assert.Equal(t, "repository.MockedRedirectRepository.Find: Redirect Not Found", err.Error())
		assert.Equal(t, shortener.ErrRedirectNotFound, errors.Cause(err))
	})
}
