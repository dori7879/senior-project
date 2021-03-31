package api

import (
	"context"
	"time"
)

// QuizSubmission represents a quiz submission in the system.
type QuizSubmission struct {
	ID int `json:"ID"`

	Grade       float32   `json:"Grade,omitempty"`
	Comments    string    `json:"Comments,omitempty"`
	SubmittedAt time.Time `json:"SubmittedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`

	StudentFullName string `json:"StudentFullName,omitempty" db:"student_fullname"`
	StudentID       int    `json:"StudentID,omitempty"`
	Student         *User  `json:"Student,omitempty"`

	QuizID int   `json:"QuizID,omitempty"`
	Quiz   *Quiz `json:"Quiz,omitempty"`

	Responses []*Response `json:"Responses"`
}

// Validate returns an error if the quiz submission contains invalid fields.
// This only performs basic validation.
func (u *QuizSubmission) Validate() error {
	if u.StudentID != 0 || u.StudentFullName != "" {
		return Errorf(EINVALID, "Enter student fullname or sign in.")
	}
	return nil
}

// QuizSubmissionService represents a service for managing quiz submissions.
type QuizSubmissionService interface {
	// Retrieves a quiz submission by ID.
	// Returns ENOTFOUND if quiz submission does not exist.
	FindQuizSubmissionByID(ctx context.Context, id int) (*QuizSubmission, error)

	// Retrieves a list of quiz submissions by filter. Also returns total count of matching
	// quiz submissions which may differ from returned results if filter.Limit is specified.
	FindQuizSubmissions(ctx context.Context, filter QuizSubmissionFilter) ([]*QuizSubmission, int, error)

	// Creates a new quiz submission.
	CreateQuizSubmission(ctx context.Context, submission *QuizSubmission) error

	// Updates a quiz submission object. Returns EUNAUTHORIZED if current quiz submission is not
	// the quiz submission that is being updated. Returns ENOTFOUND if quiz submission does not exist.
	UpdateQuizSubmission(ctx context.Context, id int, upd QuizSubmissionUpdate) (*QuizSubmission, error)

	// Permanently deletes a quiz submission and all owned dials. Returns EUNAUTHORIZED
	// if current quiz submission is not the quiz submission being deleted. Returns ENOTFOUND if
	// quiz submission does not exist.
	DeleteQuizSubmission(ctx context.Context, id int) error
}

// QuizSubmissionFilter represents a filter passed to FindQuizSubmissions().
type QuizSubmissionFilter struct {
	// Filtering fields.
	ID                *int       `json:"ID"`
	BeforeSubmittedAt *time.Time `json:"SubmittedAt"`
	AfterUpdatedAt    *time.Time `json:"UpdatedAt"`

	StudentFullName *string `json:"StudentFullName"`
	StudentID       *int    `json:"StudentID"`
	QuizID          *int    `json:"QuizID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// QuizSubmissionUpdate represents a set of fields to be updated via UpdateQuizSubmission().
type QuizSubmissionUpdate struct {
	Grade    *float32 `json:"Grade"`
	Comments *string  `json:"Comments"`

	StudentFullName *string `json:"StudentFullName"`
	StudentID       *int    `json:"StudentID"`
}
