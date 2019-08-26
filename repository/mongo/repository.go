package repository

import (
	"context"
	"time"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/pkg/errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

// Redirect is an implementation of shortener.RedirectRepository
type Redirect struct {
	client   *mongo.Client
	database string
	timeout  time.Duration
}

// Find returns a pointer for shortener.Redirect
func (r *Redirect) Find(code string) (*shortener.Redirect, error) {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	redirect := &shortener.Redirect{}
	collection := r.client.Database(r.database).Collection("redirects")
	filter := bson.M{"code": code}
	err := collection.FindOne(ctx, filter).Decode(&redirect)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return redirect, errors.Wrap(shortener.ErrRedirectNotFound, "repository.Redirect.Find")
		}
		return redirect, errors.Wrap(err, "repository.Redirect.Find")
	}
	return redirect, nil
}

// Store receives a pointer of shortener.Redirect and store in redis
func (r *Redirect) Store(redirect *shortener.Redirect) error {
	ctx, cancel := context.WithTimeout(context.Background(), r.timeout)
	defer cancel()
	collection := r.client.Database(r.database).Collection("redirects")
	_, err := collection.InsertOne(
		ctx,
		bson.M{
			"code":       redirect.Code,
			"url":        redirect.URL,
			"created_at": redirect.CreatedAt,
		},
	)
	if err != nil {
		return errors.Wrap(err, "repository.Redirect.Store")
	}
	return nil
}

func newMongoClient(mongoURL, mongoDatabase string, mongoTimeout int) (*mongo.Client, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(mongoTimeout)*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(mongoURL))
	if err != nil {
		return client, err
	}
	err = client.Ping(ctx, readpref.Primary())
	if err != nil {
		return client, err
	}
	return client, err
}

// NewRedirect will create an implementation of shortener.RedirectRepository
func NewRedirect(mongoURL, mongoDatabase string, mongoTimeout int) (*Redirect, error) {
	repo := &Redirect{
		timeout:  time.Duration(mongoTimeout) * time.Second,
		database: mongoDatabase,
	}
	client, err := newMongoClient(mongoURL, mongoDatabase, mongoTimeout)
	if err != nil {
		return repo, errors.Wrap(err, "repository.NewRedirect")
	}
	repo.client = client
	return repo, nil
}
