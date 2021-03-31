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
var _ api.HomeworkService = (*HomeworkService)(nil)

// HomeworkService represents a service for managing homeworks.
type HomeworkService struct {
	db *DB
}

// NewHomeworkService returns a new instance of HomeworkService.
func NewHomeworkService(db *DB) *HomeworkService {
	return &HomeworkService{db: db}
}

// FindHomeworkByID retrieves a homework by ID along with their associated group and owner objects.
// Returns ENOTFOUND if homework does not exist.
func (s *HomeworkService) FindHomeworkByID(ctx context.Context, id int) (*api.Homework, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch homework and their associated group and owner objects.
	hw, err := findHomeworkByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachHomeworkAssociations(ctx, tx, hw); err != nil {
		return hw, err
	}
	return hw, nil
}

// FindHomeworkByStudentLink retrieves a homework by the student link along with their associated group and owner objects.
// Returns ENOTFOUND if homework does not exist.
func (s *HomeworkService) FindHomeworkByStudentLink(ctx context.Context, link string) (*api.Homework, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch homework and their associated group and owner objects.
	hw, err := findHomeworkByStudentLink(ctx, tx, link)
	if err != nil {
		return nil, err
	} else if err := attachHomeworkAssociations(ctx, tx, hw); err != nil {
		return hw, err
	}
	return hw, nil
}

// FindHomeworkByTeacherLink retrieves a homework by the teacher link along with their associated group and owner objects.
// Returns ENOTFOUND if homework does not exist.
func (s *HomeworkService) FindHomeworkByTeacherLink(ctx context.Context, link string) (*api.Homework, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch homework and their associated group and owner objects.
	hw, err := findHomeworkByTeacherLink(ctx, tx, link)
	if err != nil {
		return nil, err
	} else if err := attachHomeworkAssociations(ctx, tx, hw); err != nil {
		return hw, err
	}
	return hw, nil
}

// FindHomeworks retrieves a list of homeworks by filter. Also returns total count of
// matching homeworks which may differ from returned results if filter.Limit is specified.
func (s *HomeworkService) FindHomeworks(ctx context.Context, filter api.HomeworkFilter) ([]*api.Homework, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findHomeworks(ctx, tx, filter)
}

// CreateHomework creates a new homework.
func (s *HomeworkService) CreateHomework(ctx context.Context, hw *api.Homework) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new homework object and attach associated group and owner objects.
	if err := createHomework(ctx, tx, hw); err != nil {
		return err
	} else if err := attachHomeworkAssociations(ctx, tx, hw); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateHomework updates a homework object. Returns EUNAUTHORIZED if current homework is
// not the homework that is being updated. Returns ENOTFOUND if homework does not exist.
func (s *HomeworkService) UpdateHomework(ctx context.Context, id int, upd api.HomeworkUpdate) (*api.Homework, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update homework & attach associated group and owner objects.
	hw, err := updateHomework(ctx, tx, id, upd)
	if err != nil {
		return hw, err
	} else if err := attachHomeworkAssociations(ctx, tx, hw); err != nil {
		return hw, err
	} else if err := tx.Commit(); err != nil {
		return hw, err
	}
	return hw, nil
}

// DeleteHomework permanently deletes a homework.
// Returns EUNAUTHORIZED if current homework is not the homework being deleted.
// Returns ENOTFOUND if homework does not exist.
func (s *HomeworkService) DeleteHomework(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteHomework(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findHomeworkByID is a helper function to fetch a homework by ID.
// Returns ENOTFOUND if homework does not exist.
func findHomeworkByID(ctx context.Context, tx *Tx, id int) (*api.Homework, error) {
	a, _, err := findHomeworks(ctx, tx, api.HomeworkFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Homework not found."}
	}
	return a[0], nil
}

// findHomeworkByStudentLink is a helper function to fetch a homework by the student link.
// Returns ENOTFOUND if homework does not exist.
func findHomeworkByStudentLink(ctx context.Context, tx *Tx, studentlink string) (*api.Homework, error) {
	a, _, err := findHomeworks(ctx, tx, api.HomeworkFilter{StudentLink: &studentlink})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Homework not found."}
	}
	return a[0], nil
}

// findHomeworkByTeacherLink is a helper function to fetch a homework by the teacher link.
// Returns ENOTFOUND if homework does not exist.
func findHomeworkByTeacherLink(ctx context.Context, tx *Tx, teacherlink string) (*api.Homework, error) {
	a, _, err := findHomeworks(ctx, tx, api.HomeworkFilter{TeacherLink: &teacherlink})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Homework not found."}
	}
	return a[0], nil
}

