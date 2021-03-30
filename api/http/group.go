package http

import (
	"encoding/json"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dori7879/senior-project/api"
	"github.com/gorilla/mux"
)

// registerGroupRoutes is a helper function for registering all group routes.
func (s *Server) registerGroupRoutes(r *mux.Router) {
	// API endpoint for creating groups.
	r.HandleFunc("/groups", s.handleGroupCreate).Methods("POST")

	r.HandleFunc("/groups/{id}", s.handleGroupView).Methods("GET")
	r.HandleFunc("/groups/{id}", s.handleGroupUpdate).Methods("PATCH")

	// Removing a group.
	r.HandleFunc("/groups/{id}", s.handleGroupDelete).Methods("DELETE")

	// Updating the value for the user's members.
	r.HandleFunc("/groups/{id}/members", s.handleAddMembers).Methods("POST")

	// Accept a share of the group as a teacher via a link
	r.HandleFunc("/groups/{link}/accept", s.handleAcceptGroupShare).Methods("POST")

	// Remove member from a group
	r.HandleFunc("/groups/{groupID}/members/{userID}", s.handleRemoveMember).Methods("DELETE")
}

// handleGroupView handles the "GET /groups/:id" route. It updates
func (s *Server) handleGroupView(w http.ResponseWriter, r *http.Request) {
	// Parse ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	var sFilter api.MemberFilter
	sFilter.Offset, _ = strconv.Atoi(r.URL.Query().Get("student_offset"))
	sFilter.Limit, _ = strconv.Atoi(r.URL.Query().Get("student_limit"))
	if sFilter.Limit == 0 {
		sFilter.Limit = 10
	}

	var tFilter api.MemberFilter
	tFilter.Offset, _ = strconv.Atoi(r.URL.Query().Get("teacher_offset"))
	tFilter.Limit, _ = strconv.Atoi(r.URL.Query().Get("teacher_limit"))
	if tFilter.Limit == 0 {
		tFilter.Limit = 10
	}

	// Fetch group from the database.
	group, err := s.GroupService.FindGroupByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	group.Owner, err = s.UserService.FindUserByID(r.Context(), group.OwnerID)
	if err != nil {
		Error(w, r, err)
		return
	}

	var teacher bool = true
	var student bool = false

	// Fetch associated memberships from the database.
	sFilter.IsTeacher = &student
	sFilter.GroupID = &group.ID
	group.Members.Users, group.Members.N, err = s.UserService.FindMembersByGroup(r.Context(), sFilter)
	if err != nil {
		Error(w, r, err)
		return
	}

	user := api.UserFromContext(r.Context())
	if user.IsTeacher {
		tFilter.IsTeacher = &teacher
		tFilter.GroupID = &group.ID
		group.Teachers.Users, group.Teachers.N, err = s.UserService.FindMembersByGroup(r.Context(), tFilter)
		if err != nil {
			Error(w, r, err)
			return
		}
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		LogError(r, err)
		return
	}
}

// handleGroupCreate handles the "POST /groups" route.
// It reads & writes data using with HTML or JSON.
func (s *Server) handleGroupCreate(w http.ResponseWriter, r *http.Request) {
	// Unmarshal data
	var group api.Group
	if err := json.NewDecoder(r.Body).Decode(&group); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Generate share link
	rand.Seed(time.Now().UnixNano())
	group.ShareLink = api.RandStringSeq(11)

	// Assign owner
	group.OwnerID = api.UserIDFromContext(r.Context())

	// Create group in the database.
	err := s.GroupService.CreateGroup(r.Context(), &group)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Write new group content to response based on accept header.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(group); err != nil {
		LogError(r, err)
		return
	}
}

// handleGroupUpdate handles the "PATCH /groups/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleGroupUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse group ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	var upd api.GroupUpdate
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the group in the database.
	_, err = s.GroupService.UpdateGroup(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleGroupDelete handles the "DELETE /groups/:id" route. This route
// permanently deletes the group and all its members.
func (s *Server) handleGroupDelete(w http.ResponseWriter, r *http.Request) {
	// Parse group ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the group from the database.
	if err := s.GroupService.DeleteGroup(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleAddMembers handles the "PUT /groups/:id/members" route.
func (s *Server) handleAddMembers(w http.ResponseWriter, r *http.Request) {
	members := &struct {
		Users []int `json:"Users"`
	}{}
	if err := json.NewDecoder(r.Body).Decode(members); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Parse group ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Add members to the group.
	if err := s.GroupService.AddStudents(r.Context(), id, members.Users); err != nil {
		Error(w, r, err)
		return
	}

	// Write response to indicate success.
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

func (s *Server) handleAcceptGroupShare(w http.ResponseWriter, r *http.Request) {
	// Parse group share link from path.
	shareLink := mux.Vars(r)["link"]

	user := api.UserFromContext(r.Context())
	if !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Fetch group from the database.
	group, err := s.GroupService.FindGroupByShareLink(r.Context(), shareLink)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Add the teacher as a member of the group
	err = s.GroupService.AddTeacher(r.Context(), group.ID, user.ID)
	if err != nil {
		Error(w, r, err)
		return
	}

	// First fetch owner, then fetch other relations for the response body
	group.Owner, err = s.UserService.FindUserByID(r.Context(), group.OwnerID)
	if err != nil {
		Error(w, r, err)
		return
	}

	var teacher bool = true
	var student bool = false

	// Fetch associated memberships from the database.
	group.Members.Users, group.Members.N, err = s.UserService.FindMembersByGroup(r.Context(), api.MemberFilter{IsTeacher: &student, GroupID: &group.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	group.Teachers.Users, group.Teachers.N, err = s.UserService.FindMembersByGroup(r.Context(), api.MemberFilter{IsTeacher: &teacher, GroupID: &group.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		LogError(r, err)
		return
	}
}

// handleRemoveMember handles the "DELETE /groups/:id/members/:id" route.
func (s *Server) handleRemoveMember(w http.ResponseWriter, r *http.Request) {
	// Parse user ID from path.
	userID, err := strconv.Atoi(mux.Vars(r)["userID"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse group ID from path.
	groupID, err := strconv.Atoi(mux.Vars(r)["groupID"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	user, err := s.UserService.FindUserByID(r.Context(), userID)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Delete the member from the group from the database.
	if err := s.GroupService.RemoveMember(r.Context(), groupID, userID, user.IsTeacher); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
