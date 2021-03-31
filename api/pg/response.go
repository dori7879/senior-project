package pg

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/dori7879/senior-project/api"
	"github.com/jackc/pgtype"
)

// Ensure service implements interface.
var _ api.ResponseService = (*ResponseService)(nil)

// ResponseService represents a service for managing responses.
type ResponseService struct {
	db *DB
}

// NewResponseService returns a new instance of ResponseService.
func NewResponseService(db *DB) *ResponseService {
	return &ResponseService{db: db}
}

// FindResponseByID retrieves a response by ID.
// Returns ENOTFOUND if response does not exist.
func (s *ResponseService) FindResponseByID(ctx context.Context, id int) (*api.Response, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch response.
	r, err := findResponseByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return r, nil
}

// FindResponses retrieves a list of responses by filter. Also returns total count of
// matching responses which may differ from returned results if filter.Limit is specified.
func (s *ResponseService) FindResponses(ctx context.Context, filter api.ResponseFilter) ([]*api.Response, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findResponses(ctx, tx, filter)
}

// CreateResponse creates a new response.
func (s *ResponseService) CreateResponse(ctx context.Context, r *api.Response) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new response object.
	if err := createResponse(ctx, tx, r); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateResponse updates a response object. Returns EUNAUTHORIZED if current response is
// not the response that is being updated. Returns ENOTFOUND if response does not exist.
func (s *ResponseService) UpdateResponse(ctx context.Context, id int, upd api.ResponseUpdate) (*api.Response, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update response.
	r, err := updateResponse(ctx, tx, id, upd)
	if err != nil {
		return r, err
	} else if err := tx.Commit(); err != nil {
		return r, err
	}
	return r, nil
}

