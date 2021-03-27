package pg

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/dori7879/senior-project/api"
)

// Ensure service implements interface.
var _ api.QuizSubmissionService = (*QuizSubmissionService)(nil)

// QuizSubmissionService represents a service for managing submissions.
type QuizSubmissionService struct {
	db *DB
}

// NewQuizSubmissionService returns a new instance of QuizSubmissionService.
func NewQuizSubmissionService(db *DB) *QuizSubmissionService {
	return &QuizSubmissionService{db: db}
}

// FindQuizSubmissionByID retrieves a submission by ID along with their associated quiz and student objects.
// Returns ENOTFOUND if submission does not exist.
func (s *QuizSubmissionService) FindQuizSubmissionByID(ctx context.Context, id int) (*api.QuizSubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch submission and their associated quiz and student objects.
	sub, err := findQuizSubmissionByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachQuizSubmissionAssociations(ctx, tx, sub); err != nil {
		return sub, err
	}
	return sub, nil
}

// FindQuizSubmissions retrieves a list of submissions by filter. Also returns total count of
// matching submissions which may differ from returned results if filter.Limit is specified.
func (s *QuizSubmissionService) FindQuizSubmissions(ctx context.Context, filter api.QuizSubmissionFilter) ([]*api.QuizSubmission, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findQuizSubmissions(ctx, tx, filter)
}

