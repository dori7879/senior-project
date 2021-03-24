package pg

import (
	"context"
	"database/sql"
	"strings"
	"time"

	"github.com/dori7879/senior-project/api"
	"github.com/jackc/pgtype"
)

// Ensure service implements interface.
var _ api.QuestionService = (*QuestionService)(nil)

// QuestionService represents a service for managing questions.
type QuestionService struct {
	db *DB
}

// NewQuestionService returns a new instance of QuestionService.
func NewQuestionService(db *DB) *QuestionService {
	return &QuestionService{db: db}
}

// FindQuestionByID retrieves a question by ID.
// Returns ENOTFOUND if question does not exist.
func (s *QuestionService) FindQuestionByID(ctx context.Context, id int) (*api.Question, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch question.
	q, err := findQuestionByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return q, nil
}

// FindQuestions retrieves a list of questions by filter. Also returns total count of
// matching questions which may differ from returned results if filter.Limit is specified.
func (s *QuestionService) FindQuestions(ctx context.Context, filter api.QuestionFilter) ([]*api.Question, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findQuestions(ctx, tx, filter)
}

// CreateQuestion creates a new question.
func (s *QuestionService) CreateQuestion(ctx context.Context, q *api.Question) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new question object.
	if err := createQuestion(ctx, tx, q); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateQuestion updates a question object. Returns EUNAUTHORIZED if current question is
// not the question that is being updated. Returns ENOTFOUND if question does not exist.
func (s *QuestionService) UpdateQuestion(ctx context.Context, id int, upd api.QuestionUpdate) (*api.Question, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update question.
	q, err := updateQuestion(ctx, tx, id, upd)
	if err != nil {
		return q, err
	} else if err := tx.Commit(); err != nil {
		return q, err
	}
	return q, nil
}

