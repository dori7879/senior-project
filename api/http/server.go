package http

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	_ "net/http/pprof"
	"strings"
	"time"

	"github.com/dori7879/senior-project/api"
	"github.com/dori7879/senior-project/api/jwt"
	"github.com/gorilla/mux"
)

// ShutdownTimeout is the time given for outstanding requests to finish before shutdown.
const ShutdownTimeout = 1 * time.Second

// Server represents an HTTP server. It is meant to wrap all HTTP functionality
// used by the application so that dependent packages (such as cmd/wtfd) do not
// need to reference the "net/http" package at all.
type Server struct {
	ln     net.Listener
	server *http.Server
	router *mux.Router

	// Bind address & domain for the server's listener.
	// If domain is specified, server is run on TLS using acme/autocert.
	Addr   string
	Domain string

	// Servics used by the various HTTP routes.
	AuthService           jwt.AuthService
	UserService           api.UserService
	GroupService          api.GroupService
	HomeworkService       api.HomeworkService
	HWSubmissionService   api.HWSubmissionService
	QuizService           api.QuizService
	QuizSubmissionService api.QuizSubmissionService
	QuestionService       api.QuestionService
	ResponseService       api.ResponseService
}

// NewServer returns a new instance of Server.
func NewServer() *Server {
	// Create a new server that wraps the net/http server & add a gorilla router.
	s := &Server{
		server: &http.Server{},
		router: mux.NewRouter(),
	}

	// Report panics to external service.
	s.router.Use(reportPanic)

	// Our router is wrapped by another function handler to perform some
	// middleware-like tasks that cannot be performed by actual middleware.
	// This includes changing route paths for JSON endpoints & overridding methods.
	s.server.Handler = http.HandlerFunc(s.serveHTTP)

	// Setup error handling routes.
	s.router.NotFoundHandler = http.HandlerFunc(s.handleNotFound)

	// Setup a base router that excludes asset handling.
	router := s.router.PathPrefix("/").Subrouter()
	router.Use(s.authenticate)

	// Register unauthenticated routes.
	{
		r := s.router.PathPrefix("/").Subrouter()
		r.Use(s.requireNoAuth)
	}

	// Register authenticated routes.
	{
		r := router.PathPrefix("/").Subrouter()
		r.Use(s.requireAuth)
	}

	return s
}

// UseTLS returns true if the cert & key file are specified.
func (s *Server) UseTLS() bool {
	return s.Domain != ""
}

// Scheme returns the URL scheme for the server.
func (s *Server) Scheme() string {
	if s.UseTLS() {
		return "https"
	}
	return "http"
}

// Port returns the TCP port for the running server.
// This is useful in tests where we allocate a random port by using ":0".
func (s *Server) Port() int {
	if s.ln == nil {
		return 0
	}
	return s.ln.Addr().(*net.TCPAddr).Port
}

// URL returns the local base URL of the running server.
func (s *Server) URL() string {
	scheme, port := s.Scheme(), s.Port()

	// Use localhost unless a domain is specified.
	domain := "localhost"
	if s.Domain != "" {
		domain = s.Domain
	}

	// Return without port if using standard ports.
	if (scheme == "http" && port == 80) || (scheme == "https" && port == 443) {
		return fmt.Sprintf("%s://%s", s.Scheme(), domain)
	}
	return fmt.Sprintf("%s://%s:%d", s.Scheme(), domain, s.Port())
}

// Open validates the server options and begins listening on the bind address.
func (s *Server) Open() (err error) {
	// Open a listener on our bind address.
	if s.ln, err = net.Listen("tcp", s.Addr); err != nil {
		return err
	}

	// Begin serving requests on the listener. We use Serve() instead of
	// ListenAndServe() because it allows us to check for listen errors (such
	// as trying to use an already open port) synchronously.
	go s.server.Serve(s.ln)

	return nil
}

// Close gracefully shuts down the server.
func (s *Server) Close() error {
	ctx, cancel := context.WithTimeout(context.Background(), ShutdownTimeout)
	defer cancel()
	return s.server.Shutdown(ctx)
}

func (s *Server) serveHTTP(w http.ResponseWriter, r *http.Request) {
	// Delegate remaining HTTP handling to the gorilla router.
	s.router.ServeHTTP(w, r)
}

// authenticate is middleware for loading session data from a cookie or API key header.
func (s *Server) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Login via API key, if available.
		if v := r.Header.Get("Authorization"); strings.HasPrefix(v, "Bearer ") {
			tokenStr := strings.TrimPrefix(v, "Bearer ")

			userID, err := s.AuthService.Validate(tokenStr)
			if err != nil {
				Error(w, r, api.Errorf(api.EINTERNAL, "Token corrupted"))
				return
			} else if userID == 0 {
				// Delegate to next HTTP handler.
				next.ServeHTTP(w, r)
			}

			if user, err := s.UserService.FindUserByID(r.Context(), userID); err != nil {
				log.Printf("cannot find session user: id=%d err=%s", userID, err)
			} else {
				// Update request context to include authenticated user.
				r = r.WithContext(api.NewContextWithUser(r.Context(), user))
			}
		}

		// Delegate to next HTTP handler.
		next.ServeHTTP(w, r)
	})
}

// requireNoAuth is middleware for requiring no authentication.
// This is used if a user goes to log in but is already logged in.
func (s *Server) requireNoAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If user is logged in, redirect to the home page.
		if userID := api.UserIDFromContext(r.Context()); userID != 0 {
			w.Header().Set("Content-type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			json.NewEncoder(w).Encode(&ErrorResponse{Error: "Only for guests"})
			return
		}

		// Delegate to next HTTP handler.
		next.ServeHTTP(w, r)
	})
}

// requireAuth is middleware for requiring authentication. This is used by
// nearly every page except for the login & oauth pages.
func (s *Server) requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// If user is logged in, delegate to next HTTP handler.
		if userID := api.UserIDFromContext(r.Context()); userID != 0 {
			next.ServeHTTP(w, r)
			return
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		json.NewEncoder(w).Encode(&ErrorResponse{Error: "You are not authorized"})
	})
}

// reportPanic is middleware for catching panics and reporting them.
func reportPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				api.ReportPanic(err)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

// handleNotFound handles requests to routes that don't exist.
func (s *Server) handleNotFound(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(&ErrorResponse{Error: "Not Found"})
}
