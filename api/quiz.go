package api

import (
	"context"
	"time"
)

// Quiz represents a quiz in the system.
type Quiz struct {
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

	TeacherFullName string `json:"TeacherFullName" db:"teacher_fullname"`
	TeacherID       int    `json:"TeacherID"`
	Teacher         *User  `json:"Teacher"`

	GroupID int    `json:"GroupID"`
	Group   *Group `json:"Group"`

	Submissions []*QuizSubmission `json:"Submissions"`
	Questions   []*Question       `json:"Questions"`
}

// Validate returns an error if the quiz contains invalid fields.
// This only performs basic validation.
func (q *Quiz) Validate() error {
	if q.Title == "" {
		return Errorf(EINVALID, "Title required.")
	} else if q.Mode != All && q.Mode != Registered {
		return Errorf(EINVALID, "Mode is incorrect.")
	}
	return nil
}

// QuizService represents a service for managing quizzes.
type QuizService interface {
	// Retrieves a quiz by ID.
	// Returns ENOTFOUND if quiz does not exist.
	FindQuizByID(ctx context.Context, id int) (*Quiz, error)

	FindQuizByStudentLink(ctx context.Context, link string) (*Quiz, error)

	FindQuizByTeacherLink(ctx context.Context, link string) (*Quiz, error)

	// Retrieves a list of quizzes by filter. Also returns total count of matching
	// quizzes which may differ from returned results if filter.Limit is specified.
	FindQuizzes(ctx context.Context, filter QuizFilter) ([]*Quiz, int, error)

	// Creates a new quiz.
	CreateQuiz(ctx context.Context, quiz *Quiz) error

	// Updates a quiz object. Returns EUNAUTHORIZED if current quiz is not
	// the quiz that is being updated. Returns ENOTFOUND if quiz does not exist.
	UpdateQuiz(ctx context.Context, id int, upd QuizUpdate) (*Quiz, error)

	// Permanently deletes a quiz and all owned dials. Returns EUNAUTHORIZED
	// if current quiz is not the quiz being deleted. Returns ENOTFOUND if
	// quiz does not exist.
	DeleteQuiz(ctx context.Context, id int) error
}

// QuizFilter represents a filter passed to FindQuizs().
type QuizFilter struct {
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

// QuizUpdate represents a set of fields to be updated via UpdateQuiz().
type QuizUpdate struct {
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