// DeleteResponse permanently deletes a response.
// Returns EUNAUTHORIZED if current response is not the response being deleted.
// Returns ENOTFOUND if response does not exist.
func (s *ResponseService) DeleteResponse(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteResponse(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findResponseByID is a helper function to fetch a response by ID.
// Returns ENOTFOUND if response does not exist.
func findResponseByID(ctx context.Context, tx *Tx, id int) (*api.Response, error) {
	a, _, err := findResponses(ctx, tx, api.ResponseFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Response not found."}
	}
	return a[0], nil
}

// findResponses returns a list of responses matching a filter. Also returns a count of
// total matching responses which may differ if filter.Limit is set.
func findResponses(ctx context.Context, tx *Tx, filter api.ResponseFilter) (_ []*api.Response, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	i := 1
	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.Type; v != nil {
		where, args = append(where, fmt.Sprintf("type = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.IsCorrect; v != nil {
		where, args = append(where, fmt.Sprintf("is_correct = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.SubmissionID; v != nil {
		where, args = append(where, fmt.Sprintf("quiz_submission_id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.QuestionID; v != nil {
		where, args = append(where, fmt.Sprintf("question_id = $%d", i)), append(args, *v)
		i++
	}

	// Execute query to fetch response rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			comments,
			is_correct,
			grade,
			type,
			open_response,
			truefalse_response,
			multiplechoice_response,
			singlechoice_response,
			quiz_submission_id,
			question_id,
		    COUNT(*) OVER()
		FROM responses
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Response objects.
	responses := make([]*api.Response, 0)
	for rows.Next() {
		var isCorrect sql.NullBool
		var grade sql.NullFloat64
		var openResponse sql.NullString
		var truefalseResponse sql.NullBool
		var multiplechoiceResponse pgtype.Int4Array
		var singlechoiceResponse sql.NullInt32

		var r api.Response
		if err := rows.Scan(
			&r.ID,
			&r.Comments,
			&isCorrect,
			&grade,
			&r.Type,
			&openResponse,
			&truefalseResponse,
			&multiplechoiceResponse,
			&singlechoiceResponse,
			&r.SubmissionID,
			&r.QuestionID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		if isCorrect.Valid {
			r.IsCorrect = isCorrect.Bool
		}
		if openResponse.Valid {
			r.OpenResponse = openResponse.String
		}
		if truefalseResponse.Valid {
			r.TrueFalseResponse = truefalseResponse.Bool
		}
		if multiplechoiceResponse.Status != pgtype.Null {
			multiplechoiceResponse.AssignTo(r.MultipleChoiceResponse)
		}
		if singlechoiceResponse.Valid {
			r.SingleChoiceResponse = int(singlechoiceResponse.Int32)
		}

		responses = append(responses, &r)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return responses, n, nil
}

// createResponse creates a new response. Sets the new database ID to r.ID and sets
// the timestamps to the current time.
func createResponse(ctx context.Context, tx *Tx, r *api.Response) error {
	// Perform basic field validation.
	if err := r.Validate(); err != nil {
		return err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var isCorrect *bool
	if r.IsCorrect {
		isCorrect = &r.IsCorrect
	}
	var grade *float32
	if r.Grade != 0 {
		grade = &r.Grade
	}
	var openResponse *string
	if r.OpenResponse != "" {
		openResponse = &r.OpenResponse
	}
	var truefalseResponse *bool
	if r.TrueFalseResponse {
		truefalseResponse = &r.TrueFalseResponse
	}
	var multiplechoiceResponse *[]int
	if r.MultipleChoiceResponse != nil {
		multiplechoiceResponse = &r.MultipleChoiceResponse
	}
	var singlechoiceResponse *int
	if r.SingleChoiceResponse != 0 {
		singlechoiceResponse = &r.SingleChoiceResponse
	}

	// Execute insertion query.
	row := tx.QueryRowContext(ctx, `
		INSERT INTO responses (
			comments,
			is_correct,
			grade,
			type,
			open_response,
			truefalse_response,
			multiplechoice_response,
			singlechoice_response,
			quiz_submission_id,
			question_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
		RETURNING id
	`,
		r.Comments,
		isCorrect,
		grade,
		r.Type,
		openResponse,
		truefalseResponse,
		multiplechoiceResponse,
		singlechoiceResponse,
		r.SubmissionID,
		r.QuestionID,
	)

	err := row.Scan(&r.ID)
	if err != nil {
		return FormatError(err)
	}

	return nil
}

// updateResponse updates fields on a response object. Returns EUNAUTHORIZED if current
// response is not the response being updated.
func updateResponse(ctx context.Context, tx *Tx, id int, upd api.ResponseUpdate) (*api.Response, error) {
	// Fetch current object state.
	currentUserID := api.UserIDFromContext(ctx)
	r, err := findResponseByID(ctx, tx, id)
	if err != nil {
		return r, err
	} else if r.Submission, err = findQuizSubmissionByID(ctx, tx, r.SubmissionID); err != nil {
		return r, err
	} else if r.Submission.Quiz, err = findQuizByID(ctx, tx, r.Submission.QuizID); err != nil {
		return r, err
	} else if currentUserID != 0 && r.Submission.StudentID != currentUserID && r.Submission.Quiz.TeacherID != currentUserID {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this response.")
	}

	// Update fields.
	if v := upd.Comments; v != nil {
		r.Comments = *v
	}
	if v := upd.IsCorrect; v != nil {
		r.IsCorrect = *v
	}
	if v := upd.Grade; v != nil {
		r.Grade = *v
	}
	if v := upd.Type; v != nil {
		r.Type = *v
	}
	if v := upd.OpenResponse; v != nil {
		r.OpenResponse = *v
	}
	if v := upd.TrueFalseResponse; v != nil {
		r.TrueFalseResponse = *v
	}
	if v := upd.MultipleChoiceResponse; v != nil {
		r.MultipleChoiceResponse = *v
	}
	if v := upd.SingleChoiceResponse; v != nil {
		r.SingleChoiceResponse = *v
	}

	// Perform basic field validation.
	if err := r.Validate(); err != nil {
		return r, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var isCorrect *bool
	if r.IsCorrect {
		isCorrect = &r.IsCorrect
	}
	var grade *float32
	if r.Grade != 0 {
		grade = &r.Grade
	}
	var openResponse *string
	if r.OpenResponse != "" {
		openResponse = &r.OpenResponse
	}
	var truefalseResponse *bool
	if r.TrueFalseResponse {
		truefalseResponse = &r.TrueFalseResponse
	}
	var multiplechoiceResponse *[]int
	if len(r.MultipleChoiceResponse) > 0 {
		multiplechoiceResponse = &r.MultipleChoiceResponse
	}
	var singlechoiceResponse *int
	if r.SingleChoiceResponse != 0 {
		singlechoiceResponse = &r.SingleChoiceResponse
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE responses
		SET comments = $1,
		    is_correct = $2,
		    grade = $3,
		    type = $4,
		    open_response = $5,
		    truefalse_response = $6,
		    multiplechoice_response = $7,
		    singlechoice_response = $8,
		WHERE id = $9
	`,
		r.Comments,
		isCorrect,
		grade,
		r.Type,
		openResponse,
		truefalseResponse,
		multiplechoiceResponse,
		singlechoiceResponse,
		id,
	); err != nil {
		return r, FormatError(err)
	}

	return r, nil
}

// deleteResponse permanently removes a response by ID. Returns EUNAUTHORIZED if current
// response is not the one being deleted.
func deleteResponse(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	currentUserID := api.UserIDFromContext(ctx)
	if r, err := findResponseByID(ctx, tx, id); err != nil {
		return err
	} else if r.Submission, err = findQuizSubmissionByID(ctx, tx, r.SubmissionID); err != nil {
		return err
	} else if r.Submission.Quiz, err = findQuizByID(ctx, tx, r.Submission.QuizID); err != nil {
		return err
	} else if currentUserID != 0 && r.Submission.Quiz.TeacherID != currentUserID {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this response.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM responses WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}
