package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerAttSubmissionPrivateRoutes is a helper function for registering private attendance submission routes.
func (s *Server) registerAttSubmissionPrivateRoutes(r *mux.Router) {
	// Listing of all attendance submissions a teacher is an owner of.
	r.HandleFunc("/attendances/submissions", s.handleAttSubmissionList).Methods("GET")
}

// registerAttSubmissionPublicRoutes is a helper function for registering public attendance submission routes.
func (s *Server) registerAttSubmissionPublicRoutes(r *mux.Router) {
	// View a single attendance.
	r.HandleFunc("/attendances/submissions/{id}", s.handleAttSubmissionView).Methods("GET")

	// API endpoint for creating attendance submissions.
	r.HandleFunc("/attendances/{attID}/submissions", s.handleAttSubmissionCreate).Methods("POST")

	r.HandleFunc("/attendances/submissions/{id}", s.handleAttSubmissionUpdate).Methods("PATCH")

	// Removing a attendance.
	r.HandleFunc("/attendances/submissions/{id}", s.handleAttSubmissionDelete).Methods("DELETE")
}

// handleAttSubmissionList handles the "GET /attendances/submissions" route. This route can optionally
// accept filter arguments and outputs a list of all attendances that the current
// user is related to.
func (s *Server) handleAttSubmissionList(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Fetch submissions from database.
	subs, n, err := s.AttSubmissionService.FindAttSubmissions(r.Context(), api.AttSubmissionFilter{StudentID: &user.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Render output based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Submissions []*api.AttSubmission `json:"Submissions"`
		N           int                  `json:"N"`
	}{
		Submissions: subs,
		N:           n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleAttSubmissionView handles the "GET /attendances/submissions/:id" route.
func (s *Server) handleAttSubmissionView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())

	// Parse ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Fetch submission from the database.
	sub, err := s.AttSubmissionService.FindAttSubmissionByID(r.Context(), id)
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

// handleAttSubmissionCreate handles the "POST /attendances/submissions" route.
func (s *Server) handleAttSubmissionCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Parse attendance submission ID from the path.
	attID, err := strconv.Atoi(mux.Vars(r)["attID"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Unmarshal data
	sub := api.AttSubmission{}
	if err := json.NewDecoder(r.Body).Decode(&sub); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	if user != nil {
		sub.StudentID = user.ID
	}

	sub.AttendanceID = attID

	// Create submission in the database.
	err = s.AttSubmissionService.CreateAttSubmission(r.Context(), &sub)
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

// handleAttSubmissionUpdate handles the "PATCH /attendances/submissions/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleAttSubmissionUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse attendance submission ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	upd := api.AttSubmissionUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher && upd.Present != nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "Student attempts to update present field"))
		return
	}

	// Update the attendance submission in the database.
	_, err = s.AttSubmissionService.UpdateAttSubmission(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleAttSubmissionDelete handles the "DELETE /attendances/submissions/:id" route. This route
// permanently deletes the attendance submission.
func (s *Server) handleAttSubmissionDelete(w http.ResponseWriter, r *http.Request) {
	// Parse attendance submission ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the attendance submission from the database.
	if err := s.AttSubmissionService.DeleteAttSubmission(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
