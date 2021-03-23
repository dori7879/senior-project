package api

import (
	"context"
)

// Response represents a response in the system.
type Response struct {
	ID int `json:"ID"`

	Comments  string  `json:"Comments"`
	IsCorrect bool    `json:"IsCorrect"`
	Grade     float32 `json:"Grade"`
	Type      int     `json:"Type"`

	OpenResponse           string `json:"OpenResponse"`
	TrueFalseResponse      bool   `json:"TrueFalseResponse" db:"truefalse_response"`
	MultipleChoiceResponse []int  `json:"MultipleChoiceResponse" db:"multiplechoice_response"`
	SingleChoiceResponse   int    `json:"SingleChoiceResponse" db:"singlechoice_response"`

	SubmissionID int             `json:"SubmissionID" db:"quiz_submission_id"`
	Submission   *QuizSubmission `json:"Submission"`
}

// Validate returns an error if the response contains invalid fields.
// This only performs basic validation.
func (u *Response) Validate() error {
	if u.Type != 0 {
		return Errorf(EINVALID, "Type required.")
	}
	return nil
}

// ResponseService represents a service for managing responses.
type ResponseService interface {
	// Retrieves a response by ID.
	// Returns ENOTFOUND if response does not exist.
	FindResponseByID(ctx context.Context, id int) (*Response, error)

	// Retrieves a list of responses by filter. Also returns total count of matching
	// responses which may differ from returned results if filter.Limit is specified.
	FindResponses(ctx context.Context, filter ResponseFilter) ([]*Response, int, error)

	// Creates a new response.
	CreateResponse(ctx context.Context, response *Response) error

	// Updates a response object. Returns EUNAUTHORIZED if current response is not
	// the response that is being updated. Returns ENOTFOUND if response does not exist.
	UpdateResponse(ctx context.Context, id int, upd ResponseUpdate) (*Response, error)

	// Permanently deletes a response and all owned dials. Returns EUNAUTHORIZED
	// if current response is not the response being deleted. Returns ENOTFOUND if
	// response does not exist.
	DeleteResponse(ctx context.Context, id int) error
}

// ResponseFilter represents a filter passed to FindResponses().
type ResponseFilter struct {
	// Filtering fields.
	ID        int  `json:"ID"`
	IsCorrect bool `json:"IsCorrect"`
	Type      int  `json:"Type"`

	SubmissionID int `json:"SubmissionID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// ResponseUpdate represents a set of fields to be updated via UpdateResponse().
type ResponseUpdate struct {
	Comments  string  `json:"Comments"`
	IsCorrect bool    `json:"IsCorrect"`
	Grade     float32 `json:"Grade"`
	Type      int     `json:"Type"`

	OpenResponse           string `json:"OpenResponse"`
	TrueFalseResponse      bool   `json:"TrueFalseResponse"`
	MultipleChoiceResponse []int  `json:"MultipleChoiceResponse"`
	SingleChoiceResponse   int    `json:"SingleChoiceResponse"`
}
