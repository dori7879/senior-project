package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"

	"api/app/requestlog"
	"api/app/router/middleware"
	"api/app/server"
	"api/config"
)

// New is a function to create the main router.
func New(s *server.Server, conf *config.Conf) *chi.Mux {
	l := s.Logger()

	r := chi.NewRouter()
	// r.Method("GET", "/", requestlog.NewHandler(s.HandleIndex, l))

	r.Route("/api/v1", func(r chi.Router) {
		r.Use(middleware.ContentTypeJSON)
		r.Use(middleware.TokenValidation(conf.AtJwtSecretKey))

		r.Method("POST", "/signup", requestlog.NewHandler(s.HandleSignUp, l))
		r.Method("POST", "/login", requestlog.NewHandler(s.HandleLogin, l))
		r.Method("POST", "/refresh-token", requestlog.NewHandler(s.HandleTokenRefresh, l))

		// Routes for homework pages
		r.Method("GET", "/homework-page", requestlog.NewHandler(s.HandleListHomeworkPage, l))
		r.Method("POST", "/homework-page", requestlog.NewHandler(s.HandleCreateHomeworkPage, l))
		r.Method("GET", "/homework-page/{id}", requestlog.NewHandler(s.HandleReadHomeworkPage, l))
		r.Method("PUT", "/homework-page/{id}", requestlog.NewHandler(s.HandleUpdateHomeworkPage, l))
		r.Method("DELETE", "/homework-page/{id}", requestlog.NewHandler(s.HandleDeleteHomeworkPage, l))

		r.Method("GET", "/homework-page/link/{str}", requestlog.NewHandler(s.HandleReadHomeworkPageByStringParam, l))

		// Routes for homeworks
		r.Method("GET", "/homework", requestlog.NewHandler(s.HandleListHomework, l))
		r.Method("POST", "/homework", requestlog.NewHandler(s.HandleCreateHomework, l))
		r.Method("GET", "/homework/{id}", requestlog.NewHandler(s.HandleReadHomework, l))
		r.Method("PUT", "/homework/{id}", requestlog.NewHandler(s.HandleUpdateHomework, l))
		r.Method("DELETE", "/homework/{id}", requestlog.NewHandler(s.HandleDeleteHomework, l))
	})

	// Routes for healthz
	r.Get("/healthz/liveness", server.HandleLive)
	r.Method("GET", "/healthz/readiness", requestlog.NewHandler(s.HandleReady, l))

	FileServer(r)

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