// DeleteQuestion permanently deletes a question.
// Returns EUNAUTHORIZED if current question is not the question being deleted.
// Returns ENOTFOUND if question does not exist.
func (s *QuestionService) DeleteQuestion(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteQuestion(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findQuestionByID is a helper function to fetch a question by ID.
// Returns ENOTFOUND if question does not exist.
func findQuestionByID(ctx context.Context, tx *Tx, id int) (*api.Question, error) {
	a, _, err := findQuestions(ctx, tx, api.QuestionFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Question not found."}
	}
	return a[0], nil
}

// findQuestions returns a list of questions matching a filter. Also returns a count of
// total matching questions which may differ if filter.Limit is set.
func findQuestions(ctx context.Context, tx *Tx, filter api.QuestionFilter) (_ []*api.Question, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	if v := filter.Type; v != nil {
		where, args = append(where, "type = $2"), append(args, *v)
	}
	if v := filter.Fixed; v != nil {
		where, args = append(where, "fixed = $3"), append(args, *v)
	}
	if v := filter.QuizID; v != nil {
		where, args = append(where, "quiz_id = $4"), append(args, *v)
	}

	// Execute query to fetch question rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			content,
			type,
			fixed,
			choices,
			open_answer,
			truefalse_answer,
			multiplechoice_answer,
			singlechoice_answer,
			created_at,
			updated_at,
			quiz_id,
		    COUNT(*) OVER()
		FROM questions
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Question objects.
	questions := make([]*api.Question, 0)
	for rows.Next() {
		// var choices pgtype.VarcharArray
		var openAnswer sql.NullString
		var truefalseAnswer sql.NullBool
		var multiplechoiceAnswer pgtype.Int4Array
		var singlechoiceAnswer sql.NullInt32
		var updatedAt sql.NullTime

		var q api.Question
		if rows.Scan(
			&q.ID,
			&q.Content,
			&q.Type,
			&q.Fixed,
			&q.Choices,
			&openAnswer,
			&truefalseAnswer,
			&multiplechoiceAnswer,
			&singlechoiceAnswer,
			&q.CreatedAt,
			&updatedAt,
			&q.QuizID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		if openAnswer.Valid {
			q.OpenAnswer = openAnswer.String
		}
		if truefalseAnswer.Valid {
			q.TrueFalseAnswer = truefalseAnswer.Bool
		}
		if multiplechoiceAnswer.Status != pgtype.Null {
			multiplechoiceAnswer.AssignTo(q.MultipleChoiceAnswer)
		}
		if updatedAt.Valid {
			q.UpdatedAt = updatedAt.Time
		}

		questions = append(questions, &q)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return questions, n, nil
}

// createQuestion creates a new question. Sets the new database ID to q.ID and sets
// the timestamps to the current time.
func createQuestion(ctx context.Context, tx *Tx, q *api.Question) error {
	// Set timestamps to the current time.
	q.CreatedAt = tx.now
	q.UpdatedAt = q.CreatedAt

	// Perform basic field validation.
	if err := q.Validate(); err != nil {
		return err
	}

	// Content is nullable so ensure we store blank fields as NULLs.
	var updatedAt *time.Time
	if !q.UpdatedAt.IsZero() {
		updatedAt = &q.UpdatedAt
	}
	var openAnswer *string
	if q.OpenAnswer != "" {
		openAnswer = &q.OpenAnswer
	}
	var truefalseAnswer *bool
	if q.TrueFalseAnswer {
		truefalseAnswer = &q.TrueFalseAnswer
	}
	var multiplechoiceAnswer *[]int
	if q.MultipleChoiceAnswer != nil {
		multiplechoiceAnswer = &q.MultipleChoiceAnswer
	}
	var singlechoiceAnswer *int
	if q.SingleChoiceAnswer != 0 {
		singlechoiceAnswer = &q.SingleChoiceAnswer
	}

	// Execute insertion query.
	result, err := tx.ExecContext(ctx, `
		INSERT INTO questions (
			content,
			type,
			fixed,
			choices,
			open_answer,
			truefalse_answer,
			multiplechoice_answer,
			singlechoice_answer,
			created_at,
			updated_at,
			quiz_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
	`,
		q.Content,
		q.Type,
		q.Fixed,
		q.Choices,
		openAnswer,
		truefalseAnswer,
		multiplechoiceAnswer,
		singlechoiceAnswer,
		q.CreatedAt,
		updatedAt,
		q.QuizID,
	)
	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	q.ID = int(id)

	return nil
}

// updateQuestion updates fields on a question object. Returns EUNAUTHORIZED if current
// question is not the question being updated.
func updateQuestion(ctx context.Context, tx *Tx, id int, upd api.QuestionUpdate) (*api.Question, error) {
	// Fetch current object state.
	q, err := findQuestionByID(ctx, tx, id)
	if err != nil {
		return q, err
	} else if q.ID != api.UserIDFromContext(ctx) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this question.")
	}

	// Update fields.
	if v := upd.Content; v != nil {
		q.Content = *v
	}
	if v := upd.Type; v != nil {
		q.Type = *v
	}
	if v := upd.Fixed; v != nil {
		q.Fixed = *v
	}
	if v := upd.Choices; v != nil {
		q.Choices = *v
	}
	if v := upd.OpenAnswer; v != nil {
		q.OpenAnswer = *v
	}
	if v := upd.TrueFalseAnswer; v != nil {
		q.TrueFalseAnswer = *v
	}
	if v := upd.MultipleChoiceAnswer; v != nil {
		q.MultipleChoiceAnswer = *v
	}
	if v := upd.SingleChoiceAnswer; v != nil {
		q.SingleChoiceAnswer = *v
	}

	// Set last updated date to current time.
	q.UpdatedAt = tx.now

	// Perform basic field validation.
	if err := q.Validate(); err != nil {
		return q, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var openAnswer *string
	if q.OpenAnswer != "" {
		openAnswer = &q.OpenAnswer
	}
	var truefalseAnswer *bool
	if q.TrueFalseAnswer {
		truefalseAnswer = &q.TrueFalseAnswer
	}
	var multiplechoiceAnswer *[]int
	if len(q.MultipleChoiceAnswer) > 0 {
		multiplechoiceAnswer = &q.MultipleChoiceAnswer
	}
	var singlechoiceAnswer *int
	if q.SingleChoiceAnswer != 0 {
		singlechoiceAnswer = &q.SingleChoiceAnswer
	}
	var updatedAt *time.Time
	if !q.UpdatedAt.IsZero() {
		updatedAt = &q.UpdatedAt
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE questions
		SET content = $1,
		    type = $2,
		    fixed = $3,
		    choices = $4,
		    open_answer = $5,
		    truefalse_answer = $6,
		    multiplechoice_answer = $7,
		    singlechoice_answer = $8,
		    updated_at = $9,
		WHERE id = $10
	`,
		q.Content,
		q.Type,
		q.Fixed,
		q.Choices,
		openAnswer,
		truefalseAnswer,
		multiplechoiceAnswer,
		singlechoiceAnswer,
		updatedAt,
		id,
	); err != nil {
		return q, FormatError(err)
	}

	return q, nil
}

// deleteQuestion permanently removes a question by ID. Returns EUNAUTHORIZED if current
// question is not the one being deleted.
func deleteQuestion(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	if q, err := findQuestionByID(ctx, tx, id); err != nil {
		return err
	} else if q.ID != api.UserIDFromContext(ctx) {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this question.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM questions WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}
