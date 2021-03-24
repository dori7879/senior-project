package http

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/dori7879/senior-project/api"
)

// Error prints & optionally logs an error message.
func Error(w http.ResponseWriter, r *http.Request, err error) {
	// Extract error code & message.
	code, message := api.ErrorCode(err), api.ErrorMessage(err)

	// Log & report internal errors.
	if code == api.EINTERNAL {
		api.ReportError(r.Context(), err, r)
		LogError(r, err)
	}

	// Print user message to response based on reqeust accept header.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(ErrorStatusCode(code))
	json.NewEncoder(w).Encode(&ErrorResponse{Error: message})
}

// ErrorResponse represents a JSON structure for error output.
type ErrorResponse struct {
	Error string `json:"error"`
}

// LogError logs an error with the HTTP route information.
func LogError(r *http.Request, err error) {
	log.Printf("[http] error: %s %s: %s", r.Method, r.URL.Path, err)
}

// lookup of application error codes to HTTP status codes.
var codes = map[string]int{
	api.ECONFLICT:       http.StatusConflict,
	api.EINVALID:        http.StatusBadRequest,
	api.ENOTFOUND:       http.StatusNotFound,
	api.ENOTIMPLEMENTED: http.StatusNotImplemented,
	api.EUNAUTHORIZED:   http.StatusUnauthorized,
	api.EINTERNAL:       http.StatusInternalServerError,
}

// ErrorStatusCode returns the associated HTTP status code for a WTF error code.
func ErrorStatusCode(code string) int {
	if v, ok := codes[code]; ok {
		return v
	}
	return http.StatusInternalServerError
}

// FromErrorStatusCode returns the associated WTF code for an HTTP status code.
func FromErrorStatusCode(code int) string {
	for k, v := range codes {
		if v == code {
			return k
		}
	}
	return api.EINTERNAL
}
