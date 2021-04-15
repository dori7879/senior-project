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
}

// registerQuizSubmissionPublicRoutes is a helper function for registering public quiz submission routes.
func (s *Server) registerQuizSubmissionPublicRoutes(r *mux.Router) {
	// View a single quiz submission.
	r.HandleFunc("/quizzes/submissions/{id}", s.handleQuizSubmissionView).Methods("GET")

	// API endpoint for creating quiz submissions.
	r.HandleFunc("/quizzes/{quizID}/submissions", s.handleQuizSubmissionCreate).Methods("POST")

	r.HandleFunc("/quizzes/submissions/{id}", s.handleQuizSubmissionUpdate).Methods("PATCH")

	// Removing a quiz.
	r.HandleFunc("/quizzes/submissions/{id}", s.handleQuizSubmissionDelete).Methods("DELETE")
}

// handleQuizSubmissionList handles the "GET /quizzes/submissions" route. This route can optionally
// accept filter arguments and outputs a list of all quizzes that the current
// user is related to.
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

	if user != nil && !user.IsTeacher && sub.StudentID != user.ID {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{}`))
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

	// Parse homework submission ID from the path.
	quizID, err := strconv.Atoi(mux.Vars(r)["quizID"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Unmarshal submission data
	sub := api.QuizSubmission{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	if user != nil {
		sub.StudentID = user.ID
	}

	sub.QuizID = quizID

	// Create submission in the database.
	err = s.QuizSubmissionService.CreateQuizSubmission(r.Context(), &sub)
	if err != nil {
		Error(w, r, err)
		return
	}

	for _, rn := range sub.Responses {
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
	upd := api.QuizSubmissionUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher && (upd.Comments != nil || upd.Grade != nil) {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "Student attempts to update comments or grade"))
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