// findHomeworks returns a list of homeworks matching a filter. Also returns a count of
// total matching homeworks which may differ if filter.Limit is set.
func findHomeworks(ctx context.Context, tx *Tx, filter api.HomeworkFilter) (_ []*api.Homework, n int, err error) {
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

	// Execute query to fetch homework rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    title,
			content,
			max_grade,
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
		FROM homeworks
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Homework objects.
	homeworks := make([]*api.Homework, 0)
	for rows.Next() {
		var content sql.NullString
		var teacherFullname sql.NullString
		var updatedAt sql.NullTime
		var openedAt sql.NullTime
		var closedAt sql.NullTime
		var teacherID sql.NullInt32
		var groupID sql.NullInt32

		var hw api.Homework
		if err := rows.Scan(
			&hw.ID,
			&hw.Title,
			&content,
			&hw.MaxGrade,
			&hw.StudentLink,
			&hw.TeacherLink,
			&hw.CourseTitle,
			&hw.Mode,
			&hw.CreatedAt,
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

		if content.Valid {
			hw.Content = content.String
		}
		if teacherFullname.Valid {
			hw.TeacherFullName = teacherFullname.String
		}
		if updatedAt.Valid {
			hw.UpdatedAt = updatedAt.Time
		}
		if openedAt.Valid {
			hw.OpenedAt = openedAt.Time
		}
		if closedAt.Valid {
			hw.ClosedAt = closedAt.Time
		}
		if teacherID.Valid {
			hw.TeacherID = int(teacherID.Int32)
		}
		if groupID.Valid {
			hw.GroupID = int(groupID.Int32)
		}

		homeworks = append(homeworks, &hw)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return homeworks, n, nil
}

// createHomework creates a new homework. Sets the new database ID to hw.ID and sets
// the timestamps to the current time.
func createHomework(ctx context.Context, tx *Tx, hw *api.Homework) error {
	// Set timestamps to the current time.
	hw.CreatedAt = tx.now
	hw.UpdatedAt = hw.CreatedAt

	// Perform basic field validation.
	if err := hw.Validate(); err != nil {
		return err
	}

	// Content is nullable so ensure we store blank fields as NULLs.
	var content *string
	if hw.Content != "" {
		content = &hw.Content
	}
	var openedAt *time.Time
	if !hw.OpenedAt.IsZero() {
		openedAt = &hw.OpenedAt
	}
	var closedAt *time.Time
	if !hw.ClosedAt.IsZero() {
		closedAt = &hw.ClosedAt
	}
	var teacherFullname *string
	if hw.TeacherFullName != "" {
		teacherFullname = &hw.TeacherFullName
	}
	var teacherID *int
	if hw.TeacherID != 0 {
		teacherID = &hw.TeacherID
	}
	var groupID *int
	if hw.GroupID != 0 {
		groupID = &hw.GroupID
	}

	// Execute insertion query.
	row := tx.QueryRowContext(ctx, `
		INSERT INTO homeworks (
			title,
			content,
			max_grade,
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
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
		RETURNING id
	`,
		hw.Title,
		content,
		hw.MaxGrade,
		hw.StudentLink,
		hw.TeacherLink,
		hw.CourseTitle,
		hw.Mode,
		hw.CreatedAt,
		openedAt,
		closedAt,
		teacherFullname,
		teacherID,
		groupID,
	)

	err := row.Scan(&hw.ID)
	if err != nil {
		return FormatError(err)
	}

	return nil
}

// updateHomework updates fields on a homework object. Returns EUNAUTHORIZED if current
// homework is not the homework being updated.
func updateHomework(ctx context.Context, tx *Tx, id int, upd api.HomeworkUpdate) (*api.Homework, error) {
	// Fetch current object state.
	currentUserID := api.UserIDFromContext(ctx)
	hw, err := findHomeworkByID(ctx, tx, id)
	if err != nil {
		return hw, err
	} else if currentUserID != 0 && hw.TeacherID != 0 && hw.TeacherID != currentUserID {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this homework.")
	}

	// Update fields.
	if v := upd.Title; v != nil {
		hw.Title = *v
	}
	if v := upd.Content; v != nil {
		hw.Content = *v
	}
	if v := upd.MaxGrade; v != nil {
		hw.MaxGrade = *v
	}
	if v := upd.CourseTitle; v != nil {
		hw.CourseTitle = *v
	}
	if v := upd.Mode; v != nil {
		hw.Mode = *v
	}
	if v := upd.OpenedAt; v != nil {
		hw.OpenedAt = *v
	}
	if v := upd.ClosedAt; v != nil {
		hw.ClosedAt = *v
	}
	if v := upd.TeacherFullName; v != nil {
		hw.TeacherFullName = *v
	}
	if v := upd.TeacherID; v != nil {
		hw.TeacherID = *v
	}
	if v := upd.GroupID; v != nil {
		hw.GroupID = *v
	}

	// Set last updated date to current time.
	hw.UpdatedAt = tx.now

	// Perform basic field validation.
	if err := hw.Validate(); err != nil {
		return hw, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var content *string
	if hw.Content != "" {
		content = &hw.Content
	}
	var openedAt *time.Time
	if !hw.OpenedAt.IsZero() {
		openedAt = &hw.OpenedAt
	}
	var closedAt *time.Time
	if !hw.ClosedAt.IsZero() {
		closedAt = &hw.ClosedAt
	}
	var teacherFullname *string
	if hw.TeacherFullName != "" {
		teacherFullname = &hw.TeacherFullName
	}
	var teacherID *int
	if hw.TeacherID != 0 {
		teacherID = &hw.TeacherID
	}
	var groupID *int
	if hw.GroupID != 0 {
		groupID = &hw.GroupID
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE homeworks
		SET title = $1,
		    content = $2,
		    max_grade = $3,
		    course_title = $4,
		    mode = $5,
		    opened_at = $6,
		    closed_at = $7,
		    teacher_fullname = $8,
		    teacher_id = $9,
		    group_id = $10
		WHERE id = $11
	`,
		hw.Title,
		content,
		hw.MaxGrade,
		hw.CourseTitle,
		hw.Mode,
		openedAt,
		closedAt,
		teacherFullname,
		teacherID,
		groupID,
		id,
	); err != nil {
		return hw, FormatError(err)
	}

	return hw, nil
}

// deleteHomework permanently removes a homework by ID. Returns EUNAUTHORIZED if current
// homework is not the one being deleted.
func deleteHomework(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	currentUserID := api.UserIDFromContext(ctx)
	if hw, err := findHomeworkByID(ctx, tx, id); err != nil {
		return err
	} else if currentUserID != 0 && hw.TeacherID != 0 && hw.TeacherID != currentUserID {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this homework.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM homeworks WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}

// attachHomeworkAssociations attaches group and owner objects associated with the homework.
func attachHomeworkAssociations(ctx context.Context, tx *Tx, hw *api.Homework) (err error) {
	if hw.TeacherID == 0 {
		return nil
	} else if hw.Teacher, err = findUserByID(ctx, tx, hw.TeacherID); err != nil {
		return fmt.Errorf("attach homework user: %w", err)
	} else if hw.GroupID == 0 {
		return nil
	} else if hw.Group, err = findGroupByID(ctx, tx, hw.GroupID); err != nil {
		return fmt.Errorf("attach homework group: %w", err)
	}
	return nil
}
