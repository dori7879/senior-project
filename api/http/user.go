package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
	"golang.org/x/crypto/bcrypt"
)

// registerUserRoutes is a helper function for registering user and auth routes.
func (s *Server) registerUserRoutes(r *mux.Router) {
	r.HandleFunc("/users/password", s.handlePasswordChange).Methods("PATCH")
	r.HandleFunc("/users/suggestions", s.handleUserSuggestions).Methods("GET")

	r.HandleFunc("/profile", s.handleProfileView).Methods("GET")
	r.HandleFunc("/profile", s.handleProfileUpdate).Methods("PUT")
	r.HandleFunc("/profile", s.handleProfileDelete).Methods("DELETE")
}

// registerAujthRoutes is a helper function for registering auth routes for unauthenticated users.
func (s *Server) registerAuthRoutes(r *mux.Router) {
	r.HandleFunc("/login", s.handleLogin).Methods("POST")
	r.HandleFunc("/signup", s.handleSignup).Methods("POST")
}

func (s *Server) handleLogin(w http.ResponseWriter, r *http.Request) {
	auth := &api.Auth{}
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
	if err := json.NewEncoder(w).Encode(struct {
		Token string `json:"Token"`
	}{
		Token: token,
	}); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleSignup(w http.ResponseWriter, r *http.Request) {
	// Parse password first
	in := &struct {
		Password string `json:"Password"`
		api.User
	}{}
	if err := json.NewDecoder(r.Body).Decode(in); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON field"))
		return
	}

	// Parse rest of the fields
	var user api.User
	user.FirstName = in.FirstName
	user.LastName = in.LastName
	user.Email = in.Email
	user.IsTeacher = in.IsTeacher

	// Generate password hash and assign it to user
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 12)
	if err != nil {
		Error(w, r, err)
		return
	}
	user.PasswordHash = passwordHash

	// Create user
	err = s.UserService.CreateUser(r.Context(), &user)
	if err != nil {
		Error(w, r, api.Errorf(api.EINTERNAL, "Could not create user"))
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(`{}`))
}

func (s *Server) handlePasswordChange(w http.ResponseWriter, r *http.Request) {
	passwords := &struct {
		Old string `json:"OldPassword"`
		New string `json:"NewPassword"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(passwords); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	user := api.UserFromContext(r.Context())

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(passwords.Old)); err == bcrypt.ErrMismatchedHashAndPassword {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "Incorrect old password"))
		return
	} else if err != nil {
		Error(w, r, err)
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(passwords.New), 12)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Construct update instance
	upd := api.UserUpdate{
		PasswordHash: &passwordHash,
	}

	_, err = s.UserService.UpdateUser(r.Context(), user.ID, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func (s *Server) handleUserSuggestions(w http.ResponseWriter, r *http.Request) {
	var uFilter api.UserFilter

	uFilter.Offset, _ = strconv.Atoi(r.URL.Query().Get("offset"))
	uFilter.Limit, _ = strconv.Atoi(r.URL.Query().Get("limit"))
	if uFilter.Limit == 0 {
		uFilter.Limit = 10
	}
	emailSubStr := r.URL.Query().Get("email")
	uFilter.EmailSubStr = &emailSubStr
	isTeacher, _ := strconv.ParseBool(r.URL.Query().Get("isTeacher"))
	uFilter.IsTeacher = &isTeacher

	users, n, err := s.UserService.FindUsers(r.Context(), uFilter)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(struct {
		Users []*api.User `json:"Users"`
		N     int         `json:"n"`
	}{
		Users: users,
		N:     n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleProfileView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	var err error
	uFilter := api.UserFilter{}
	uFilter.ID = &user.ID
	uFilter.IsTeacher = &user.IsTeacher

	user.SharedGroups.Groups, _, err = s.GroupService.FindGroupsByMember(r.Context(), uFilter)
	if err != nil {
		Error(w, r, err)
		return
	}

	gFilter := api.GroupFilter{
		OwnerID: &user.ID,
	}

	user.OwnedGroups.Groups, _, err = s.GroupService.FindGroups(r.Context(), gFilter)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		LogError(r, err)
		return
	}
}

func (s *Server) handleProfileUpdate(w http.ResponseWriter, r *http.Request) {
	upd := api.UserUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	if _, err := s.UserService.UpdateUser(r.Context(), api.UserIDFromContext(r.Context()), upd); err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}

func (s *Server) handleProfileDelete(w http.ResponseWriter, r *http.Request) {
	// Parse password to check for it's correctness
	rawPassword := &struct {
		Content string `json:"Password"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(rawPassword); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid password JSON field"))
		return
	}

	user := api.UserFromContext(r.Context())

	if err := bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(rawPassword.Content)); err == bcrypt.ErrMismatchedHashAndPassword {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "Incorrect password"))
		return
	} else if err != nil {
		Error(w, r, err)
		return
	}

	if err := s.UserService.DeleteUser(r.Context(), user.ID); err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{}`))
}
