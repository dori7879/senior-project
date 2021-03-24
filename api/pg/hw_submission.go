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
var _ api.HWSubmissionService = (*HWSubmissionService)(nil)

// HWSubmissionService represents a service for managing submissions.
type HWSubmissionService struct {
	db *DB
}

// NewHWSubmissionService returns a new instance of HWSubmissionService.
func NewHWSubmissionService(db *DB) *HWSubmissionService {
	return &HWSubmissionService{db: db}
}

// FindHWSubmissionByID retrieves a submission by ID along with their associated homework and student objects.
// Returns ENOTFOUND if submission does not exist.
func (s *HWSubmissionService) FindHWSubmissionByID(ctx context.Context, id int) (*api.HWSubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch submission and their associated homework and student objects.
	sub, err := findHWSubmissionByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachHWSubmissionAssociations(ctx, tx, sub); err != nil {
		return sub, err
	}
	return sub, nil
}

// FindHWSubmissions retrieves a list of submissions by filter. Also returns total count of
// matching submissions which may differ from returned results if filter.Limit is specified.
func (s *HWSubmissionService) FindHWSubmissions(ctx context.Context, filter api.HWSubmissionFilter) ([]*api.HWSubmission, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findHWSubmissions(ctx, tx, filter)
}

// CreateHWSubmission creates a new submission.
func (s *HWSubmissionService) CreateHWSubmission(ctx context.Context, sub *api.HWSubmission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new submission object and attach associated homework and student objects.
	if err := createHWSubmission(ctx, tx, sub); err != nil {
		return err
	} else if err := attachHWSubmissionAssociations(ctx, tx, sub); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateHWSubmission updates a submission object. Returns EUNAUTHORIZED if current submission is
// not the submission that is being updated. Returns ENOTFOUND if submission does not exist.
func (s *HWSubmissionService) UpdateHWSubmission(ctx context.Context, id int, upd api.HWSubmissionUpdate) (*api.HWSubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update submission & attach associated homework and student objects.
	sub, err := updateHWSubmission(ctx, tx, id, upd)
	if err != nil {
		return sub, err
	} else if err := attachHWSubmissionAssociations(ctx, tx, sub); err != nil {
		return sub, err
	} else if err := tx.Commit(); err != nil {
		return sub, err
	}
	return sub, nil
}

// DeleteHWSubmission permanently deletes a submission.
// Returns EUNAUTHORIZED if current submission is not the submission being deleted.
// Returns ENOTFOUND if submission does not exist.
func (s *HWSubmissionService) DeleteHWSubmission(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteHWSubmission(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findHWSubmissionByID is a helper function to fetch a submission by ID.
// Returns ENOTFOUND if submission does not exist.
func findHWSubmissionByID(ctx context.Context, tx *Tx, id int) (*api.HWSubmission, error) {
	a, _, err := findHWSubmissions(ctx, tx, api.HWSubmissionFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Homework submission not found."}
	}
	return a[0], nil
}

// findHWSubmissions returns a list of submissions matching a filter. Also returns a count of
// total matching submissions which may differ if filter.Limit is set.
func findHWSubmissions(ctx context.Context, tx *Tx, filter api.HWSubmissionFilter) (_ []*api.HWSubmission, n int, err error) {
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
	if v := filter.HomeworkID; v != nil {
		where, args = append(where, "homework_id = $6"), append(args, *v)
	}

	// Execute query to fetch submission rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    response,
			grade,
			comments,
			submitted_at,
			updated_at,
			student_fullname,
			student_id,
			homework_id,
		    COUNT(*) OVER()
		FROM hw_submissions
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into HWSubmission objects.
	submissions := make([]*api.HWSubmission, 0)
	for rows.Next() {
		var response sql.NullString
		var studentFullname sql.NullString
		var grade sql.NullFloat64
		var updatedAt sql.NullTime
		var studentID sql.NullInt32

		var sub api.HWSubmission
		if rows.Scan(
			&sub.ID,
			&response,
			&grade,
			&sub.Comments,
			&sub.SubmittedAt,
			&updatedAt,
			&studentFullname,
			&studentID,
			&sub.HomeworkID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		if response.Valid {
			sub.Response = response.String
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

// createHWSubmission creates a new submission. Sets the new database ID to sub.ID and sets
// the timestamps to the current time.
func createHWSubmission(ctx context.Context, tx *Tx, sub *api.HWSubmission) error {
	// Set timestamps to the current time.
	sub.SubmittedAt = tx.now

	// Perform basic field validation.
	if err := sub.Validate(); err != nil {
		return err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var response *string
	if sub.Response != "" {
		response = &sub.Response
	}
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
		INSERT INTO hw_submissions (
			response,
			grade,
			comments,
			submitted_at,
			updated_at,
			student_fullname,
			student_id,
			homework_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`,
		response,
		grade,
		sub.Comments,
		sub.SubmittedAt,
		updatedAt,
		studentFullname,
		studentID,
		sub.HomeworkID,
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

// updateHWSubmission updates fields on a submission object. Returns EUNAUTHORIZED if current
// submission is not the submission being updated.
func updateHWSubmission(ctx context.Context, tx *Tx, id int, upd api.HWSubmissionUpdate) (*api.HWSubmission, error) {
	// Fetch current object state.
	sub, err := findHWSubmissionByID(ctx, tx, id)
	if err != nil {
		return sub, err
	} else if sub.ID != api.UserIDFromContext(ctx) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this homework submission.")
	}

	// Update fields.
	if v := upd.Response; v != nil {
		sub.Response = *v
	}
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
	var response *string
	if sub.Response != "" {
		response = &sub.Response
	}
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
		UPDATE hw_submissions
		SET response = $1,
		    grade = $2,
		    comments = $3,
		    student_fullname = $4,
		    student_id = $5
		WHERE id = $6
	`,
		response,
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

// deleteHWSubmission permanently removes a submission by ID. Returns EUNAUTHORIZED if current
// submission is not the one being deleted.
func deleteHWSubmission(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	if sub, err := findHWSubmissionByID(ctx, tx, id); err != nil {
		return err
	} else if sub.ID != api.UserIDFromContext(ctx) {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this homework submission.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM hw_submissions WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}

// attachHWSubmissionAssociations attaches homework and student objects associated with the submission.
func attachHWSubmissionAssociations(ctx context.Context, tx *Tx, sub *api.HWSubmission) (err error) {
	if sub.Student, err = findUserByID(ctx, tx, sub.StudentID); err != nil {
		return fmt.Errorf("attach homework submission user: %w", err)
	} else if sub.Homework, err = findHomeworkByID(ctx, tx, sub.HomeworkID); err != nil {
		return fmt.Errorf("attach homework submission homework: %w", err)
	}
	return nil
}
