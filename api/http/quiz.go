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

// registerQuizPrivateRoutes is a helper function for registering private quiz routes.
func (s *Server) registerQuizPrivateRoutes(r *mux.Router) {
	// Listing of all quizzes a teacher is an owner of.
	r.HandleFunc("/quizzes", s.handleQuizList).Methods("GET")

	// View a single quiz.
	r.HandleFunc("/quizzes/{id}", s.handleQuizView).Methods("GET")
}

// registerQuizPublicRoutes is a helper function for registering public quiz routes.
func (s *Server) registerQuizPublicRoutes(r *mux.Router) {
	// API endpoint for creating quizzes.
	r.HandleFunc("/quizzes", s.handleQuizCreate).Methods("POST")

	// View a single quiz.
	r.HandleFunc("/quizzes/shared/{link}/teacher", s.handleQuizTeacherView).Methods("GET")
	r.HandleFunc("/quizzes/shared/{link}/student", s.handleQuizStudentView).Methods("GET")

	r.HandleFunc("/quizzes/{id}", s.handleQuizUpdate).Methods("PATCH")

	// Removing a quiz.
	r.HandleFunc("/quizzes/{id}", s.handleQuizDelete).Methods("DELETE")

}

// handleQuizList handles the "GET /quizzes" route. This route can optionally
// accept filter arguments and outputs a list of all quizzes that the current
// user is related to.
func (s *Server) handleQuizList(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user == nil {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You must be logged in"))
		return
	} else if !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Fetch quizzes from database.
	quizzes, n, err := s.QuizService.FindQuizzes(r.Context(), api.QuizFilter{TeacherID: &user.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Render output based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(struct {
		Quizs []*api.Quiz `json:"Quizs"`
		N     int         `json:"N"`
	}{
		Quizs: quizzes,
		N:     n,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizView handles the "GET /quizzes/:id" route.
func (s *Server) handleQuizView(w http.ResponseWriter, r *http.Request) {
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

	// Fetch quiz from the database.
	quiz, err := s.QuizService.FindQuizByID(r.Context(), id)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Fetch associated submissions and questions from the database.
	quiz.Submissions, _, err = s.QuizSubmissionService.FindQuizSubmissions(r.Context(), api.QuizSubmissionFilter{QuizID: &quiz.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	quiz.Questions, _, err = s.QuestionService.FindQuestions(r.Context(), api.QuestionFilter{QuizID: &quiz.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(quiz); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizTeacherView handles the "GET /quizzes/shared/:link/teacher" route.
func (s *Server) handleQuizTeacherView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Parse teacher link from path.
	link := mux.Vars(r)["link"]

	// Fetch quiz from the database.
	quiz, err := s.QuizService.FindQuizByTeacherLink(r.Context(), link)
	if err != nil {
		Error(w, r, err)
		return
	}

	// Fetch associated submissions and questions from the database.
	quiz.Submissions, _, err = s.QuizSubmissionService.FindQuizSubmissions(r.Context(), api.QuizSubmissionFilter{QuizID: &quiz.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	quiz.Questions, _, err = s.QuestionService.FindQuestions(r.Context(), api.QuestionFilter{QuizID: &quiz.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(quiz); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizStudentView handles the "GET /quizzes/shared/:link/studnet" route.
func (s *Server) handleQuizStudentView(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a student"))
		return
	}

	// Parse student link from path.
	link := mux.Vars(r)["link"]

	// Fetch quiz from the database.
	quiz, err := s.QuizService.FindQuizByStudentLink(r.Context(), link)
	if err != nil {
		Error(w, r, err)
		return
	}

	quiz.Questions, _, err = s.QuestionService.FindQuestions(r.Context(), api.QuestionFilter{QuizID: &quiz.ID})
	if err != nil {
		Error(w, r, err)
		return
	}

	if quiz.TeacherID != user.ID {
		w.Header().Set("Content-type", "application/json")
		w.Write([]byte(`{}`))
		return
	}

	// Format returned data based on HTTP accept header.
	w.Header().Set("Content-type", "application/json")
	if err := json.NewEncoder(w).Encode(quiz); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizCreate handles the "POST /quizzes" route.
func (s *Server) handleQuizCreate(w http.ResponseWriter, r *http.Request) {
	user := api.UserFromContext(r.Context())
	if user != nil && !user.IsTeacher {
		Error(w, r, api.Errorf(api.EUNAUTHORIZED, "You are not a teacher"))
		return
	}

	// Unmarshal quiz data first
	quiz := api.Quiz{}
	if err := json.NewDecoder(r.Body).Decode(&quiz); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	rand.Seed(time.Now().UnixNano())
	quiz.StudentLink = api.RandStringSeq(11)
	quiz.TeacherLink = api.RandStringSeq(11)

	if user != nil {
		quiz.TeacherID = user.ID
	}

	// Create quiz in the database.
	err := s.QuizService.CreateQuiz(r.Context(), &quiz)
	if err != nil {
		Error(w, r, err)
		return
	}

	for _, q := range quiz.Questions {
		q.QuizID = quiz.ID

		err = s.QuestionService.CreateQuestion(r.Context(), q)
		if err != nil {
			Error(w, r, err)
			return
		}
	}

	// Write new quiz content to response based on accept header.
	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(struct {
		StudentLink string `json:"StudentLink"`
		TeacherLink string `json:"TeacherLink"`
	}{
		StudentLink: quiz.StudentLink,
		TeacherLink: quiz.TeacherLink,
	}); err != nil {
		LogError(r, err)
		return
	}
}

// handleQuizUpdate handles the "PATCH /quizzes/:id" route. This route
// reads in the updated fields and issues an update in the database.
func (s *Server) handleQuizUpdate(w http.ResponseWriter, r *http.Request) {
	// Parse quiz ID from the path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Parse fields into an update object.
	upd := api.QuizUpdate{}
	if err := json.NewDecoder(r.Body).Decode(&upd); err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid JSON body"))
		return
	}

	// Update the quiz in the database.
	_, err = s.QuizService.UpdateQuiz(r.Context(), id, upd)
	if err != nil {
		Error(w, r, err)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}

// handleQuizDelete handles the "DELETE /quizzes/:id" route. This route
// permanently deletes the quiz.
func (s *Server) handleQuizDelete(w http.ResponseWriter, r *http.Request) {
	// Parse quiz ID from path.
	id, err := strconv.Atoi(mux.Vars(r)["id"])
	if err != nil {
		Error(w, r, api.Errorf(api.EINVALID, "Invalid ID format"))
		return
	}

	// Delete the quiz from the database.
	if err := s.QuizService.DeleteQuiz(r.Context(), id); err != nil {
		Error(w, r, err)
		return
	}

	// Response part
	w.Header().Set("Content-type", "application/json")
	w.Write([]byte(`{}`))
}
