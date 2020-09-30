package router

import (
	"github.com/go-chi/chi"

	"api/app/requestlog"
	"api/app/server"
)

func New(s *server.Server) *chi.Mux {
	l := s.Logger()

	r := chi.NewRouter()
	r.Method("GET", "/", requestlog.NewHandler(s.HandleIndex, l))

	return r
}
