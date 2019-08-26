package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	h "github.com/allisson/go-url-shortener/http"
	mr "github.com/allisson/go-url-shortener/repository/mongo"
	rr "github.com/allisson/go-url-shortener/repository/redis"
	"github.com/allisson/go-url-shortener/service"
)

func redirectRepository() shortener.RedirectRepository {
	switch os.Getenv("STORAGE_ENGINE") {
	case "redis":
		redisURL := os.Getenv("REDIS_URL")
		repo, err := rr.NewRedirect(redisURL)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	case "mongo":
		mongoURL := os.Getenv("MONGODB_URL")
		mongoDatabase := os.Getenv("MONGODB_DATABASE")
		mongoTimeout, _ := strconv.Atoi(os.Getenv("MONGODB_TIMEOUT"))
		repo, err := mr.NewRedirect(mongoURL, mongoDatabase, mongoTimeout)
		if err != nil {
			log.Fatal(err)
		}
		return repo
	}
	return nil
}

func httpPort() string {
	port := "8000"
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	return fmt.Sprintf(":%s", port)
}

func main() {
	repo := redirectRepository()
	service := service.NewRedirect(repo)
	handler := h.NewHandler(service)
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Get("/{code}", handler.Get)
	r.Post("/", handler.Post)
	log.Fatal(http.ListenAndServe(httpPort(), r))
}
