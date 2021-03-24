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
var _ api.QuizService = (*QuizService)(nil)

// QuizService represents a service for managing quizzes.
type QuizService struct {
	db *DB
}

// NewQuizService returns a new instance of QuizService.
func NewQuizService(db *DB) *QuizService {
	return &QuizService{db: db}
}

// FindQuizByID retrieves a quiz by ID along with their associated group and owner objects.
// Returns ENOTFOUND if quiz does not exist.
func (s *QuizService) FindQuizByID(ctx context.Context, id int) (*api.Quiz, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch quiz and their associated group and owner objects.
	qz, err := findQuizByID(ctx, tx, id)
	if err != nil {
		return nil, err
	} else if err := attachQuizAssociations(ctx, tx, qz); err != nil {
		return qz, err
	}
	return qz, nil
}

// FindQuizzes retrieves a list of quizzes by filter. Also returns total count of
// matching quizzes which may differ from returned results if filter.Limit is specified.
func (s *QuizService) FindQuizzes(ctx context.Context, filter api.QuizFilter) ([]*api.Quiz, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findQuizzes(ctx, tx, filter)
}

// CreateQuiz creates a new quiz.
func (s *QuizService) CreateQuiz(ctx context.Context, qz *api.Quiz) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new quiz object and attach associated group and owner objects.
	if err := createQuiz(ctx, tx, qz); err != nil {
		return err
	} else if err := attachQuizAssociations(ctx, tx, qz); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateQuiz updates a quiz object. Returns EUNAUTHORIZED if current quiz is
// not the quiz that is being updated. Returns ENOTFOUND if quiz does not exist.
func (s *QuizService) UpdateQuiz(ctx context.Context, id int, upd api.QuizUpdate) (*api.Quiz, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update quiz & attach associated group and owner objects.
	qz, err := updateQuiz(ctx, tx, id, upd)
	if err != nil {
		return qz, err
	} else if err := attachQuizAssociations(ctx, tx, qz); err != nil {
		return qz, err
	} else if err := tx.Commit(); err != nil {
		return qz, err
	}
	return qz, nil
}

// DeleteQuiz permanently deletes a quiz.
// Returns EUNAUTHORIZED if current quiz is not the quiz being deleted.
// Returns ENOTFOUND if quiz does not exist.
func (s *QuizService) DeleteQuiz(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteQuiz(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findQuizByID is a helper function to fetch a quiz by ID.
// Returns ENOTFOUND if quiz does not exist.
func findQuizByID(ctx context.Context, tx *Tx, id int) (*api.Quiz, error) {
	a, _, err := findQuizzes(ctx, tx, api.QuizFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Quiz not found."}
	}
	return a[0], nil
}

// findQuizByStudentLink is a helper function to fetch a quiz by studentlink.
// Returns ENOTFOUND if quiz does not exist.
func findQuizByStudentLink(ctx context.Context, tx *Tx, studentlink string) (*api.Quiz, error) {
	a, _, err := findQuizzes(ctx, tx, api.QuizFilter{StudentLink: &studentlink})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Quiz not found."}
	}
	return a[0], nil
}

// findQuizByTeacherLink is a helper function to fetch a quiz by teacherlink.
// Returns ENOTFOUND if quiz does not exist.
func findQuizByTeacherLink(ctx context.Context, tx *Tx, teacherlink string) (*api.Quiz, error) {
	a, _, err := findQuizzes(ctx, tx, api.QuizFilter{TeacherLink: &teacherlink})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Quiz not found."}
	}
	return a[0], nil
}

