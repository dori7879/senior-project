package http

import (
	"encoding/json"
	"net/http"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// registerUserRoutes is a helper function for registering all user and auth routes.
func (s *Server) registerUserRoutes(r *mux.Router) {
	r.HandleFunc("/login", s.handleLogin).Methods("POST")
	// r.HandleFunc("/signup", s.handleSignup).Methods("POST")

	// r.HandleFunc("/password/change", s.handlePasswordChange).Methods("PATCH")

	// r.HandleFunc("/profile", s.handleProfileView).Methods("GET")
	// r.HandleFunc("/profile", s.handleProfileUpdate).Methods("PUT")
	// r.HandleFunc("/profile", s.handleProfileDelete).Methods("DELETE")
}

type authResponse struct {
	Token string `json:"Token"`
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	var auth *api.Auth
	if err := json.NewDecoder(r.Body).Decode(auth); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Fetch user from the database.
	user, err := s.UserService.FindUserByEmail(r.Context(), auth.Email)
	if err != nil {
		Error(w, r, err)
		return
	}

	token, err := s.AuthService.Login(auth, user)
	if err == bcrypt.ErrMismatchedHashAndPassword {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "Incorrect password or email"))
		return
	} else if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(authResponse{Token: token}); err != nil {
		LogError(r, err)
		return
	}
}

// // Parse dial ID from path.
// id, err := strconv.Atoi(mux.Vars(r)["id"])
// if err != nil {
// 	Error(w, r, wtf.Errorf(wtf.EINVALID, "Invalid ID format"))
// 	return
// }

// // Write response to indicate success.
// w.Header().Set("Content-type", "application/json")
// w.Write([]byte(`{}`))
