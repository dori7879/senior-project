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
var _ api.AttendanceService = (*AttendanceService)(nil)

// AttendanceService represents a service for managing attendances.
type AttendanceService struct {
	db *DB
}

// NewAttendanceService returns a new instance of AttendanceService.
func NewAttendanceService(db *DB) *AttendanceService {
	return &AttendanceService{db: db}
}

// FindAttendanceByID retrieves a attendance by ID along with their associated group and owner objects.
// Returns ENOTFOUND if attendance does not exist.
func (s *AttendanceService) FindAttendanceByID(ctx context.Context, id int) (*api.Attendance, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch attendance and their associated group and owner objects.
	att, err := findAttendanceByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachAttendanceAssociations(ctx, tx, att); err != nil {
		return att, err
	}
	return att, nil
}

// FindAttendanceByStudentLink retrieves a attendance by the student link along with their associated group and owner objects.
// Returns ENOTFOUND if attendance does not exist.
func (s *AttendanceService) FindAttendanceByStudentLink(ctx context.Context, link string) (*api.Attendance, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch attendance and their associated group and owner objects.
	att, err := findAttendanceByStudentLink(ctx, tx, link)
	if err != nil {
		return nil, err
	} else if err := attachAttendanceAssociations(ctx, tx, att); err != nil {
		return att, err
	}
	return att, nil
}

// FindAttendanceByTeacherLink retrieves a attendance by the teacher link along with their associated group and owner objects.
// Returns ENOTFOUND if attendance does not exist.
func (s *AttendanceService) FindAttendanceByTeacherLink(ctx context.Context, link string) (*api.Attendance, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch attendance and their associated group and owner objects.
	att, err := findAttendanceByTeacherLink(ctx, tx, link)
	if err != nil {
		return nil, err
	} else if err := attachAttendanceAssociations(ctx, tx, att); err != nil {
		return att, err
	}
	return att, nil
}

// FindAttendances retrieves a list of attendances by filter. Also returns total count of
// matching attendances which may differ from returned results if filter.Limit is specified.
func (s *AttendanceService) FindAttendances(ctx context.Context, filter api.AttendanceFilter) ([]*api.Attendance, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findAttendances(ctx, tx, filter)
}

// CreateAttendance creates a new attendance.
func (s *AttendanceService) CreateAttendance(ctx context.Context, att *api.Attendance) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new attendance object and attach associated group and owner objects.
	if err := createAttendance(ctx, tx, att); err != nil {
		return err
	} else if err := attachAttendanceAssociations(ctx, tx, att); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateAttendance updates a attendance object. Returns EUNAUTHORIZED if current attendance is
// not the attendance that is being updated. Returns ENOTFOUND if attendance does not exist.
func (s *AttendanceService) UpdateAttendance(ctx context.Context, id int, upd api.AttendanceUpdate) (*api.Attendance, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update attendance & attach associated group and owner objects.
	att, err := updateAttendance(ctx, tx, id, upd)
	if err != nil {
		return att, err
	} else if err := attachAttendanceAssociations(ctx, tx, att); err != nil {
		return att, err
	} else if err := tx.Commit(); err != nil {
		return att, err
	}
	return att, nil
}

// DeleteAttendance permanently deletes a attendance.
// Returns EUNAUTHORIZED if current attendance is not the attendance being deleted.
// Returns ENOTFOUND if attendance does not exist.
func (s *AttendanceService) DeleteAttendance(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteAttendance(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findAttendanceByID is a helper function to fetch a attendance by ID.
// Returns ENOTFOUND if attendance does not exist.
func findAttendanceByID(ctx context.Context, tx *Tx, id int) (*api.Attendance, error) {
	a, _, err := findAttendances(ctx, tx, api.AttendanceFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Attendance not found."}
	}
	return a[0], nil
}

// findAttendanceByStudentLink is a helper function to fetch a attendance by the student link.
// Returns ENOTFOUND if attendance does not exist.
func findAttendanceByStudentLink(ctx context.Context, tx *Tx, studentlink string) (*api.Attendance, error) {
	a, _, err := findAttendances(ctx, tx, api.AttendanceFilter{StudentLink: &studentlink})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Attendance not found."}
	}
	return a[0], nil
}

// findAttendanceByTeacherLink is a helper function to fetch a attendance by the teacher link.
// Returns ENOTFOUND if attendance does not exist.
func findAttendanceByTeacherLink(ctx context.Context, tx *Tx, teacherlink string) (*api.Attendance, error) {
	a, _, err := findAttendances(ctx, tx, api.AttendanceFilter{TeacherLink: &teacherlink})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Attendance not found."}
	}
	return a[0], nil
}

