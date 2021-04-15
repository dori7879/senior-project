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
var _ api.AttSubmissionService = (*AttSubmissionService)(nil)

// AttSubmissionService represents a service for managing submissions.
type AttSubmissionService struct {
	db *DB
}

// NewAttSubmissionService returns a new instance of AttSubmissionService.
func NewAttSubmissionService(db *DB) *AttSubmissionService {
	return &AttSubmissionService{db: db}
}

// FindAttSubmissionByID retrieves a submission by ID along with their associated attendance and student objects.
// Returns ENOTFOUND if submission does not exist.
func (s *AttSubmissionService) FindAttSubmissionByID(ctx context.Context, id int) (*api.AttSubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch submission and their associated attendance and student objects.
	sub, err := findAttSubmissionByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachAttSubmissionAssociations(ctx, tx, sub); err != nil {
		return sub, err
	}
	return sub, nil
}

// FindAttSubmissions retrieves a list of submissions by filter. Also returns total count of
// matching submissions which may differ from returned results if filter.Limit is specified.
func (s *AttSubmissionService) FindAttSubmissions(ctx context.Context, filter api.AttSubmissionFilter) ([]*api.AttSubmission, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findAttSubmissions(ctx, tx, filter)
}

// CreateAttSubmission creates a new submission.
func (s *AttSubmissionService) CreateAttSubmission(ctx context.Context, sub *api.AttSubmission) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new submission object and attach associated attendance and student objects.
	if err := createAttSubmission(ctx, tx, sub); err != nil {
		return err
	} else if err := attachAttSubmissionAssociations(ctx, tx, sub); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateAttSubmission updates a submission object. Returns EUNAUTHORIZED if current submission is
// not the submission that is being updated. Returns ENOTFOUND if submission does not exist.
func (s *AttSubmissionService) UpdateAttSubmission(ctx context.Context, id int, upd api.AttSubmissionUpdate) (*api.AttSubmission, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update submission & attach associated attendance and student objects.
	sub, err := updateAttSubmission(ctx, tx, id, upd)
	if err != nil {
		return sub, err
	} else if err := attachAttSubmissionAssociations(ctx, tx, sub); err != nil {
		return sub, err
	} else if err := tx.Commit(); err != nil {
		return sub, err
	}
	return sub, nil
}

