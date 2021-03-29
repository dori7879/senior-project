package api

import (
	"context"
)

// Group represents a group in the system.
type Group struct {
	ID int `json:"ID"`

	Title     string `json:"Title"`
	ShareLink string `json:"ShareLink"`
	OwnerID   int    `json:"OwnerID"`
	Owner     *User  `json:"Owner,omitempty"`

	Teachers struct {
		Users []*User `json:"Users"`
		N     int     `json:"N"`
	} `json:"Teachers"`

	Members struct {
		Users []*User `json:"Users"`
		N     int     `json:"N"`
	} `json:"Members"`
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

	FindGroupByShareLink(ctx context.Context, link string) (*Group, error)

	// Retrieves a list of groups by filter. Also returns total count of matching
	// groups which may differ from returned results if filter.Limit is specified.
	FindGroups(ctx context.Context, filter GroupFilter) ([]*Group, int, error)

	// Retrieves a list of groups for either a teacher who has been shared with
	// or a student who is a member of.
	FindGroupsByMember(ctx context.Context, filter UserFilter) ([]*Group, int, error)

	// Creates a new group.
	CreateGroup(ctx context.Context, group *Group) error

	// Updates a group object. Returns EUNAUTHORIZED if current group is not
	// the group that is being updated. Returns ENOTFOUND if group does not exist.
	UpdateGroup(ctx context.Context, id int, upd GroupUpdate) (*Group, error)

	// Permanently deletes a group and all owned dials. Returns EUNAUTHORIZED
	// if current group is not the group being deleted. Returns ENOTFOUND if
	// group does not exist.
	DeleteGroup(ctx context.Context, id int) error

	AddStudents(ctx context.Context, id int, users []int) error

	AddTeacher(ctx context.Context, groupID int, teacherID int) error

	RemoveMember(ctx context.Context, groupID, userID int, isTeacher bool) error
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

// MemberFilter represents a filter passed to FindUsersByGroup().
type MemberFilter struct {
	// Filtering fields.
	GroupID   *int  `json:"GroupID"`
	IsTeacher *bool `json:"IsTeacher"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// GroupUpdate represents a set of fields to be updated via UpdateGroup().
type GroupUpdate struct {
	Title *string `json:"Title"`
}
