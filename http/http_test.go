package http

import (
	"bytes"
	"context"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	shortener "github.com/allisson/go-url-shortener"
	"github.com/allisson/go-url-shortener/mocks"
	js "github.com/allisson/go-url-shortener/serializer/json"
	ms "github.com/allisson/go-url-shortener/serializer/msgpack"
	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler(t *testing.T) {
	redirect := &shortener.Redirect{
		Code:      "github-allisson",
		URL:       "https://github.com/allisson",
		CreatedAt: 949407194000,
	}

	t.Run("Get with invalid code", func(t *testing.T) {
		redirectService := mocks.RedirectService{}
		redirectService.On("Find", "invalid-code").Return(nil, shortener.ErrRedirectNotFound)
		handler := NewHandler(&redirectService)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/invalid-code", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("code", "invalid-code")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		handler.Get(rr, req)
		assert.Equal(t, http.StatusNotFound, rr.Code)
	})

	t.Run("Get with valid code", func(t *testing.T) {
		redirectService := mocks.RedirectService{}
		redirectService.On("Find", "github-allisson").Return(redirect, nil)
		handler := NewHandler(&redirectService)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("GET", "/github-allisson", nil)
		rctx := chi.NewRouteContext()
		rctx.URLParams.Add("code", "github-allisson")
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		handler.Get(rr, req)
		assert.Equal(t, http.StatusMovedPermanently, rr.Code)
		assert.Equal(t, redirect.URL, rr.Header().Get("Location"))
	})

	t.Run("Post with json", func(t *testing.T) {
		serializer := js.Redirect{}
		body, _ := serializer.Encode(redirect)
		redirectService := mocks.RedirectService{}
		redirectService.On("Store", redirect).Return(nil)
		handler := NewHandler(&redirectService)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req.Header.Set("Content-Type", "application/json")
		handler.Post(rr, req)
		responseBody, _ := ioutil.ReadAll(rr.Body)
		assert.Equal(t, "application/json", rr.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, body, responseBody)
	})

	t.Run("Post with msgpack", func(t *testing.T) {
		serializer := ms.Redirect{}
		body, _ := serializer.Encode(redirect)
		redirectService := mocks.RedirectService{}
		redirectService.On("Store", redirect).Return(nil)
		handler := NewHandler(&redirectService)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req.Header.Set("Content-Type", "application/x-msgpack")
		handler.Post(rr, req)
		responseBody, _ := ioutil.ReadAll(rr.Body)
		assert.Equal(t, "application/x-msgpack", rr.Header().Get("Content-Type"))
		assert.Equal(t, http.StatusCreated, rr.Code)
		assert.Equal(t, body, responseBody)
	})

	t.Run("Post with wrong content type", func(t *testing.T) {
		serializer := ms.Redirect{}
		body, _ := serializer.Encode(redirect) // msgpack body
		redirectService := mocks.RedirectService{}
		redirectService.On("Store", redirect).Return(nil)
		handler := NewHandler(&redirectService)
		rr := httptest.NewRecorder()
		req := httptest.NewRequest("POST", "/", bytes.NewReader(body))
		rctx := chi.NewRouteContext()
		req = req.WithContext(context.WithValue(req.Context(), chi.RouteCtxKey, rctx))
		req.Header.Set("Content-Type", "application/json")
		handler.Post(rr, req)
		assert.Equal(t, http.StatusInternalServerError, rr.Code)
	})
}
