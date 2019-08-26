package http

import (
	"io/ioutil"
	"log"
	"net/http"

	shortener "github.com/allisson/go-url-shortener"
	js "github.com/allisson/go-url-shortener/serializer/json"
	ms "github.com/allisson/go-url-shortener/serializer/msgpack"
	"github.com/go-chi/chi"
	"github.com/pkg/errors"
)

func makeResponse(w http.ResponseWriter, contentType string, body []byte, statusCode int) {
	w.Header().Set("Content-Type", contentType)
	w.WriteHeader(statusCode)
	_, err := w.Write(body)
	if err != nil {
		log.Println(err)
	}
}

// Handler ...
type Handler struct {
	redirectService shortener.RedirectService
}

func (h *Handler) serializer(contentType string) shortener.RedirectSerializer {
	switch contentType {
	case "application/x-msgpack":
		return &ms.Redirect{}
	case "application/json":
		return &js.Redirect{}
	default:
		return &js.Redirect{}
	}
}

// Get handler
func (h *Handler) Get(w http.ResponseWriter, r *http.Request) {
	code := chi.URLParam(r, "code")
	redirect, err := h.redirectService.Find(code)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectNotFound {
			http.Error(w, http.StatusText(http.StatusNotFound), http.StatusNotFound)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, redirect.URL, http.StatusMovedPermanently)
}

// Post handler
func (h *Handler) Post(w http.ResponseWriter, r *http.Request) {
	contentType := r.Header.Get("Content-Type")
	requestBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	redirect, err := h.serializer(contentType).Encode(requestBody)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	err = h.redirectService.Store(redirect)
	if err != nil {
		if errors.Cause(err) == shortener.ErrRedirectInvalid {
			http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
			return
		}
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	responseBody, err := h.serializer(contentType).Decode(redirect)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}
	makeResponse(w, contentType, responseBody, http.StatusCreated)
}

// NewHandler ...
func NewHandler(redirectService shortener.RedirectService) *Handler {
	return &Handler{redirectService: redirectService}
}
