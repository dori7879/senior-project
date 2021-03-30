package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerQuestionRoutes is a helper function for registering all question routes.
func (s *Server) registerQuestionRoutes(r *mux.Router) {
	// API endpoint for creating questions.
	r.HandleFunc("/questions", s.handleQuestionCreate).Methods("POST")

	r.HandleFunc("/questions/{id}", s.handleQuestionUpdate).Methods("PATCH")

	// Removing a question.
	r.HandleFunc("/questions/{id}", s.handleQuestionDelete).Methods("DELETE")
}

// handleQuestionCreate handles the "POST /questions" route.
func (s *Server) handleQuestionCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Unmarshal data
	var q api.Question
	if err := json.NewDecoder(r.Body).Decode(&q); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Create the question in the database.
	err := s.QuestionService.CreateQuestion(r.Context(), &q)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&q); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuestionUpdate handles the "PATCH /questions/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleQuestionUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse question ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	var upd api.QuestionUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the question in the database.
	_, err = s.QuestionService.UpdateQuestion(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleQuestionDelete handles the "DELETE /questions/:id" route. This route
// permanently deletes the question.
func (s *Server) handleQuestionDelete(w http.ResponseWriter, r *http.Request) {
	// Parse question ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the question from the database.
	if err := s.QuestionService.DeleteQuestion(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
