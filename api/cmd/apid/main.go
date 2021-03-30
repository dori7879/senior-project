package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"runtime/debug"

	"github.com/dori7879/senior-project/api"
	"github.com/dori7879/senior-project/api/http"
	"github.com/dori7879/senior-project/api/jwt"
	"github.com/dori7879/senior-project/api/pg"
)

// main is the entry point to our application binary. However, it has some poor
// usability so we mainly use it to delegate out to our Main type.
func main() {
	// Setup signal handlers.
	ctx, cancel := context.WithCancel(context.Background())
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() { <-c; cancel() }()

	// Instantiate a new type to represent our application.
	// This type lets us shared setup code with our end-to-end tests.
	m := NewMain()

	// Parse command line flags & load configuration.
	if err := m.ParseFlags(ctx, os.Args[1:]); err == flag.ErrHelp {
		os.Exit(1)
	} else if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Execute program.
	if err := m.Run(ctx); err != nil {
		m.Close()
		fmt.Fprintln(os.Stderr, err)
		api.ReportError(ctx, err)
		os.Exit(1)
	}

	// Wait for CTRL-C.
	<-ctx.Done()

	// Clean up program.
	if err := m.Close(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Main represents the program.
type Main struct {
	// Configuration path and parsed config data.
	Config Config

	// SQLite database used by SQLite service implementations.
	DB *pg.DB

	// HTTP server for handling HTTP communication.
	// SQLite services are attached to it before running.
	HTTPServer *http.Server
}

// NewMain returns a new instance of Main.
func NewMain() *Main {
	var config Config
	return &Main{
		Config: config,

		DB:         pg.NewDB(""),
		HTTPServer: http.NewServer(),
	}
}

// Close gracefully stops the program.
func (m *Main) Close() error {
	if m.HTTPServer != nil {
		if err := m.HTTPServer.Close(); err != nil {
			return err
		}
	}
	if m.DB != nil {
		if err := m.DB.Close(); err != nil {
			return err
		}
	}
	return nil
}

// ParseFlags parses the command line arguments & loads the config.
//
// This exists separately from the Run() function so that we can skip it
// during end-to-end tests. Those tests will configure manually and call Run().
func (m *Main) ParseFlags(ctx context.Context, args []string) error {
	flag.StringVar(&m.Config.DB.DSN, "dsn", "main_user:mysecretuserpassword@/astyqbaga?parseTime=true", "MySQL/MariaDB database DSN")
	flag.StringVar(&m.Config.FS.HashKey, "fs-hash-key", "00000000000000000000000000000000000000000000000000", "Hash key for naming files")
	flag.StringVar(&m.Config.HTTP.Addr, "addr", ":4000", "HTTP network address")
	flag.StringVar(&m.Config.HTTP.Domain, "domain", "", "HTTP network address")
	flag.StringVar(&m.Config.SignKey, "sign-key", "00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000", "Sign key for JWT")
	flag.StringVar(&m.Config.VerifyKey, "verify-key", "0000000000000000000000000000000000000000000000000000000000000000", "Verification key for JWT")
	flag.Parse()

	return nil
}

// Run executes the program. The configuration should already be set up before
// calling this function.
func (m *Main) Run(ctx context.Context) (err error) {
	api.ReportError = defaultReportError
	api.ReportPanic = defaultReportPanic

	// Open the database. This will instantiate the MariaDB connection
	// and execute any pending migration files.
	m.DB.DSN = m.Config.DB.DSN
	if err := m.DB.Open(); err != nil {
		return fmt.Errorf("cannot open db: %w", err)
	}

	_, err = os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %w", err)
	}

	// Instantiate PG-backed services.
	authService := jwt.NewAuthService([]byte(m.Config.SignKey), []byte(m.Config.VerifyKey))
	userService := pg.NewUserService(m.DB)
	groupService := pg.NewGroupService(m.DB)
	homeworkService := pg.NewHomeworkService(m.DB)
	hwSubmissionService := pg.NewHWSubmissionService(m.DB)
	quizService := pg.NewQuizService(m.DB)
	quizSubmissionService := pg.NewQuizSubmissionService(m.DB)
	questionService := pg.NewQuestionService(m.DB)
	responseService := pg.NewResponseService(m.DB)

	// Copy configuration settings to the HTTP server.
	m.HTTPServer.Addr = m.Config.HTTP.Addr
	m.HTTPServer.Domain = m.Config.HTTP.Domain

	// Attach underlying services to the HTTP server.
	m.HTTPServer.AuthService = authService
	m.HTTPServer.UserService = userService
	m.HTTPServer.GroupService = groupService
	m.HTTPServer.HomeworkService = homeworkService
	m.HTTPServer.HWSubmissionService = hwSubmissionService
	m.HTTPServer.QuizService = quizService
	m.HTTPServer.QuizSubmissionService = quizSubmissionService
	m.HTTPServer.QuestionService = questionService
	m.HTTPServer.ResponseService = responseService

	// Start the HTTP server.
	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	// Start the HTTP server.
	if err := m.HTTPServer.Open(); err != nil {
		return err
	}

	log.Printf("running: url=%q dsn=%q", m.HTTPServer.URL(), m.Config.DB.DSN)

	return nil
}

// Config represents the CLI configuration file.
type Config struct {
	DB struct {
		DSN string
	}

	FS struct {
		HashKey string
	}

	SignKey   string
	VerifyKey string

	HTTP struct {
		Addr   string
		Domain string
	}
}

// defaultReportError reports internal errors to stdout.
func defaultReportError(ctx context.Context, err error, args ...interface{}) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	log.Printf("error: %v", trace)
}

// defaultReportPanic reports panics to stdout.
func defaultReportPanic(err interface{}) {
	trace := fmt.Sprintf("%s\n%s", err, debug.Stack())
	log.Printf("panic: %v", trace)
}
