package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"api/app/requestlog"
	"api/app/router/middleware"
	"api/app/server"
)

// New is a function to create the main router.
func New(s *server.Server) *chi.Mux {
	l := s.Logger()

	r := chi.NewRouter()
	r.Method("GET", "/", requestlog.NewHandler(s.HandleIndex, l))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJson)

		// Routes for homework pages
		r.Method("GET", "/homework", requestlog.NewHandler(s.HandleListHomework, l))
		r.Method("POST", "/homework", requestlog.NewHandler(s.HandleCreateHomework, l))
		r.Method("GET", "/homework/{id}", requestlog.NewHandler(s.HandleReadHomework, l))
		r.Method("PUT", "/homework/{id}", requestlog.NewHandler(s.HandleUpdateHomework, l))
		r.Method("DELETE", "/homework/{id}", requestlog.NewHandler(s.HandleDeleteHomework, l))
	})

	FileServer(r)

	// Routes for healthz
	r.Get("/healthz/liveness", server.HandleLive)
	r.Method("GET", "/healthz/readiness", requestlog.NewHandler(s.HandleReady, l))

	return r
}

// FileServer is serving static files.
func FileServer(router *chi.Mux) {
	root := "./web"
	fs := http.FileServer(http.Dir(root))

	router.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		if _, err := os.Stat(root + r.RequestURI); os.IsNotExist(err) {
			http.StripPrefix(r.RequestURI, fs).ServeHTTP(w, r)
		} else {
			fs.ServeHTTP(w, r)
		}
	})
}
