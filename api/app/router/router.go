package router

import (
	"net/http"
	"os"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"

	"api/app/requestlog"
	"api/app/router/middleware"
	"api/app/server"
	"api/config"
)

// New is a function to create the main router.
func New(s *server.Server, conf *config.Conf) *chi.Mux {
	l := s.Logger()

	r := chi.NewRouter()

	r.Route("/api/v1", func(r chi.Router) {
		// Basic CORS
		// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
		r.Use(cors.Handler(cors.Options{
			// AllowedOrigins: []string{"https://foo.com"}, // Use this to allow specific origin hosts
			AllowedOrigins: []string{"*"},
			// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
			AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
			AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
			ExposedHeaders:   []string{"Link"},
			AllowCredentials: false,
			MaxAge:           300, // Maximum value not ignored by any of major browsers
		}))

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

		r.Method("GET", "/homework-page/student/{str}", requestlog.NewHandler(s.HandleReadHomeworkPageByStudentLink, l))
		r.Method("GET", "/homework-page/teacher/{str}", requestlog.NewHandler(s.HandleReadHomeworkPageByTeacherLink, l))

		// Routes for homeworks
		r.Method("GET", "/homework", requestlog.NewHandler(s.HandleListHomework, l))
		r.Method("POST", "/homework", requestlog.NewHandler(s.HandleCreateHomework, l))
		r.Method("GET", "/homework/{id}", requestlog.NewHandler(s.HandleReadHomework, l))
		r.Method("PUT", "/homework/{id}", requestlog.NewHandler(s.HandleUpdateHomework, l))
		r.Method("DELETE", "/homework/{id}", requestlog.NewHandler(s.HandleDeleteHomework, l))

		// Routes for quizzes
		r.Method("GET", "/quiz", requestlog.NewHandler(s.HandleListQuiz, l))
		r.Method("POST", "/quiz", requestlog.NewHandler(s.HandleCreateQuiz, l))
		r.Method("GET", "/quiz/{id}", requestlog.NewHandler(s.HandleReadQuiz, l))
		r.Method("PUT", "/quiz/{id}", requestlog.NewHandler(s.HandleUpdateQuiz, l))
		r.Method("DELETE", "/quiz/{id}", requestlog.NewHandler(s.HandleDeleteQuiz, l))

		r.Method("GET", "/quiz/student/{str}", requestlog.NewHandler(s.HandleReadQuizByStudentLink, l))
		r.Method("GET", "/quiz/teacher/{str}", requestlog.NewHandler(s.HandleReadQuizByTeacherLink, l))

		// Routes for quiz submissions
		r.Method("GET", "/quiz/submission", requestlog.NewHandler(s.HandleListQuizSubmission, l))
		r.Method("POST", "/quiz/submission", requestlog.NewHandler(s.HandleCreateQuizSubmission, l))
		r.Method("GET", "/quiz/submission/{id}", requestlog.NewHandler(s.HandleReadQuizSubmission, l))
		r.Method("PUT", "/quiz/submission/{id}", requestlog.NewHandler(s.HandleUpdateQuizSubmission, l))
		r.Method("DELETE", "/quiz/submission/{id}", requestlog.NewHandler(s.HandleDeleteQuizSubmission, l))
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
