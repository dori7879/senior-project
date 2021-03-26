package pg

import (
	"context"
	"strings"

	"github.com/dori7879/senior-project/api"
)

// Ensure service implements interface.
var _ api.GroupService = (*GroupService)(nil)

// GroupService represents a service for managing groups.
type GroupService struct {
	db *DB
}

// NewGroupService returns a new instance of GroupService.
func NewGroupService(db *DB) *GroupService {
	return &GroupService{db: db}
}

// FindGroupByID retrieves a group by ID.
// Returns ENOTFOUND if group does not exist.
func (s *GroupService) FindGroupByID(ctx context.Context, id int) (*api.Group, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch group.
	group, err := findGroupByID(ctx, tx, id)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// FindGroupByShareLink retrieves a group by share link.
// Returns ENOTFOUND if group does not exist.
func (s *GroupService) FindGroupByShareLink(ctx context.Context, link string) (*api.Group, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Fetch group.
	group, err := findGroupByShareLink(ctx, tx, link)
	if err != nil {
		return nil, err
	}
	return group, nil
}

// FindGroups retrieves a list of groups by filter. Also returns total count of
// matching groups which may differ from returned results if filter.Limit is specified.
func (s *GroupService) FindGroups(ctx context.Context, filter api.GroupFilter) ([]*api.Group, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findGroups(ctx, tx, filter)
}

// FindGroupsByMember retrieves a list of groups for either a teacher who has been shared with
// or a student who is a member of.
func (s *GroupService) FindGroupsByMember(ctx context.Context, filter api.UserFilter) ([]*api.Group, int, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, 0, err
	}
	defer tx.Rollback()
	return findGroupsByMember(ctx, tx, filter)
}

// CreateGroup creates a new group.
func (s *GroupService) CreateGroup(ctx context.Context, group *api.Group) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	// Create a new group object.
	if err := createGroup(ctx, tx, group); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateGroup updates a group object. Returns EUNAUTHORIZED if current group is
// not the group that is being updated. Returns ENOTFOUND if group does not exist.
func (s *GroupService) UpdateGroup(ctx context.Context, id int, upd api.GroupUpdate) (*api.Group, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()

	// Update group.
	group, err := updateGroup(ctx, tx, id, upd)
	if err != nil {
		return group, err
	} else if err := tx.Commit(); err != nil {
		return group, err
	}
	return group, nil
}

// DeleteGroup permanently deletes a group.
// Returns EUNAUTHORIZED if current group is not the group being deleted.
// Returns ENOTFOUND if group does not exist.
func (s *GroupService) DeleteGroup(ctx context.Context, id int) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if err := deleteGroup(ctx, tx, id); err != nil {
		return err
	}
	return tx.Commit()
}

// findGroupByID is a helper function to fetch a group by ID.
// Returns ENOTFOUND if group does not exist.
func findGroupByID(ctx context.Context, tx *Tx, id int) (*api.Group, error) {
	a, _, err := findGroups(ctx, tx, api.GroupFilter{ID: &id})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Group not found."}
	}
	return a[0], nil
}

// findGroupByShareLink is a helper function to fetch a group by link.
// Returns ENOTFOUND if group does not exist.
func findGroupByShareLink(ctx context.Context, tx *Tx, link string) (*api.Group, error) {
	a, _, err := findGroups(ctx, tx, api.GroupFilter{ShareLink: &link})
	if err != nil {
		return nil, err
	} else if len(a) == 0 {
		return nil, &api.Error{Code: api.ENOTFOUND, Message: "Group not found."}
	}
	return a[0], nil
}