// CreateQuizSubmission creates a new submission.
func (s *QuizSubmissionService) CreateQuizSubmission(ctx context.Context, sub *api.QuizSubmission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new submission object and attach associated quiz and student objects.
	if err := createQuizSubmission(ctx, tx, sub); err != nil {
		return err
	} else if err := attachQuizSubmissionAssociations(ctx, tx, sub); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateQuizSubmission updates a submission object. Returns EUNAUTHORIZED if current submission is
// not the submission that is being updated. Returns ENOTFOUND if submission does not exist.
func (s *QuizSubmissionService) UpdateQuizSubmission(ctx context.Context, id int, upd api.QuizSubmissionUpdate) (*api.QuizSubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update submission & attach associated quiz and student objects.
	sub, err := updateQuizSubmission(ctx, tx, id, upd)
	if err != nil {
		return sub, err
	} else if err := attachQuizSubmissionAssociations(ctx, tx, sub); err != nil {
		return sub, err
	} else if err := tx.Commit(); err != nil {
		return sub, err
	}
	return sub, nil
}

// DeleteQuizSubmission permanently deletes a submission.
// Returns EUNAUTHORIZED if current submission is not the submission being deleted.
// Returns ENOTFOUND if submission does not exist.
func (s *QuizSubmissionService) DeleteQuizSubmission(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteQuizSubmission(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findQuizSubmissionByID is a helper function to fetch a submission by ID.
// Returns ENOTFOUND if submission does not exist.
func findQuizSubmissionByID(ctx context.Context, tx *Tx, id int) (*api.QuizSubmission, error) {
	a, _, err := findQuizSubmissions(ctx, tx, api.QuizSubmissionFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Quiz submission not found."}
	}
	return a[0], nil
}

// findQuizSubmissions returns a list of submissions matching a filter. Also returns a count of
// total matching submissions which may differ if filter.Limit is set.
func findQuizSubmissions(ctx context.Context, tx *Tx, filter api.QuizSubmissionFilter) (_ []*api.QuizSubmission, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	if v := filter.BeforeSubmittedAt; v != nil {
		where, args = append(where, "submitted_at < $2"), append(args, *v)
	}
	if v := filter.AfterUpdatedAt; v != nil {
		where, args = append(where, "updated_at >= $3"), append(args, *v)
	}
	if v := filter.StudentFullName; v != nil {
		where, args = append(where, "student_fullname = $4"), append(args, *v)
	}
	if v := filter.StudentID; v != nil {
		where, args = append(where, "student_id = $5"), append(args, *v)
	}
	if v := filter.QuizID; v != nil {
		where, args = append(where, "quiz_id = $6"), append(args, *v)
	}

	// Execute query to fetch submission rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
			grade,
			comments,
			submitted_at,
			updated_at,
			student_fullname,
			student_id,
			quiz_id,
		    COUNT(*) OVER()
		FROM quiz_submissions
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into QuizSubmission objects.
	submissions := make([]*api.QuizSubmission, 0)
	for rows.Next() {
		var studentFullname sql.NullString
		var grade sql.NullFloat64
		var updatedAt sql.NullTime
		var studentID sql.NullInt32

		var sub api.QuizSubmission
		if rows.Scan(
			&sub.ID,
			&grade,
			&sub.Comments,
			&sub.SubmittedAt,
			&updatedAt,
			&studentFullname,
			&studentID,
			&sub.QuizID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		if grade.Valid {
			sub.Grade = float32(grade.Float64)
		}
		if studentFullname.Valid {
			sub.StudentFullName = studentFullname.String
		}
		if updatedAt.Valid {
			sub.UpdatedAt = updatedAt.Time
		}
		if studentID.Valid {
			sub.StudentID = int(studentID.Int32)
		}

		submissions = append(submissions, &sub)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return submissions, n, nil
}

// createQuizSubmission creates a new submission. Sets the new database ID to sub.ID and sets
// the timestamps to the current time.
func createQuizSubmission(ctx context.Context, tx *Tx, sub *api.QuizSubmission) error {
	// Set timestamps to the current time.
	sub.SubmittedAt = tx.now

	// Perform basic field validation.
	if err := sub.Validate(); err != nil {
		return err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var updatedAt *time.Time
	if !sub.UpdatedAt.IsZero() {
		updatedAt = &sub.UpdatedAt
	}
	var studentFullname *string
	if sub.StudentFullName != "" {
		studentFullname = &sub.StudentFullName
	}
	var studentID *int
	if sub.StudentID != 0 {
		studentID = &sub.StudentID
	}
	var grade *float32
	if sub.Grade != 0 {
		grade = &sub.Grade
	}

	// Execute insertion query.
	result, err := tx.ExecContext(ctx, `
		INSERT INTO quiz_submissions (
			grade,
			comments,
			submitted_at,
			updated_at,
			student_fullname,
			student_id,
			quiz_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`,
		grade,
		sub.Comments,
		sub.SubmittedAt,
		updatedAt,
		studentFullname,
		studentID,
		sub.QuizID,
	)
	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	sub.ID = int(id)

	return nil
}

// updateQuizSubmission updates fields on a submission object. Returns EUNAUTHORIZED if current
// submission is not the submission being updated.
func updateQuizSubmission(ctx context.Context, tx *Tx, id int, upd api.QuizSubmissionUpdate) (*api.QuizSubmission, error) {
	// Fetch current object state.
	sub, err := findQuizSubmissionByID(ctx, tx, id)
	if err != nil {
		return sub, err
	} else if sub.ID != api.UserIDFromContext(ctx) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this quiz submission.")
	}

	// Update fields.
	if v := upd.Grade; v != nil {
		sub.Grade = *v
	}
	if v := upd.Comments; v != nil {
		sub.Comments = *v
	}
	if v := upd.StudentFullName; v != nil {
		sub.StudentFullName = *v
	}
	if v := upd.StudentID; v != nil {
		sub.StudentID = *v
	}

	// Set last updated date to current time.
	sub.UpdatedAt = tx.now

	// Perform basic field validation.
	if err := sub.Validate(); err != nil {
		return sub, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var grade *float32
	if sub.Grade != 0 {
		grade = &sub.Grade
	}
	var studentFullname *string
	if sub.StudentFullName != "" {
		studentFullname = &sub.StudentFullName
	}
	var studentID *int
	if sub.StudentID != 0 {
		studentID = &sub.StudentID
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE quiz_submissions
		SET grade = $1,
		    comments = $2,
		    student_fullname = $3,
		    student_id = $4
		WHERE id = $5
	`,
		grade,
		sub.Comments,
		studentFullname,
		studentID,
		id,
	); err != nil {
		return sub, FormatError(err)
	}

	return sub, nil
}

// deleteQuizSubmission permanently removes a submission by ID. Returns EUNAUTHORIZED if current
// submission is not the one being deleted.
func deleteQuizSubmission(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	if sub, err := findQuizSubmissionByID(ctx, tx, id); err != nil {
		return err
	} else if sub.Quiz, err = findQuizByID(ctx, tx, sub.QuizID); err != nil {
		return err
	} else if sub.Quiz.TeacherID != api.UserIDFromContext(ctx) {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this quiz submission.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM quiz_submissions WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}

// attachQuizSubmissionAssociations attaches quiz and student objects associated with the submission.
func attachQuizSubmissionAssociations(ctx context.Context, tx *Tx, sub *api.QuizSubmission) (err error) {
	if sub.Student, err = findUserByID(ctx, tx, sub.StudentID); err != nil {
		return fmt.Errorf("attach quiz submission user: %w", err)
	} else if sub.Quiz, err = findQuizByID(ctx, tx, sub.QuizID); err != nil {
		return fmt.Errorf("attach quiz submission quiz: %w", err)
	}
	return nil
}
