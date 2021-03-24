package api

import (
	"context"
	"time"
)

// Question represents a question in the system.
type Question struct {
	ID int `json:"ID"`

	Content string `json:"Content"`
	Type    int    `json:"Type"`
	Fixed   bool   `json:"Fixed"`

	Choices []*string `json:"Choices"`

	OpenAnswer           string `json:"OpenAnswer"`
	TrueFalseAnswer      bool   `json:"TrueFalseAnswer" db:"truefalse_answer"`
	MultipleChoiceAnswer []int  `json:"MultipleChoiceAnswer" db:"multiplechoice_answer"`
	SingleChoiceAnswer   int    `json:"SingleChoiceAnswer" db:"singlechoice_answer"`

	CreatedAt time.Time `json:"CreatedAt"`
	UpdatedAt time.Time `json:"UpdatedAt"`

	QuizID int `json:"QuizID"`
}

// Validate returns an error if the question contains invalid fields.
// This only performs basic validation.
func (u *Question) Validate() error {
	if u.Content == "" {
		return Errorf(EINVALID, "Content required.")
	}
	return nil
}

// QuestionService represents a service for managing questions.
type QuestionService interface {
	// Retrieves a question by ID.
	// Returns ENOTFOUND if question does not exist.
	FindQuestionByID(ctx context.Context, id int) (*Question, error)

	// Retrieves a list of questions by filter. Also returns total count of matching
	// questions which may differ from returned results if filter.Limit is specified.
	FindQuestions(ctx context.Context, filter QuestionFilter) ([]*Question, int, error)

	// Creates a new question.
	CreateQuestion(ctx context.Context, question *Question) error

	// Updates a question object. Returns EUNAUTHORIZED if current question is not
	// the question that is being updated. Returns ENOTFOUND if question does not exist.
	UpdateQuestion(ctx context.Context, id int, upd QuestionUpdate) (*Question, error)

	// Permanently deletes a question and all owned dials. Returns EUNAUTHORIZED
	// if current question is not the question being deleted. Returns ENOTFOUND if
	// question does not exist.
	DeleteQuestion(ctx context.Context, id int) error
}

// QuestionFilter represents a filter passed to FindQuestions().
type QuestionFilter struct {
	// Filtering fields.
	ID    *int  `json:"ID"`
	Type  *int  `json:"Type"`
	Fixed *bool `json:"Fixed"`

	QuizID *int `json:"QuizID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// QuestionUpdate represents a set of fields to be updated via UpdateQuestion().
type QuestionUpdate struct {
	Content *string `json:"Content"`
	Type    *int    `json:"Type"`
	Fixed   *bool   `json:"Fixed"`

	Choices *[]*string `json:"Choices"`

	OpenAnswer           *string `json:"OpenAnswer"`
	TrueFalseAnswer      *bool   `json:"TrueFalseAnswer"`
	MultipleChoiceAnswer *[]int  `json:"MultipleChoiceAnswer"`
	SingleChoiceAnswer   *int    `json:"SingleChoiceAnswer"`
}
