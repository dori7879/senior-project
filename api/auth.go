package api

// Auth represents a auth in the system.
type Auth struct {
	Email    string `json:"Email"`
	Password string `json:"Password"`
}

// Validate returns an error if the auth contains invalid fields.
// This only performs basic validation.
func (u *Auth) Validate() error {
	if u.Email == "" {
		return Errorf(EINVALID, "Email required.")
	} else if u.Password == "" {
		return Errorf(EINVALID, "Password required.")
	}
	return nil
}

// AuthService represents a service for managing auths.
type AuthService interface {
	Login(auth *Auth, user *User) (string, error)

	Validate(tokenStr string) (int, error)
}
