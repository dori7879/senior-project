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

// registerHomeworkPrivateRoutes is a helper function for registering private homework routes.
func (s *Server) registerHomeworkPrivateRoutes(r *mux.Router) {
	// Listing of all homeworks a teacher is an owner of.
	r.HandleFunc("/homeworks", s.handleHomeworkList).Methods("GET")

	// View a single homework.
	r.HandleFunc("/homeworks/{id}", s.handleHomeworkView).Methods("GET")
}

// registerHomeworkPublicRoutes is a helper function for registering public homework routes.
func (s *Server) registerHomeworkPublicRoutes(r *mux.Router) {
	// API endpoint for creating homeworks.
	r.HandleFunc("/homeworks", s.handleHomeworkCreate).Methods("POST")

	// View a single homework.
	r.HandleFunc("/homeworks/shared/{link}/teacher", s.handleHomeworkTeacherView).Methods("GET")
	r.HandleFunc("/homeworks/shared/{link}/student", s.handleHomeworkStudentView).Methods("GET")

	r.HandleFunc("/homeworks/{id}", s.handleHomeworkUpdate).Methods("PATCH")

	// Removing a homework.
	r.HandleFunc("/homeworks/{id}", s.handleHomeworkDelete).Methods("DELETE")
}

// handleHomeworkList handles the "GET /homeworks" route. This route can optionally
// accept filter arguments and outputs a list of all homeworks that the current
// user is related to.
func (s *Server) handleHomeworkList(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Fetch homeworks from database.
	homeworks, n, err := s.HomeworkService.FindHomeworks(r.Context(), api.HomeworkFilter{TeacherID: &user.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Render output based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Homeworks []*api.Homework `json:"Homeworks"`
		N         int             `json:"N"`
	}{
		Homeworks: homeworks,
		N:         n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleHomeworkView handles the "GET /homeworks/:id" route.
func (s *Server) handleHomeworkView(w http.ResponseWriter, r *http.Request) {
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

	// Fetch homework from the database.
	homework, err := s.HomeworkService.FindHomeworkByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Fetch associated submissions from the database.
	homework.Submissions, _, err = s.HWSubmissionService.FindHWSubmissions(r.Context(), api.HWSubmissionFilter{HomeworkID: &homework.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	if user != nil && homework.TeacherID != user.ID {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(homework); err != nil {
		LogError(r, err)
		return
	}
}

// handleHomeworkTeacherView handles the "GET /homeworks/shared/:link/teacher" route.
func (s *Server) handleHomeworkTeacherView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Parse teacher link from path.
	link := mux.Vars(r)["link"]

	// Fetch homework from the database.
	homework, err := s.HomeworkService.FindHomeworkByTeacherLink(r.Context(), link)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Fetch associated submissions from the database.
	homework.Submissions, _, err = s.HWSubmissionService.FindHWSubmissions(r.Context(), api.HWSubmissionFilter{HomeworkID: &homework.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(homework); err != nil {
		LogError(r, err)
		return
	}
}

// handleHomeworkStudentView handles the "GET /homeworks/shared/:link/studnet" route.
func (s *Server) handleHomeworkStudentView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Parse student link from path.
	link := mux.Vars(r)["link"]

	// Fetch homework from the database.
	homework, err := s.HomeworkService.FindHomeworkByStudentLink(r.Context(), link)
	if err != nil {
		Error(w, r, err)
		return
	}

	if homework.GroupID != 0 {
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
			if v.ID == homework.GroupID {
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
	if err := json.NewEncoder(w).Encode(homework); err != nil {
		LogError(r, err)
		return
	}
}

// handleHomeworkCreate handles the "POST /homeworks" route.
func (s *Server) handleHomeworkCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Unmarshal data based on HTTP request's content type.
	homework := api.Homework{}
	if err := json.NewDecoder(r.Body).Decode(&homework); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	rand.Seed(time.Now().UnixNano())
	homework.StudentLink = api.RandStringSeq(11)
	homework.TeacherLink = api.RandStringSeq(11)

	if user != nil {
		homework.TeacherID = user.ID
	}

	// Create homework in the database.
	err := s.HomeworkService.CreateHomework(r.Context(), &homework)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Write new homework content to response based on accept header.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		StudentLink string `json:"StudentLink"`
		TeacherLink string `json:"TeacherLink"`
	}{
		StudentLink: homework.StudentLink,
		TeacherLink: homework.TeacherLink,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleHomeworkUpdate handles the "PATCH /homeworks/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleHomeworkUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse homework ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	upd := api.HomeworkUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the homework in the database.
	_, err = s.HomeworkService.UpdateHomework(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleHomeworkDelete handles the "DELETE /homeworks/:id" route. This route
// permanently deletes the homework.
func (s *Server) handleHomeworkDelete(w http.ResponseWriter, r *http.Request) {
	// Parse homework ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the homework from the database.
	if err := s.HomeworkService.DeleteHomework(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
