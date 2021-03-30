package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerResponseRoutes is a helper function for registering all response routes.
func (s *Server) registerResponseRoutes(r *mux.Router) {
	// API endpoint for creating responses.
	r.HandleFunc("/responses", s.handleResponseCreate).Methods("POST")

	r.HandleFunc("/responses/{id}", s.handleResponseUpdate).Methods("PATCH")

	// Removing a response.
	r.HandleFunc("/responses/{id}", s.handleResponseDelete).Methods("DELETE")
}

// handleResponseCreate handles the "POST /responses" route.
func (s *Server) handleResponseCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Unmarshal data
	var rn api.Response
	if err := json.NewDecoder(r.Body).Decode(&rn); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Create the response in the database.
	err := s.ResponseService.CreateResponse(r.Context(), &rn)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(&rn); err != nil {
		LogError(r, err)
		return
	}
}

// handleResponseUpdate handles the "PATCH /responses/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleResponseUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse response ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	var upd api.ResponseUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the response in the database.
	_, err = s.ResponseService.UpdateResponse(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleResponseDelete handles the "DELETE /responses/:id" route. This route
// permanently deletes the response.
func (s *Server) handleResponseDelete(w http.ResponseWriter, r *http.Request) {
	// Parse response ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the response from the database.
	if err := s.ResponseService.DeleteResponse(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
