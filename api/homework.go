package api

import (
	"context"
	"time"
)

const (
	// All
	All = "all"
	// Registered
	Registered = "registered"
)

// Homework represents a homework in the system.
type Homework struct {
	ID int `json:"ID"`

	Title       string    `json:"Title"`
	Content     string    `json:"Content"`
	MaxGrade    float32   `json:"MaxGrade"`
	StudentLink string    `json:"StudentLink"`
	TeacherLink string    `json:"TeacherLink"`
	CourseTitle string    `json:"CourseTitle"`
	Mode        string    `json:"Mode"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	OpenedAt    time.Time `json:"OpenedAt"`
	ClosedAt    time.Time `json:"ClosedAt"`

	TeacherFullName string `json:"TeacherFullName"`
	TeacherID       int    `json:"TeacherID"`
	Teacher         *User  `json:"Teacher"`

	GroupID int    `json:"GroupID"`
	Group   *Group `json:"Group"`

	Submissions []*HWSubmission `json:"Submissions,omitempty"`
}

// Validate returns an error if the homework contains invalid fields.
// This only performs basic validation.
func (u *Homework) Validate() error {
	if u.Title == "" {
		return Errorf(EINVALID, "Title required.")
	} else if u.Mode != All && u.Mode != Registered {
		return Errorf(EINVALID, "Mode is incorrect.")
	}
	return nil
}

// HomeworkService represents a service for managing homeworks.
type HomeworkService interface {
	// Retrieves a homework by ID.
	// Returns ENOTFOUND if homework does not exist.
	FindHomeworkByID(ctx context.Context, id int) (*Homework, error)

	FindHomeworkByStudentLink(ctx context.Context, link string) (*Homework, error)

	FindHomeworkByTeacherLink(ctx context.Context, link string) (*Homework, error)

	// Retrieves a list of homeworks by filter. Also returns total count of matching
	// homeworks which may differ from returned results if filter.Limit is specified.
	FindHomeworks(ctx context.Context, filter HomeworkFilter) ([]*Homework, int, error)

	// Creates a new homework.
	CreateHomework(ctx context.Context, homework *Homework) error

	// Updates a homework object. Returns EUNAUTHORIZED if current homework is not
	// the homework that is being updated. Returns ENOTFOUND if homework does not exist.
	UpdateHomework(ctx context.Context, id int, upd HomeworkUpdate) (*Homework, error)

	// Permanently deletes a homework and all owned dials. Returns EUNAUTHORIZED
	// if current homework is not the homework being deleted. Returns ENOTFOUND if
	// homework does not exist.
	DeleteHomework(ctx context.Context, id int) error
}

// HomeworkFilter represents a filter passed to FindHomeworks().
type HomeworkFilter struct {
	// Filtering fields.
	ID          *int    `json:"ID"`
	Title       *string `json:"Title"`
	StudentLink *string `json:"StudentLink"`
	TeacherLink *string `json:"TeacherLink"`
	CourseTitle *string `json:"CourseTitle"`
	Mode        *string `json:"Mode"`

	TeacherFullName *string `json:"TeacherFullName"`
	TeacherID       *int    `json:"TeacherID"`
	GroupID         *int    `json:"GroupID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// HomeworkUpdate represents a set of fields to be updated via UpdateHomework().
type HomeworkUpdate struct {
	Title       *string    `json:"Title"`
	Content     *string    `json:"Content"`
	MaxGrade    *float32   `json:"MaxGrade"`
	CourseTitle *string    `json:"CourseTitle"`
	Mode        *string    `json:"Mode"`
	OpenedAt    *time.Time `json:"OpenedAt"`
	ClosedAt    *time.Time `json:"ClosedAt"`

	TeacherFullName *string `json:"TeacherFullName"`
	TeacherID       *int    `json:"TeacherID"`
	GroupID         *int    `json:"GroupID"`
}