// findGroups returns a list of groups matching a filter. Also returns a count of
// total matching groups which may differ if filter.Limit is set.
func findGroups(ctx context.Context, tx *Tx, filter api.GroupFilter) (_ []*api.Group, n int, err error) {
	// Build WHERE clause.
	where, args := []string{"1 = 1"}, []interface{}{}
	if v := filter.ID; v != nil {
		where, args = append(where, "id = $1"), append(args, *v)
	}
	if v := filter.Title; v != nil {
		where, args = append(where, "title = $2"), append(args, *v)
	}
	if v := filter.ShareLink; v != nil {
		where, args = append(where, "share_link = $3"), append(args, *v)
	}
	if v := filter.OwnerID; v != nil {
		where, args = append(where, "owner_id = $4"), append(args, *v)
	}

	// Execute query to fetch group rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    id,
		    title,
		    share_link,
		    owner_id,
		    COUNT(*) OVER()
		FROM groups
		WHERE `+strings.Join(where, " AND ")+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		args...,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Group objects.
	groups := make([]*api.Group, 0)
	for rows.Next() {
		var group api.Group

		if rows.Scan(
			&group.ID,
			&group.Title,
			&group.ShareLink,
			&group.OwnerID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		groups = append(groups, &group)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return groups, n, nil
}

// findGroupsByMember returns a list of groups by teacher ID or student ID.
func findGroupsByMember(ctx context.Context, tx *Tx, filter api.UserFilter) (_ []*api.Group, n int, err error) {
	var m2m string
	if *filter.IsTeacher {
		m2m = `LEFT JOIN teachers_groups t on t.group_id = g.id
		WHERE t.teacher_id = ?`
	} else {
		m2m = `LEFT JOIN students_groups s on s.group_id = g.id
		WHERE s.student_id = ?`
	}

	// Execute query to fetch group rows.
	rows, err := tx.QueryContext(ctx, `
		SELECT 
		    g.id,
		    g.title,
		    g.share_link,
		    g.owner_id,
		    COUNT(*) OVER()
		FROM groups g
		`+m2m+`
		ORDER BY id ASC
		`+FormatLimitOffset(filter.Limit, filter.Offset),
		*filter.ID,
	)
	if err != nil {
		return nil, n, err
	}
	defer rows.Close()

	// Deserialize rows into Group objects.
	groups := make([]*api.Group, 0)
	for rows.Next() {
		var group api.Group

		if rows.Scan(
			&group.ID,
			&group.Title,
			&group.ShareLink,
			&group.OwnerID,
			&n,
		); err != nil {
			return nil, 0, err
		}

		groups = append(groups, &group)
	}
	if err := rows.Err(); err != nil {
		return nil, 0, err
	}

	return groups, n, nil
}

// createGroup creates a new group. Sets the new database ID to group.ID and sets
// the timestamps to the current time.
func createGroup(ctx context.Context, tx *Tx, group *api.Group) error {
	// Perform basic field validation.
	if err := group.Validate(); err != nil {
		return err
	}

	// Execute insertion query.
	result, err := tx.ExecContext(ctx, `
		INSERT INTO groups (
			title,
			share_link,
			owner_id
		)
		VALUES ($1, $2, $3)
	`,
		group.Title,
		group.ShareLink,
		group.OwnerID,
	)
	if err != nil {
		return FormatError(err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	group.ID = int(id)

	return nil
}

// updateGroup updates fields on a group object. Returns EUNAUTHORIZED if current
// group is not the group being updated.
func updateGroup(ctx context.Context, tx *Tx, id int, upd api.GroupUpdate) (*api.Group, error) {
	// Fetch current object state.
	group, err := findGroupByID(ctx, tx, id)
	if err != nil {
		return group, err
	} else if group.ID != api.UserIDFromContext(ctx) {
		return nil, api.Errorf(api.EUNAUTHORIZED, "You are not allowed to update this group.")
	}

	// Update fields.
	if v := upd.Title; v != nil {
		group.Title = *v
	}

	// Perform basic field validation.
	if err := group.Validate(); err != nil {
		return group, err
	}

	// Execute update query.
	if _, err := tx.ExecContext(ctx, `
		UPDATE groups
		SET title = $1
		WHERE id = $2
	`,
		group.Title,
		id,
	); err != nil {
		return group, FormatError(err)
	}

	return group, nil
}

// deleteGroup permanently removes a group by ID. Returns EUNAUTHORIZED if current
// group is not the one being deleted.
func deleteGroup(ctx context.Context, tx *Tx, id int) error {
	// Verify object exists.
	if group, err := findGroupByID(ctx, tx, id); err != nil {
		return err
	} else if group.ID != api.UserIDFromContext(ctx) {
		return api.Errorf(api.EUNAUTHORIZED, "You are not allowed to delete this group.")
	}

	// Remove row from database.
	if _, err := tx.ExecContext(ctx, `DELETE FROM groups WHERE id = $1`, id); err != nil {
		return FormatError(err)
	}
	return nil
}
