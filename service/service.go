package service

import (
	"time"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/pkg/errors"
	"github.com/teris-io/shortid"
	"gopkg.in/dealancer/validate.v2"
)

// Redirect is an implementation of shortener.RedirectService
type Redirect struct {
	redirectRepository shortener.RedirectRepository
}

// Find ...
func (r *Redirect) Find(code string) (*shortener.Redirect, error) {
	return r.redirectRepository.Find(code)
}

// Store ...
func (r *Redirect) Store(redirect *shortener.Redirect) error {
	if err := validate.Validate(redirect); err != nil {
		return errors.Wrap(shortener.ErrRedirectInvalid, "service.Redirect.Store")
	}
	redirect.Code = shortid.MustGenerate()
	redirect.CreatedAt = time.Now().UTC().Unix()
	return r.redirectRepository.Store(redirect)
}

// NewRedirect will create an implementation of shortener.RedirectService
func NewRedirect(redirectRepository shortener.RedirectRepository) *Redirect {
	return &Redirect{redirectRepository: redirectRepository}
}
