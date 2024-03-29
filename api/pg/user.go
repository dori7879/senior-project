package pg

import (
	"context"
	"fmt"
	"strings"

	"github.com/dori7879/senior-project/api"
)

// Ensure service implements interface.
var _ api.UserService = (*UserService)(nil)

// UserService represents a service for managing users.
type UserService struct {
	db *DB
}

// NewUserService returns a new instance of UserService.
func NewUserService(db *DB) *UserService {
	return &UserService{db: db}
}

// FindUserByID retrieves a user by ID.
// Returns ENOTFOUND if user does not exist.
func (s *UserService) FindUserByID(ctx context.Context, id int) (*api.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch user.
	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUserByEmail retrieves a user by Email.
// Returns ENOTFOUND if user does not exist.
func (s *UserService) FindUserByEmail(ctx context.Context, email string) (*api.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch user.
	user, err := findUserByEmail(ctx, tx, email)
	if err != nil {
		return nil, err
	}
	return user, nil
}

// FindUsers retrieves a list of users by filter. Also returns total count of
// matching users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindUsers(ctx context.Context, filter api.UserFilter) ([]*api.User, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findUsers(ctx, tx, filter)
}

// Retrieves a list of users by group. Also returns total count of matching
// users which may differ from returned results if filter.Limit is specified.
func (s *UserService) FindMembersByGroup(ctx context.Context, filter api.MemberFilter) ([]*api.User, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findUsersByGroup(ctx, tx, filter)
}

// CreateUser creates a new user.
func (s *UserService) CreateUser(ctx context.Context, user *api.User) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new user object.
	if err := createUser(ctx, tx, user); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateUser updates a user object. Returns EUNAUTHORIZED if current user is
// not the user that is being updated. Returns ENOTFOUND if user does not exist.
func (s *UserService) UpdateUser(ctx context.Context, id int, upd api.UserUpdate) (*api.User, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update user.
	user, err := updateUser(ctx, tx, id, upd)
	if err != nil {
		return user, err
	} else if err := tx.Commit(); err != nil {
		return user, err
	}
	return user, nil
}

