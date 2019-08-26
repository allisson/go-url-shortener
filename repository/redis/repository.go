package repository

import (
	"fmt"
	"strconv"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/go-redis/redis"
	"github.com/pkg/errors"
)

// Redirect is an implementation of shortener.RedirectRepository
type Redirect struct {
	client *redis.Client
}

func (r *Redirect) generateKey(code string) string {
	return fmt.Sprintf("redirect:%s", code)
}

// Find returns a pointer for shortener.Redirect
func (r *Redirect) Find(code string) (*shortener.Redirect, error) {
	redirect := &shortener.Redirect{}
	key := r.generateKey(code)
	data, err := r.client.HGetAll(key).Result()
	if err != nil {
		return redirect, errors.Wrap(err, "repository.Redirect.Find")
	}
	if len(data) == 0 {
		return redirect, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
	}
	createdAt, err := strconv.ParseInt(data["created_at"], 10, 64)
	if err != nil {
		return redirect, errors.Wrap(err, "repository.Redirect.Find")
	}
	redirect.Code = data["code"]
	redirect.URL = data["url"]
	redirect.CreatedAt = createdAt
	return redirect, err
}

// Store receives a pointer of shortener.Redirect and store in redis
func (r *Redirect) Store(redirect *shortener.Redirect) error {
	key := r.generateKey(redirect.Code)
	data := map[string]interface{}{
		"code":       redirect.Code,
		"url":        redirect.URL,
		"created_at": redirect.CreatedAt,
	}
	_, err := r.client.HMSet(key, data).Result()
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

func newRedisClient(redisURL string) (*redis.Client, error) {
	options, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, err
	}
	client := redis.NewClient(options)
	_, err = client.Ping().Result()
	return client, err
}

// NewRedirect will create an implementation of shortener.RedirectRepository
func NewRedirect(redisURL string) (*Redirect, error) {
	repo := &Redirect{}
	client, err := newRedisClient(redisURL)
	if err != nil {
		return repo, errors.Wrap(err, "repository.NewRedirect")
	}
	repo.client = client
	return repo, nil
}
