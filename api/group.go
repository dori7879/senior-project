package api

import (
	"context"
)

// Group represents a group in the system.
type Group struct {
	ID int `json:"ID"`

	Title     string  `json:"Title"`
	ShareLink string  `json:"ShareLink"`
	OwnerID   int     `json:"OwnerID"`
	Owner     *User   `json:"Owner"`
	Teachers  []*User `json:"Teachers"`
	Members   []*User `json:"Members"`
}

// Validate returns an error if the group contains invalid fields.
// This only performs basic validation.
func (u *Group) Validate() error {
	if u.Title == "" {
		return Errorf(EINVALID, "Title required.")
	}
	return nil
}

// GroupService represents a service for managing groups.
type GroupService interface {
	// Retrieves a group by ID.
	// Returns ENOTFOUND if group does not exist.
	FindGroupByID(ctx context.Context, id int) (*Group, error)

	// Retrieves a list of groups by filter. Also returns total count of matching
	// groups which may differ from returned results if filter.Limit is specified.
	FindGroups(ctx context.Context, filter GroupFilter) ([]*Group, int, error)

	// Creates a new group.
	CreateGroup(ctx context.Context, group *Group) error

	// Updates a group object. Returns EUNAUTHORIZED if current group is not
	// the group that is being updated. Returns ENOTFOUND if group does not exist.
	UpdateGroup(ctx context.Context, id int, upd GroupUpdate) (*Group, error)

	// Permanently deletes a group and all owned dials. Returns EUNAUTHORIZED
	// if current group is not the group being deleted. Returns ENOTFOUND if
	// group does not exist.
	DeleteGroup(ctx context.Context, id int) error
}

// GroupFilter represents a filter passed to FindGroups().
type GroupFilter struct {
	// Filtering fields.
	ID        *int    `json:"ID"`
	Title     *string `json:"Title"`
	ShareLink *string `json:"ShareLink"`
	OwnerID   *int    `json:"OwnerID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// GroupUpdate represents a set of fields to be updated via UpdateGroup().
type GroupUpdate struct {
	Title *string `json:"Title"`
}
