package api

import (
	"context"
	"time"
)

// HWSubmission represents a hw submission in the system.
type HWSubmission struct {
	ID int `json:"ID"`

	Response    string    `json:"Response"`
	Grade       float32   `json:"Grade"`
	Comments    string    `json:"Comments"`
	SubmittedAt time.Time `json:"SubmittedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`

	StudentFullName string `json:"StudentFullName"`
	StudentID       int    `json:"StudentID"`
	Student         *User  `json:"Student"`

	HomeworkID int       `json:"HomeworkID"`
	Homework   *Homework `json:"Homework"`
}

// Validate returns an error if the hw submission contains invalid fields.
// This only performs basic validation.
func (u *HWSubmission) Validate() error {
	if u.Response == "" {
		return Errorf(EINVALID, "Response required.")
	}
	return nil
}

// HWSubmissionService represents a service for managing hw submissions.
type HWSubmissionService interface {
	// Retrieves a hw submission by ID.
	// Returns ENOTFOUND if hw submission does not exist.
	FindHWSubmissionByID(ctx context.Context, id int) (*HWSubmission, error)

	// Retrieves a list of hw submissions by filter. Also returns total count of matching
	// hw submissions which may differ from returned results if filter.Limit is specified.
	FindHWSubmissions(ctx context.Context, filter HWSubmissionFilter) ([]*HWSubmission, int, error)

	// Creates a new hw submission.
	CreateHWSubmission(ctx context.Context, submission *HWSubmission) error

	// Updates a hw submission object. Returns EUNAUTHORIZED if current hw submission is not
	// the hw submission that is being updated. Returns ENOTFOUND if hw submission does not exist.
	UpdateHWSubmission(ctx context.Context, id int, upd HWSubmissionUpdate) (*HWSubmission, error)

	// Permanently deletes a hw submission and all owned dials. Returns EUNAUTHORIZED
	// if current hw submission is not the hw submission being deleted. Returns ENOTFOUND if
	// hw submission does not exist.
	DeleteHWSubmission(ctx context.Context, id int) error
}

// HWSubmissionFilter represents a filter passed to FindHWSubmissions().
type HWSubmissionFilter struct {
	// Filtering fields.
	ID          *int       `json:"ID"`
	SubmittedAt *time.Time `json:"SubmittedAt"`
	UpdatedAt   *time.Time `json:"UpdatedAt"`

	StudentFullName *string `json:"StudentFullName"`
	StudentID       *int    `json:"StudentID"`
	HomeworkID      *int    `json:"HomeworkID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// HWSubmissionUpdate represents a set of fields to be updated via UpdateHWSubmission().
type HWSubmissionUpdate struct {
	Response    *string    `json:"Response"`
	Grade       *float32   `json:"Grade"`
	Comments    *string    `json:"Comments"`
	SubmittedAt *time.Time `json:"SubmittedAt"`
	UpdatedAt   *time.Time `json:"UpdatedAt"`

	StudentFullName *string `json:"StudentFullName"`
	StudentID       *int    `json:"StudentID"`
}
