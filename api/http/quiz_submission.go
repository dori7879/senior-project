package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerQuizSubmissionPrivateRoutes is a helper function for registering private quiz submission routes.
func (s *Server) registerQuizSubmissionPrivateRoutes(r *mux.Router) {
	// Listing of all quiz submissions a student is an owner of.
	r.HandleFunc("/quizzes/submissions", s.handleQuizSubmissionList).Methods("GET")

	// View a single quiz.
	r.HandleFunc("/quizzes/submissions/{id}", s.handleQuizSubmissionView).Methods("GET")
}

// registerQuizSubmissionPublicRoutes is a helper function for registering public quiz submission routes.
func (s *Server) registerQuizSubmissionPublicRoutes(r *mux.Router) {
	// API endpoint for creating quiz submissions.
	r.HandleFunc("/quizzes/submissions", s.handleQuizSubmissionCreate).Methods("POST")

	r.HandleFunc("/quizzes/submissions/{id}", s.handleQuizSubmissionUpdate).Methods("PATCH")

	// Removing a quiz.
	r.HandleFunc("/quizzes/submissions/{id}", s.handleQuizSubmissionDelete).Methods("DELETE")
}

// handleQuizSubmissionList handles the "GET /quizzes/submissions" route. This route can optionally
// accept filter arguments and outputs a list of all quizzes that the current
// user is a member of.
func (s *Server) handleQuizSubmissionList(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Fetch submissions from database.
	subs, n, err := s.QuizSubmissionService.FindQuizSubmissions(r.Context(), api.QuizSubmissionFilter{StudentID: &user.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Render output based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Submissions []*api.QuizSubmission `json:"Submissions"`
		N           int                   `json:"N"`
	}{
		Submissions: subs,
		N:           n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizSubmissionView handles the "GET /quizzes/submissions/:id" route.
func (s *Server) handleQuizSubmissionView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Parse ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Fetch submission from the database.
	sub, err := s.QuizSubmissionService.FindQuizSubmissionByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(sub); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizSubmissionCreate handles the "POST /quizzes/submissions" route.
func (s *Server) handleQuizSubmissionCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Unmarshal submission data first
	var sub api.QuizSubmission
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Then unmarshal responses JSON data
	var responses []*api.Response
	if err := json.NewDecoder(r.Body).Decode(&responses); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Create submission in the database.
	err := s.QuizSubmissionService.CreateQuizSubmission(r.Context(), &sub)
	if err != nil {
		Error(w, r, err)
		return
	}

	for _, rn := range responses {
		rn.SubmissionID = sub.ID

		err = s.ResponseService.CreateResponse(r.Context(), rn)
		if err != nil {
			Error(w, r, err)
			return
		}
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&sub); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizSubmissionUpdate handles the "PATCH /quizzes/submissions/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleQuizSubmissionUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse quiz submission ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	var upd api.QuizSubmissionUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the quiz submission in the database.
	_, err = s.QuizSubmissionService.UpdateQuizSubmission(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleQuizSubmissionDelete handles the "DELETE /quizzes/submissions/:id" route. This route
// permanently deletes the quiz submission.
func (s *Server) handleQuizSubmissionDelete(w http.ResponseWriter, r *http.Request) {
	// Parse quiz submission ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the quiz submission from the database.
	if err := s.QuizSubmissionService.DeleteQuizSubmission(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
