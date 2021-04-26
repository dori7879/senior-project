package api

import (
	"context"
	"time"
)

// AttSubmission represents a attendance submission in the system.
type AttSubmission struct {
	ID int `json:"ID"`

	Present     bool      `json:"Present"`
	PIN         string    `json:"PIN"`
	SubmittedAt time.Time `json:"SubmittedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`

	StudentFullName string `json:"StudentFullName,omitempty"`
	StudentID       int    `json:"StudentID,omitempty"`
	Student         *User  `json:"Student,omitempty"`

	AttendanceID int         `json:"AttendanceID"`
	Attendance   *Attendance `json:"Attendance,omitempty"`
}

// Validate returns an error if the attendance submission contains invalid fields.
// This only performs basic validation.
func (u *AttSubmission) Validate() error {
	return nil
}

// AttSubmissionService represents a service for managing attendance submissions.
type AttSubmissionService interface {
	// Retrieves a attendance submission by ID.
	// Returns ENOTFOUND if attendance submission does not exist.
	FindAttSubmissionByID(ctx context.Context, id int) (*AttSubmission, error)

	// Retrieves a list of attendance submissions by filter. Also returns total count of matching
	// attendance submissions which may differ from returned results if filter.Limit is specified.
	FindAttSubmissions(ctx context.Context, filter AttSubmissionFilter) ([]*AttSubmission, int, error)

	// Creates a new attendance submission.
	CreateAttSubmission(ctx context.Context, submission *AttSubmission) error

	// Updates a attendance submission object. Returns EUNAUTHORIZED if current attendance submission is not
	// the attendance submission that is being updated. Returns ENOTFOUND if attendance submission does not exist.
	UpdateAttSubmission(ctx context.Context, id int, upd AttSubmissionUpdate) (*AttSubmission, error)

	// Permanently deletes a attendance submission and all owned dials. Returns EUNAUTHORIZED
	// if current attendance submission is not the attendance submission being deleted. Returns ENOTFOUND if
	// attendance submission does not exist.
	DeleteAttSubmission(ctx context.Context, id int) error
}

// AttSubmissionFilter represents a filter passed to FindAttSubmissions().
type AttSubmissionFilter struct {
	// Filtering fields.
	ID                *int       `json:"ID"`
	BeforeSubmittedAt *time.Time `json:"SubmittedAt"`
	AfterUpdatedAt    *time.Time `json:"UpdatedAt"`

	StudentFullName *string `json:"StudentFullName" db:"student_fullname"`
	StudentID       *int    `json:"StudentID"`
	AttendanceID    *int    `json:"AttendanceID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// AttSubmissionUpdate represents a set of fields to be updated via UpdateAttSubmission().
type AttSubmissionUpdate struct {
	Present *bool `json:"Present"`

	StudentFullName *string `json:"StudentFullName"`
	StudentID       *int    `json:"StudentID"`
}