// findAttendances returns a list of attendances matching a filter. Also returns a count of
// total matching attendances which may differ if filter.Limit is set.
func findAttendances(ctx context.Context, tx *Tx, filter api.AttendanceFilter) (_ []*api.Attendance, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	i := 1
	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.Title; v != nil {
		where, args = append(where, fmt.Sprintf("title = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.StudentLink; v != nil {
		where, args = append(where, fmt.Sprintf("student_link = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.TeacherLink; v != nil {
		where, args = append(where, fmt.Sprintf("teacher_link = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.CourseTitle; v != nil {
		where, args = append(where, fmt.Sprintf("course_title = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.Mode; v != nil {
		where, args = append(where, fmt.Sprintf("mode = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.TeacherFullName; v != nil {
		where, args = append(where, fmt.Sprintf("teacher_fullname = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.TeacherID; v != nil {
		where, args = append(where, fmt.Sprintf("teacher_id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.GroupID; v != nil {
		where, args = append(where, fmt.Sprintf("group_id = $%d", i)), append(args, *v)
		i++
	}

	// Execute query to fetch attendance rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    title,
			pin,
			student_link,
			teacher_link,
			course_title,
			mode,
			created_at,
			updated_at,
			opened_at,
			closed_at,
			teacher_fullname,
			teacher_id,
			group_id,
		    COUNT(*) OVER()
		FROM attendances
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Attendance objects.
	attendances := make([]*api.Attendance, 0)
	for rows.Next() {
		var teacherFullname sql.NullString
		var updatedAt sql.NullTime
		var openedAt sql.NullTime
		var closedAt sql.NullTime
		var teacherID sql.NullInt32
		var groupID sql.NullInt32

		var att api.Attendance
		if err := rows.Scan(
			&att.ID,
			&att.Title,
			&att.PIN,
			&att.StudentLink,
			&att.TeacherLink,
			&att.CourseTitle,
			&att.Mode,
			&att.CreatedAt,
			&updatedAt,
			&openedAt,
			&closedAt,
			&teacherFullname,
			&teacherID,
			&groupID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		if teacherFullname.Valid {
			att.TeacherFullName = teacherFullname.String
		}
		if updatedAt.Valid {
			att.UpdatedAt = updatedAt.Time
		}
		if openedAt.Valid {
			att.OpenedAt = openedAt.Time
		}
		if closedAt.Valid {
			att.ClosedAt = closedAt.Time
		}
		if teacherID.Valid {
			att.TeacherID = int(teacherID.Int32)
		}
		if groupID.Valid {
			att.GroupID = int(groupID.Int32)
		}

		attendances = append(attendances, &att)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return attendances, n, nil
}

// createAttendance creates a new attendance. Sets the new database ID to att.ID and sets
// the timestamps to the current time.
func createAttendance(ctx context.Context, tx *Tx, att *api.Attendance) error {
	// Set timestamps to the current time.
	att.CreatedAt = tx.now
	att.UpdatedAt = att.CreatedAt

	// Perform basic field validation.
	if err := att.Validate(); err != nil {
		return err
	}

	// Content is nullable so ensure we store blank fields as NULLs.
	var openedAt *time.Time
	if !att.OpenedAt.IsZero() {
		openedAt = &att.OpenedAt
	}
	var closedAt *time.Time
	if !att.ClosedAt.IsZero() {
		closedAt = &att.ClosedAt
	}
	var teacherFullname *string
	if att.TeacherFullName != "" {
		teacherFullname = &att.TeacherFullName
	}
	var teacherID *int
	if att.TeacherID != 0 {
		teacherID = &att.TeacherID
	}
	var groupID *int
	if att.GroupID != 0 {
		groupID = &att.GroupID
	}

	// Execute insertion query.
	row := tx.QueryRowContext(ctx, `
		INSERT INTO attendances (
			title,
			pin,
			student_link,
			teacher_link,
			course_title,
			mode,
			created_at,
			opened_at,
			closed_at,
			teacher_fullname,
			teacher_id,
			group_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)
		RETURNING id
	`,
		att.Title,
		att.PIN,
		att.StudentLink,
		att.TeacherLink,
		att.CourseTitle,
		att.Mode,
		att.CreatedAt,
		openedAt,
		closedAt,
		teacherFullname,
		teacherID,
		groupID,
	)

	err := row.Scan(&att.ID)
	if err != nil {
		return FormatError(err)
	}

	return nil
}

// updateAttendance updates fields on a attendance object. Returns EUNAUTHORIZED if current
// attendance is not the attendance being updated.
func updateAttendance(ctx context.Context, tx *Tx, id int, upd api.AttendanceUpdate) (*api.Attendance, error) {
	// Fetch current object state.
	currentUserID := api.UserIDFromContext(ctx)
	att, err := findAttendanceByID(ctx, tx, id)
	if err != nil {
		return att, err
	} else if currentUserID != 0 && att.TeacherID != 0 && att.TeacherID != currentUserID {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this attendance.")
	}

	// Update fields.
	if v := upd.Title; v != nil {
		att.Title = *v
	}
	if v := upd.CourseTitle; v != nil {
		att.CourseTitle = *v
	}
	if v := upd.Mode; v != nil {
		att.Mode = *v
	}
	if v := upd.PIN; v != nil {
		att.PIN = *v
	}
	if v := upd.OpenedAt; v != nil {
		att.OpenedAt = *v
	}
	if v := upd.ClosedAt; v != nil {
		att.ClosedAt = *v
	}
	if v := upd.TeacherFullName; v != nil {
		att.TeacherFullName = *v
	}
	if v := upd.TeacherID; v != nil {
		att.TeacherID = *v
	}
	if v := upd.GroupID; v != nil {
		att.GroupID = *v
	}

	// Set last updated date to current time.
	att.UpdatedAt = tx.now

	// Perform basic field validation.
	if err := att.Validate(); err != nil {
		return att, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var openedAt *time.Time
	if !att.OpenedAt.IsZero() {
		openedAt = &att.OpenedAt
	}
	var closedAt *time.Time
	if !att.ClosedAt.IsZero() {
		closedAt = &att.ClosedAt
	}
	var teacherFullname *string
	if att.TeacherFullName != "" {
		teacherFullname = &att.TeacherFullName
	}
	var teacherID *int
	if att.TeacherID != 0 {
		teacherID = &att.TeacherID
	}
	var groupID *int
	if att.GroupID != 0 {
		groupID = &att.GroupID
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE attendances
		SET title = $1,
		    pin = $2,
		    course_title = $3,
		    mode = $4,
		    opened_at = $5,
		    closed_at = $6,
		    teacher_fullname = $7,
		    teacher_id = $8,
		    group_id = $9
		WHERE id = $10
	`,
		att.Title,
		att.PIN,
		att.CourseTitle,
		att.Mode,
		openedAt,
		closedAt,
		teacherFullname,
		teacherID,
		groupID,
		id,
	); err != nil {
		return att, FormatError(err)
	}

	return att, nil
}

// deleteAttendance permanently removes a attendance by ID. Returns EUNAUTHORIZED if current
// attendance is not the one being deleted.
func deleteAttendance(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	currentUserID := api.UserIDFromContext(ctx)
	if att, err := findAttendanceByID(ctx, tx, id); err != nil {
		return err
	} else if currentUserID != 0 && att.TeacherID != 0 && att.TeacherID != currentUserID {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this attendance.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM attendances WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}

// attachAttendanceAssociations attaches group and owner objects associated with the attendance.
func attachAttendanceAssociations(ctx context.Context, tx *Tx, att *api.Attendance) (err error) {
	if att.TeacherID == 0 {
		return nil
	} else if att.Teacher, err = findUserByID(ctx, tx, att.TeacherID); err != nil {
		return fmt.Errorf("attach attendance user: %w", err)
	} else if att.GroupID == 0 {
		return nil
	} else if att.Group, err = findGroupByID(ctx, tx, att.GroupID); err != nil {
		return fmt.Errorf("attach attendance group: %w", err)
	}
	return nil
}
