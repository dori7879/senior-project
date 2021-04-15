package api

import (
	"context"
	"time"
)

// Attendance represents a attendance in the system.
type Attendance struct {
	ID int `json:"ID"`

	Title       string    `json:"Title"`
	PIN         string    `json:"PIN"`
	StudentLink string    `json:"StudentLink"`
	TeacherLink string    `json:"TeacherLink"`
	CourseTitle string    `json:"CourseTitle"`
	Mode        string    `json:"Mode"`
	CreatedAt   time.Time `json:"CreatedAt"`
	UpdatedAt   time.Time `json:"UpdatedAt"`
	OpenedAt    time.Time `json:"OpenedAt"`
	ClosedAt    time.Time `json:"ClosedAt"`

	TeacherFullName string `json:"TeacherFullName"`
	TeacherID       int    `json:"TeacherID"`
	Teacher         *User  `json:"Teacher"`

	GroupID int    `json:"GroupID"`
	Group   *Group `json:"Group"`

	Submissions []*AttSubmission `json:"Submissions,omitempty"`
}

// Validate returns an error if the attendance contains invalid fields.
// This only performs basic validation.
func (u *Attendance) Validate() error {
	if u.Title == "" {
		return Errorf(EINVALID, "Title required.")
	} else if u.Mode != All && u.Mode != Registered {
		return Errorf(EINVALID, "Mode is incorrect.")
	}
	return nil
}

// AttendanceService represents a service for managing attendances.
type AttendanceService interface {
	// Retrieves a attendance by ID.
	// Returns ENOTFOUND if attendance does not exist.
	FindAttendanceByID(ctx context.Context, id int) (*Attendance, error)

	FindAttendanceByStudentLink(ctx context.Context, link string) (*Attendance, error)

	FindAttendanceByTeacherLink(ctx context.Context, link string) (*Attendance, error)

	// Retrieves a list of attendances by filter. Also returns total count of matching
	// attendances which may differ from returned results if filter.Limit is specified.
	FindAttendances(ctx context.Context, filter AttendanceFilter) ([]*Attendance, int, error)

	// Creates a new attendance.
	CreateAttendance(ctx context.Context, attendance *Attendance) error

	// Updates a attendance object. Returns EUNAUTHORIZED if current attendance is not
	// the attendance that is being updated. Returns ENOTFOUND if attendance does not exist.
	UpdateAttendance(ctx context.Context, id int, upd AttendanceUpdate) (*Attendance, error)

	// Permanently deletes a attendance and all owned dials. Returns EUNAUTHORIZED
	// if current attendance is not the attendance being deleted. Returns ENOTFOUND if
	// attendance does not exist.
	DeleteAttendance(ctx context.Context, id int) error
}

// AttendanceFilter represents a filter passed to FindAttendances().
type AttendanceFilter struct {
	// Filtering fields.
	ID          *int    `json:"ID"`
	Title       *string `json:"Title"`
	StudentLink *string `json:"StudentLink"`
	TeacherLink *string `json:"TeacherLink"`
	CourseTitle *string `json:"CourseTitle"`
	Mode        *string `json:"Mode"`

	TeacherFullName *string `json:"TeacherFullName"`
	TeacherID       *int    `json:"TeacherID"`
	GroupID         *int    `json:"GroupID"`

	// Restrict to subset of results.
	Offset int `json:"Offset"`
	Limit  int `json:"Limit"`
}

// AttendanceUpdate represents a set of fields to be updated via UpdateAttendance().
type AttendanceUpdate struct {
	Title       *string    `json:"Title"`
	CourseTitle *string    `json:"CourseTitle"`
	Mode        *string    `json:"Mode"`
	OpenedAt    *time.Time `json:"OpenedAt"`
	ClosedAt    *time.Time `json:"ClosedAt"`

	TeacherFullName *string `json:"TeacherFullName"`
	TeacherID       *int    `json:"TeacherID"`
	GroupID         *int    `json:"GroupID"`
}
