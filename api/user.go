package api

import (
	"context"
	"time"
)

// User represents a user in the system.
type User struct {
	ID int `json:"ID"`

	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
	Email     string `json:"Email"`
	IsTeacher bool   `json:"IsTeacher"`

	PasswordHash []byte `json:"-"`

	// Timestamps for user creation & last update.
	DateJoined time.Time `json:"DateJoined"`

	Groups []*int `json:"Groups"`
}

// Validate returns an error if the user contains invalid fields.
// This only performs basic validation.
func (u *User) Validate() error {
	if u.Email == "" {
		return Errorf(EINVALID, "Email required.")
	}
	return nil
}

// UserService represents a service for managing users.
type UserService interface {
	// Retrieves a user by ID.
	// Returns ENOTFOUND if user does not exist.
	FindUserByID(ctx context.Context, id int) (*User, error)

	// Retrieves a user by email.
	// Returns ENOTFOUND if user does not exist.
	FindUserByEmail(ctx context.Context, email string) (*User, error)

	// Retrieves a list of users by filter. Also returns total count of matching
	// users which may differ from returned results if filter.Limit is specified.
	FindUsers(ctx context.Context, filter UserFilter) ([]*User, int, error)

	// Creates a new user.
	CreateUser(ctx context.Context, user *User) error

	// Updates a user object. Returns EUNAUTHORIZED if current user is not
	// the user that is being updated. Returns ENOTFOUND if user does not exist.
	UpdateUser(ctx context.Context, id int, upd UserUpdate) (*User, error)

	// Permanently deletes a user and all owned dials. Returns EUNAUTHORIZED
	// if current user is not the user being deleted. Returns ENOTFOUND if
	// user does not exist.
	DeleteUser(ctx context.Context, id int) error
}

// UserFilter represents a filter passed to FindUsers().
type UserFilter struct {
	// Filtering fields.
	ID        *int    `json:"ID"`
	Email     *string `json:"Email"`
	IsTeacher *string `json:"IsTeacher"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// UserUpdate represents a set of fields to be updated via UpdateUser().
type UserUpdate struct {
	FirstName *string `json:"FirstName"`
	LastName  *string `json:"LastName"`
	Email     *string `json:"Email"`
	IsTeacher *bool   `json:"IsTeacher"`
}