// DeleteAttSubmission permanently deletes a submission.
// Returns EUNAUTHORIZED if current submission is not the submission being deleted.
// Returns ENOTFOUND if submission does not exist.
func (s *AttSubmissionService) DeleteAttSubmission(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteAttSubmission(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findAttSubmissionByID is a helper function to fetch a submission by ID.
// Returns ENOTFOUND if submission does not exist.
func findAttSubmissionByID(ctx context.Context, tx *Tx, id int) (*api.AttSubmission, error) {
	a, _, err := findAttSubmissions(ctx, tx, api.AttSubmissionFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Attendance submission not found."}
	}
	return a[0], nil
}

// findAttSubmissions returns a list of submissions matching a filter. Also returns a count of
// total matching submissions which may differ if filter.Limit is set.
func findAttSubmissions(ctx context.Context, tx *Tx, filter api.AttSubmissionFilter) (_ []*api.AttSubmission, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	i := 1
	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.BeforeSubmittedAt; v != nil {
		where, args = append(where, fmt.Sprintf("submitted_at < $%d", i)), append(args, *v)
		i++
	}
	if v := filter.AfterUpdatedAt; v != nil {
		where, args = append(where, fmt.Sprintf("updated_at >= $%d", i)), append(args, *v)
		i++
	}
	if v := filter.StudentFullName; v != nil {
		where, args = append(where, fmt.Sprintf("student_fullname = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.StudentID; v != nil {
		where, args = append(where, fmt.Sprintf("student_id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.AttendanceID; v != nil {
		where, args = append(where, fmt.Sprintf("attendance_id = $%d", i)), append(args, *v)
		i++
	}

	// Execute query to fetch submission rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    present,
			submitted_at,
			updated_at,
			student_fullname,
			student_id,
			attendance_id,
		    COUNT(*) OVER()
		FROM att_submissions
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into AttSubmission objects.
	submissions := make([]*api.AttSubmission, 0)
	for rows.Next() {
		var studentFullname sql.NullString
		var updatedAt sql.NullTime
		var studentID sql.NullInt32

		var sub api.AttSubmission
		if err := rows.Scan(
			&sub.ID,
			&sub.Present,
			&sub.SubmittedAt,
			&updatedAt,
			&studentFullname,
			&studentID,
			&sub.AttendanceID,
			&n,
		); err != nil {
			return nil, 0, err
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

// createAttSubmission creates a new submission. Sets the new database ID to sub.ID and sets
// the timestamps to the current time.
func createAttSubmission(ctx context.Context, tx *Tx, sub *api.AttSubmission) error {
	var err error

	// Set present to false, then check whether PIN is correctly entered
	sub.Present = false
	if sub.Attendance, err = findAttendanceByID(ctx, tx, sub.AttendanceID); err != nil {
		return err
	} else if sub.PIN == sub.Attendance.PIN {
		sub.Present = true
	}

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

	// Execute insertion query.
	row := tx.QueryRowContext(ctx, `
		INSERT INTO att_submissions (
			present,
			submitted_at,
			updated_at,
			student_fullname,
			student_id,
			attendance_id
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		sub.Present,
		sub.SubmittedAt,
		updatedAt,
		studentFullname,
		studentID,
		sub.AttendanceID,
	)

	err = row.Scan(&sub.ID)
	if err != nil {
		return FormatError(err)
	}

	return nil
}

// updateAttSubmission updates fields on a submission object. Returns EUNAUTHORIZED if current
// submission is not the submission being updated.
func updateAttSubmission(ctx context.Context, tx *Tx, id int, upd api.AttSubmissionUpdate) (*api.AttSubmission, error) {
	// Fetch current object state.
	currentUserID := api.UserIDFromContext(ctx)
	sub, err := findAttSubmissionByID(ctx, tx, id)
	if err != nil {
		return sub, err
	} else if sub.Attendance, err = findAttendanceByID(ctx, tx, sub.AttendanceID); err != nil {
		return sub, err
	} else if (currentUserID != 0 && sub.StudentID != 0 && sub.Attendance.TeacherID != 0) &&
		(sub.StudentID != currentUserID && sub.Attendance.TeacherID != currentUserID) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this attendance submission.")
	}

	// Update fields.
	if v := upd.Present; v != nil {
		sub.Present = *v
		sub.UpdatedAt = tx.now
	}
	if v := upd.StudentFullName; v != nil {
		sub.StudentFullName = *v
	}
	if v := upd.StudentID; v != nil {
		sub.StudentID = *v
	}

	// Perform basic field validation.
	if err := sub.Validate(); err != nil {
		return sub, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var present *bool
	if sub.Present {
		present = &sub.Present
	}
	var studentFullname *string
	if sub.StudentFullName != "" {
		studentFullname = &sub.StudentFullName
	}
	var studentID *int
	if sub.StudentID != 0 {
		studentID = &sub.StudentID
	}
	var updatedAt *time.Time
	if !sub.UpdatedAt.IsZero() {
		updatedAt = &sub.UpdatedAt
	}

	// Execute update query.
	if _, err = tx.ExecContext(ctx, `
		UPDATE att_submissions
		SET present = $1,
			student_fullname = $2,
			student_id = $3,
			updated_at = $4
		WHERE id = $5
	`,
		present,
		studentFullname,
		studentID,
		updatedAt,
		id,
	); err != nil {
		return sub, FormatError(err)
	}

	return sub, nil
}

// deleteAttSubmission permanently removes a submission by ID. Returns EUNAUTHORIZED if current
// submission is not the one being deleted.
func deleteAttSubmission(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	currentUserID := api.UserIDFromContext(ctx)
	if sub, err := findAttSubmissionByID(ctx, tx, id); err != nil {
		return err
	} else if sub.Attendance, err = findAttendanceByID(ctx, tx, sub.AttendanceID); err != nil {
		return err
	} else if (currentUserID != 0 && sub.StudentID != 0 && sub.Attendance.TeacherID != 0) && sub.Attendance.TeacherID != currentUserID {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this attendance submission.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM att_submissions WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}

// attachAttSubmissionAssociations attaches attendance and student objects associated with the submission.
func attachAttSubmissionAssociations(ctx context.Context, tx *Tx, sub *api.AttSubmission) (err error) {
	if sub.Attendance, err = findAttendanceByID(ctx, tx, sub.AttendanceID); err != nil {
		return fmt.Errorf("attach attendance submission attendance: %w", err)
	} else if sub.StudentID == 0 {
		return nil
	} else if sub.Student, err = findUserByID(ctx, tx, sub.StudentID); err != nil {
		return fmt.Errorf("attach attendance submission user: %w", err)
	}
	return nil
}
