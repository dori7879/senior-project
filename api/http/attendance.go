package http

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerAttendancePrivateRoutes is a helper function for registering private attendance routes.
func (s *Server) registerAttendancePrivateRoutes(r *mux.Router) {
	// Listing of all attendances a teacher is an owner of.
	r.HandleFunc("/attendances", s.handleAttendanceList).Methods("GET")

	// View a single attendance.
	r.HandleFunc("/attendances/{id}", s.handleAttendanceView).Methods("GET")
}

// registerAttendancePublicRoutes is a helper function for registering public attendance routes.
func (s *Server) registerAttendancePublicRoutes(r *mux.Router) {
	// API endpoint for creating attendances.
	r.HandleFunc("/attendances", s.handleAttendanceCreate).Methods("POST")

	// View a single attendance.
	r.HandleFunc("/attendances/shared/{link}/teacher", s.handleAttendanceTeacherView).Methods("GET")
	r.HandleFunc("/attendances/shared/{link}/student", s.handleAttendanceStudentView).Methods("GET")

	r.HandleFunc("/attendances/{id}", s.handleAttendanceUpdate).Methods("PATCH")
	r.HandleFunc("/attendances/{id}/renew", s.handleAttendancePINRenew).Methods("PATCH")

	// Removing a attendance.
	r.HandleFunc("/attendances/{id}", s.handleAttendanceDelete).Methods("DELETE")
}

// handleAttendanceList handles the "GET /attendances" route. This route can optionally
// accept filter arguments and outputs a list of all attendances that the current
// user is related to.
func (s *Server) handleAttendanceList(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Fetch attendances from database.
	attendances, n, err := s.AttendanceService.FindAttendances(r.Context(), api.AttendanceFilter{TeacherID: &user.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Render output based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Attendances []*api.Attendance `json:"Attendances"`
		N           int               `json:"N"`
	}{
		Attendances: attendances,
		N:           n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleAttendanceView handles the "GET /attendances/:id" route.
func (s *Server) handleAttendanceView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Parse ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Fetch attendance from the database.
	attendance, err := s.AttendanceService.FindAttendanceByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Fetch associated submissions from the database.
	attendance.Submissions, _, err = s.AttSubmissionService.FindAttSubmissions(r.Context(), api.AttSubmissionFilter{AttendanceID: &attendance.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	if user != nil && attendance.TeacherID != user.ID {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(attendance); err != nil {
		LogError(r, err)
		return
	}
}

// handleAttendanceTeacherView handles the "GET /attendances/shared/:link/teacher" route.
func (s *Server) handleAttendanceTeacherView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Parse teacher link from path.
	link := mux.Vars(r)["link"]

	// Fetch attendance from the database.
	attendance, err := s.AttendanceService.FindAttendanceByTeacherLink(r.Context(), link)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Fetch associated submissions from the database.
	attendance.Submissions, _, err = s.AttSubmissionService.FindAttSubmissions(r.Context(), api.AttSubmissionFilter{AttendanceID: &attendance.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(attendance); err != nil {
		LogError(r, err)
		return
	}
}

// handleAttendanceStudentView handles the "GET /attendances/shared/:link/studnet" route.
func (s *Server) handleAttendanceStudentView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Parse student link from path.
	link := mux.Vars(r)["link"]

	// Fetch attendance from the database.
	attendance, err := s.AttendanceService.FindAttendanceByStudentLink(r.Context(), link)
	if err != nil {
		Error(w, r, err)
		return
	}

	if attendance.GroupID != 0 {
		if user == nil {
			Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not in the group. Also, you are not logged in."))
			return
		}

		isTeacher := false

		// Check whether the user is a member of the group
		gs, _, err := s.GroupService.FindGroupsByMember(r.Context(), api.UserFilter{ID: &user.ID, IsTeacher: &isTeacher})
		if err != nil {
			Error(w, r, err)
			return
		}

		var inGroup bool

		for _, v := range gs {
			if v.ID == attendance.GroupID {
				inGroup = true
			}
		}

		if !inGroup {
			Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not in the group"))
			return
		}
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(attendance); err != nil {
		LogError(r, err)
		return
	}
}

// handleAttendanceCreate handles the "POST /attendances" route.
func (s *Server) handleAttendanceCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Unmarshal data based on HTTP request's content type.
	attendance := api.Attendance{}
	if err := json.NewDecoder(r.Body).Decode(&attendance); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	rand.Seed(time.Now().UnixNano())
	attendance.StudentLink = api.RandStringSeq(11)
	attendance.TeacherLink = api.RandStringSeq(11)
	attendance.PIN = api.RandDigitSeq(6)

	if user != nil {
		attendance.TeacherID = user.ID
	}

	// Create attendance in the database.
	err := s.AttendanceService.CreateAttendance(r.Context(), &attendance)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Write new attendance content to response based on accept header.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		StudentLink string `json:"StudentLink"`
		TeacherLink string `json:"TeacherLink"`
	}{
		StudentLink: attendance.StudentLink,
		TeacherLink: attendance.TeacherLink,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleAttendanceUpdate handles the "PATCH /attendances/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleAttendanceUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse attendance ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	upd := api.AttendanceUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the attendance in the database.
	_, err = s.AttendanceService.UpdateAttendance(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

func (s *Server) handleAttendancePINRenew(w http.ResponseWriter, r *http.Request) {
	// Parse attendance ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	pin := api.RandDigitSeq(6)
	upd := api.AttendanceUpdate{PIN: &pin}

	// Update the attendance in the database.
	_, err = s.AttendanceService.UpdateAttendance(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(fmt.Sprintf(`{"pin": "%s"}`, pin)))
}

// handleAttendanceDelete handles the "DELETE /attendances/:id" route. This route
// permanently deletes the attendance.
func (s *Server) handleAttendanceDelete(w http.ResponseWriter, r *http.Request) {
	// Parse attendance ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the attendance from the database.
	if err := s.AttendanceService.DeleteAttendance(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