// DeleteUser permanently deletes a user.
// Returns EUNAUTHORIZED if current user is not the user being deleted.
// Returns ENOTFOUND if user does not exist.
func (s *UserService) DeleteUser(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteUser(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findUserByID is a helper function to fetch a user by ID.
// Returns ENOTFOUND if user does not exist.
func findUserByID(ctx context.Context, tx *Tx, id int) (*api.User, error) {
	a, _, err := findUsers(ctx, tx, api.UserFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}

// findUserByEmail is a helper function to fetch a user by email.
// Returns ENOTFOUND if user does not exist.
func findUserByEmail(ctx context.Context, tx *Tx, email string) (*api.User, error) {
	a, _, err := findUsers(ctx, tx, api.UserFilter{Email: &email})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "User not found."}
	}
	return a[0], nil
}

// findUsers returns a list of users matching a filter. Also returns a count of
// total matching users which may differ if filter.Limit is set.
func findUsers(ctx context.Context, tx *Tx, filter api.UserFilter) (_ []*api.User, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	i := 1
	if v := filter.ID; v != nil {
		where, args = append(where, fmt.Sprintf("id = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.Email; v != nil {
		where, args = append(where, fmt.Sprintf("email = $%d", i)), append(args, *v)
		i++
	}
	if v := filter.IsTeacher; v != nil {
		where, args = append(where, fmt.Sprintf("is_teacher = $%d", i)), append(args, *v)
		i++
	}

	// Add Email substring filtering w/ built in argument
	if v := filter.EmailSubStr; v != nil {
		where = append(where, fmt.Sprintf("email LIKE '%%%s%%'", *v))
	}

	// Execute query to fetch user rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
			id,
			first_name,
			last_name,
			email,
			password_hash,
			date_joined,
			is_teacher,
			COUNT(*) OVER()
		FROM users
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into User objects.
	users := make([]*api.User, 0)
	for rows.Next() {
		var user api.User

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PasswordHash,
			&user.DateJoined,
			&user.IsTeacher,
			&n,
		); err != nil {
			return nil, 0, err
		}

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, n, nil
}

// findUsersByGroup returns a list of users by group filter. Also returns a count of
// total matching users which may differ if filter.Limit is set.
func findUsersByGroup(ctx context.Context, tx *Tx, filter api.MemberFilter) (_ []*api.User, n int, err error) {
	var m2m string
	if *filter.IsTeacher {
		m2m = `FROM teachers_groups t
		LEFT JOIN users u on u.id = t.teacher_id
		WHERE t.group_id = $1`
	} else {
		m2m = `FROM students_groups s
		LEFT JOIN users u on u.id = s.student_id
		WHERE s.group_id = $1`
	}

	args := []interface{}{}
	args = append(args, *filter.GroupID)

	if v := filter.UserID; v != nil {
		if *filter.IsTeacher {
			m2m += " and t.teacher_id = $2"
		} else {
			m2m += " and s.student_id = $2"
		}
		args = append(args, *v)
	}

	// Execute query to fetch user rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    u.id,
		    u.first_name,
		    u.last_name,
		    u.email,
		    u.password_hash,
		    u.date_joined,
			u.is_teacher,
		    COUNT(*) OVER()
		`+m2m+`
		ORDER BY u.id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into User objects.
	users := make([]*api.User, 0)
	for rows.Next() {
		var user api.User

		if err := rows.Scan(
			&user.ID,
			&user.FirstName,
			&user.LastName,
			&user.Email,
			&user.PasswordHash,
			&user.DateJoined,
			&user.IsTeacher,
			&n,
		); err != nil {
			return nil, 0, err
		}

		users = append(users, &user)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return users, n, nil
}

// createUser creates a new user. Sets the new database ID to user.ID and sets
// the timestamps to the current time.
func createUser(ctx context.Context, tx *Tx, user *api.User) error {
	// Set timestamps to the current time.
	user.DateJoined = tx.now

	// Perform basic field validation.
	if err := user.Validate(); err != nil {
		return err
	}

	// Execute insertion query.
	row := tx.QueryRowContext(ctx, `
		INSERT INTO users (
			first_name,
			last_name,
			email,
			password_hash,
			date_joined,
			is_teacher
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.PasswordHash,
		user.DateJoined,
		user.IsTeacher,
	)

	err := row.Scan(&user.ID)
	if err != nil {
		return FormatError(err)
	}

	return nil
}

// updateUser updates fields on a user object. Returns EUNAUTHORIZED if current
// user is not the user being updated.
func updateUser(ctx context.Context, tx *Tx, id int, upd api.UserUpdate) (*api.User, error) {
	// Fetch current object state.
	user, err := findUserByID(ctx, tx, id)
	if err != nil {
		return user, err
	} else if user.ID != api.UserIDFromContext(ctx) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this user.")
	}

	// Update fields.
	if v := upd.FirstName; v != nil {
		user.FirstName = *v
	}
	if v := upd.LastName; v != nil {
		user.LastName = *v
	}
	if v := upd.Email; v != nil {
		user.Email = *v
	}
	if v := upd.IsTeacher; v != nil {
		user.IsTeacher = *v
	}
	if v := upd.PasswordHash; v != nil {
		user.PasswordHash = *v
	}

	// Perform basic field validation.
	if err := user.Validate(); err != nil {
		return user, err
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE users
		SET first_name = $1,
			last_name = $2,
		    email = $3,
		    is_teacher = $4,
			password_hash = $5
		WHERE id = $6
	`,
		user.FirstName,
		user.LastName,
		user.Email,
		user.IsTeacher,
		user.PasswordHash,
		id,
	); err != nil {
		return user, FormatError(err)
	}

	return user, nil
}

// deleteUser permanently removes a user by ID. Returns EUNAUTHORIZED if current
// user is not the one being deleted.
func deleteUser(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	if user, err := findUserByID(ctx, tx, id); err != nil {
		return err
	} else if user.ID != api.UserIDFromContext(ctx) {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this user.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM users WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}
