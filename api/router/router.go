package router

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/rs/zerolog"
	"gorm.io/gorm"

	"myapp/api/resource/book"
	"myapp/api/resource/health"
	"myapp/api/resource/provider"
	"myapp/api/router/middleware"
	"myapp/api/router/middleware/requestlog"
)

func New(l *zerolog.Logger, v *validator.Validate, db *gorm.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Get("/livez", health.Read)

	r.Route("/v1", func(r chi.Router) {
		r.Use(middleware.RequestID)
		r.Use(middleware.ContentTypeJSON)

		bookAPI := book.New(l, v, db)
		r.Method(http.MethodGet, "/books", requestlog.NewHandler(bookAPI.List, l))
		r.Method(http.MethodPost, "/books", requestlog.NewHandler(bookAPI.Create, l))
		r.Method(http.MethodGet, "/books/{id}", requestlog.NewHandler(bookAPI.Read, l))
		r.Method(http.MethodPut, "/books/{id}", requestlog.NewHandler(bookAPI.Update, l))

		providerAPI := provider.New(l, v, db)
		r.Method(http.MethodGet, "/providers", requestlog.NewHandler(providerAPI.List, l))
		r.Method(http.MethodPost, "/providers", requestlog.NewHandler(providerAPI.Create, l))
		r.Method(http.MethodGet, "/providers/{id}", requestlog.NewHandler(providerAPI.Read, l))
		r.Method(http.MethodPut, "/providers/{id}", requestlog.NewHandler(providerAPI.Update, l))
		r.Method(http.MethodDelete, "/providers/{id}", requestlog.NewHandler(providerAPI.Delete, l))
	})

	return r
}
