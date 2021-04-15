package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerHWSubmissionPrivateRoutes is a helper function for registering private homework submission routes.
func (s *Server) registerHWSubmissionPrivateRoutes(r *mux.Router) {
	// Listing of all homework submissions a teacher is an owner of.
	r.HandleFunc("/homeworks/submissions", s.handleHWSubmissionList).Methods("GET")
}

// registerHWSubmissionPublicRoutes is a helper function for registering public homework submission routes.
func (s *Server) registerHWSubmissionPublicRoutes(r *mux.Router) {
	// View a single homework.
	r.HandleFunc("/homeworks/submissions/{id}", s.handleHWSubmissionView).Methods("GET")

	// API endpoint for creating homework submissions.
	r.HandleFunc("/homeworks/{hwID}/submissions", s.handleHWSubmissionCreate).Methods("POST")

	r.HandleFunc("/homeworks/submissions/{id}", s.handleHWSubmissionUpdate).Methods("PATCH")

	// Removing a homework.
	r.HandleFunc("/homeworks/submissions/{id}", s.handleHWSubmissionDelete).Methods("DELETE")
}

// handleHWSubmissionList handles the "GET /homeworks/submissions" route. This route can optionally
// accept filter arguments and outputs a list of all homeworks that the current
// user is related to.
func (s *Server) handleHWSubmissionList(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Fetch submissions from database.
	subs, n, err := s.HWSubmissionService.FindHWSubmissions(r.Context(), api.HWSubmissionFilter{StudentID: &user.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Render output based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Submissions []*api.HWSubmission `json:"Submissions"`
		N           int                 `json:"N"`
	}{
		Submissions: subs,
		N:           n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleHWSubmissionView handles the "GET /homeworks/submissions/:id" route.
func (s *Server) handleHWSubmissionView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())

	// Parse ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Fetch submission from the database.
	sub, err := s.HWSubmissionService.FindHWSubmissionByID(r.Context(), id)
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

// handleHWSubmissionCreate handles the "POST /homeworks/submissions" route.
func (s *Server) handleHWSubmissionCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Parse homework submission ID from the path.
	hwID, err := strconv.Atoi(mux.Vars(r)["hwID"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Unmarshal data
	sub := api.HWSubmission{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	if user != nil {
		sub.StudentID = user.ID
	}

	sub.HomeworkID = hwID

	// Create submission in the database.
	err = s.HWSubmissionService.CreateHWSubmission(r.Context(), &sub)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&sub); err != nil {
		LogError(r, err)
		return
	}
}

// handleHWSubmissionUpdate handles the "PATCH /homeworks/submissions/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleHWSubmissionUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse homework submission ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	upd := api.HWSubmissionUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher && (upd.Comments != nil || upd.Grade != nil) {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "Student attempts to update comments or grade"))
		return
	}

	// Update the homework submission in the database.
	_, err = s.HWSubmissionService.UpdateHWSubmission(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleHWSubmissionDelete handles the "DELETE /homeworks/submissions/:id" route. This route
// permanently deletes the homework submission.
func (s *Server) handleHWSubmissionDelete(w http.ResponseWriter, r *http.Request) {
	// Parse homework submission ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the homework submission from the database.
	if err := s.HWSubmissionService.DeleteHWSubmission(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