// findQuizzes returns a list of quizzes matching a filter. Also returns a count of
// total matching quizzes which may differ if filter.Limit is set.
func findQuizzes(ctx context.Context, tx *Tx, filter api.QuizFilter) (_ []*api.Quiz, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	if v := filter.Title; v != nil {
		where, args = append(where, "title = $2"), append(args, *v)
	}
	if v := filter.StudentLink; v != nil {
		where, args = append(where, "student_link = $3"), append(args, *v)
	}
	if v := filter.TeacherLink; v != nil {
		where, args = append(where, "teacher_link = $4"), append(args, *v)
	}
	if v := filter.CourseTitle; v != nil {
		where, args = append(where, "course_title = $5"), append(args, *v)
	}
	if v := filter.Mode; v != nil {
		where, args = append(where, "mode = $6"), append(args, *v)
	}
	if v := filter.TeacherFullName; v != nil {
		where, args = append(where, "teacher_fullname = $7"), append(args, *v)
	}
	if v := filter.TeacherID; v != nil {
		where, args = append(where, "teacher_id = $8"), append(args, *v)
	}
	if v := filter.GroupID; v != nil {
		where, args = append(where, "group_id = $9"), append(args, *v)
	}

	// Execute query to fetch quiz rows.
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
		FROM quizzes
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Quiz objects.
	quizzes := make([]*api.Quiz, 0)
	for rows.Next() {
		var content sql.NullString
		var teacherFullname sql.NullString
		var updatedAt sql.NullTime
		var openedAt sql.NullTime
		var closedAt sql.NullTime
		var teacherID sql.NullInt32
		var groupID sql.NullInt32

		var qz api.Quiz
		if rows.Scan(
			&qz.ID,
			&qz.Title,
			&content,
			&qz.MaxGrade,
			&qz.StudentLink,
			&qz.TeacherLink,
			&qz.CourseTitle,
			&qz.Mode,
			&qz.CreatedAt,
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
			qz.Content = content.String
		}
		if teacherFullname.Valid {
			qz.TeacherFullName = teacherFullname.String
		}
		if updatedAt.Valid {
			qz.UpdatedAt = updatedAt.Time
		}
		if openedAt.Valid {
			qz.OpenedAt = openedAt.Time
		}
		if closedAt.Valid {
			qz.ClosedAt = closedAt.Time
		}
		if teacherID.Valid {
			qz.TeacherID = int(teacherID.Int32)
		}
		if groupID.Valid {
			qz.GroupID = int(groupID.Int32)
		}

		quizzes = append(quizzes, &qz)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return quizzes, n, nil
}

// createQuiz creates a new quiz. Sets the new database ID to qz.ID and sets
// the timestamps to the current time.
func createQuiz(ctx context.Context, tx *Tx, qz *api.Quiz) error {
	// Set timestamps to the current time.
	qz.CreatedAt = tx.now
	qz.UpdatedAt = qz.CreatedAt

	// Perform basic field validation.
	if err := qz.Validate(); err != nil {
		return err
	}

	// Content is nullable so ensure we store blank fields as NULLs.
	var content *string
	if qz.Content != "" {
		content = &qz.Content
	}
	var openedAt *time.Time
	if !qz.OpenedAt.IsZero() {
		openedAt = &qz.OpenedAt
	}
	var closedAt *time.Time
	if !qz.ClosedAt.IsZero() {
		closedAt = &qz.ClosedAt
	}
	var teacherFullname *string
	if qz.TeacherFullName != "" {
		teacherFullname = &qz.TeacherFullName
	}
	var teacherID *int
	if qz.TeacherID != 0 {
		teacherID = &qz.TeacherID
	}
	var groupID *int
	if qz.GroupID != 0 {
		groupID = &qz.GroupID
	}

	// Execute insertion query.
	result, err := tx.ExecContext(ctx, `
		INSERT INTO quizzes (
			title,
			content,
			max_grade,
			course_title,
			mode,
			opened_at,
			closed_at,
			teacher_fullname,
			teacher_id,
			group_id
		)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`,
		qz.Title,
		content,
		qz.MaxGrade,
		qz.CourseTitle,
		qz.Mode,
		openedAt,
		closedAt,
		teacherFullname,
		teacherID,
		groupID,
	)
	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	qz.ID = int(id)

	return nil
}

// updateQuiz updates fields on a quiz object. Returns EUNAUTHORIZED if current
// quiz is not the quiz being updated.
func updateQuiz(ctx context.Context, tx *Tx, id int, upd api.QuizUpdate) (*api.Quiz, error) {
	// Fetch current object state.
	qz, err := findQuizByID(ctx, tx, id)
	if err != nil {
		return qz, err
	} else if qz.ID != api.UserIDFromContext(ctx) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this quiz.")
	}

	// Update fields.
	if v := upd.Title; v != nil {
		qz.Title = *v
	}
	if v := upd.Content; v != nil {
		qz.Content = *v
	}
	if v := upd.MaxGrade; v != nil {
		qz.MaxGrade = *v
	}
	if v := upd.CourseTitle; v != nil {
		qz.CourseTitle = *v
	}
	if v := upd.Mode; v != nil {
		qz.Mode = *v
	}
	if v := upd.OpenedAt; v != nil {
		qz.OpenedAt = *v
	}
	if v := upd.ClosedAt; v != nil {
		qz.ClosedAt = *v
	}
	if v := upd.TeacherFullName; v != nil {
		qz.TeacherFullName = *v
	}
	if v := upd.TeacherID; v != nil {
		qz.TeacherID = *v
	}
	if v := upd.GroupID; v != nil {
		qz.GroupID = *v
	}

	// Set last updated date to current time.
	qz.UpdatedAt = tx.now

	// Perform basic field validation.
	if err := qz.Validate(); err != nil {
		return qz, err
	}

	// These fields are nullable so ensure we store blank fields as NULLs.
	var content *string
	if qz.Content != "" {
		content = &qz.Content
	}
	var openedAt *time.Time
	if !qz.OpenedAt.IsZero() {
		openedAt = &qz.OpenedAt
	}
	var closedAt *time.Time
	if !qz.ClosedAt.IsZero() {
		closedAt = &qz.ClosedAt
	}
	var teacherFullname *string
	if qz.TeacherFullName != "" {
		teacherFullname = &qz.TeacherFullName
	}
	var teacherID *int
	if qz.TeacherID != 0 {
		teacherID = &qz.TeacherID
	}
	var groupID *int
	if qz.GroupID != 0 {
		groupID = &qz.GroupID
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE quizzes
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
		qz.Title,
		content,
		qz.MaxGrade,
		qz.CourseTitle,
		qz.Mode,
		openedAt,
		closedAt,
		teacherFullname,
		teacherID,
		groupID,
		id,
	); err != nil {
		return qz, FormatError(err)
	}

	return qz, nil
}

// deleteQuiz permanently removes a quiz by ID. Returns EUNAUTHORIZED if current
// quiz is not the one being deleted.
func deleteQuiz(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	if qz, err := findQuizByID(ctx, tx, id); err != nil {
		return err
	} else if qz.ID != api.UserIDFromContext(ctx) {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this quiz.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM quizzes WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}

// attachQuizAssociations attaches group and owner objects associated with the quiz.
func attachQuizAssociations(ctx context.Context, tx *Tx, qz *api.Quiz) (err error) {
	if qz.Teacher, err = findUserByID(ctx, tx, qz.TeacherID); err != nil {
		return fmt.Errorf("attach quiz user: %w", err)
	} else if qz.Group, err = findGroupByID(ctx, tx, qz.GroupID); err != nil {
		return fmt.Errorf("attach quiz group: %w", err)
	}
	return nil
}
