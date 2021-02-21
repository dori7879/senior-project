package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"api/app/router/middleware"
	"api/model"
	"api/repository"
	"api/util/linkgeneration"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

// HandleListQuiz is a handler that lists quizzes
func (server *Server) HandleListQuiz(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to list quiz pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "student" {
		server.logger.Warn().Err(errors.New("Student tries to list quiz pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	quizzes, err := repository.ListQuizzesByOwner(server.db, claims["sub"].(string))
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	if quizzes == nil {
		fmt.Fprintf(w, "[]")
		return
	}

	dtos := quizzes.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleCreateQuiz is a handler that creates a quiz
func (server *Server) HandleCreateQuiz(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Registered student tries to create a quiz page")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	form := &model.QuizForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v", "form": "%+v"}`, serverErrFormDecodingFailure, form)
		return
	}

	if err := server.validator.Struct(form); err != nil {
		handleValidationError(w, server.logger, err)
		return
	}

	quizModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v", "form": "%+v"}`, serverErrFormDecodingFailure, form)
		return
	}

	// Add registered user if the request from one
	rand.Seed(time.Now().UnixNano())
	studentRandomString := linkgeneration.RandStringSeq(11)
	teacherRandomString := linkgeneration.RandStringSeq(11)

	quizModel.StudentLink = studentRandomString
	quizModel.TeacherLink = teacherRandomString

	if claims != nil && claims["role"].(string) == "teacher" {
		user, err := repository.GetUserByEmail(server.db, claims["sub"].(string))
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
			return
		}
		server.logger.Warn().Err(err).Msgf("TeacherID: %v", &user.ID)
		quizModel.TeacherID = user.ID
	}

	quiz, err := repository.CreateQuiz(server.db, quizModel)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
		return
	}

	for _, oq := range form.OpenQuestions {
		oqModel, err := oq.ToModel()
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"error": "%v", "form": "%+v"}`, serverErrFormDecodingFailure, form)
			return
		}

		oqModel.QuizID = quiz.ID

		_, err = repository.CreateOpenQuestion(server.db, oqModel)
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
			return
		}
	}

	for _, tfq := range form.TrueFalseQuestions {
		tfqModel, err := tfq.ToModel()
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"error": "%v", "form": "%+v"}`, serverErrFormDecodingFailure, form)
			return
		}

		tfqModel.QuizID = quiz.ID

		_, err = repository.CreateTrueFalseQuestion(server.db, tfqModel)
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
			return
		}
	}

	for _, mcq := range form.MultipleChoiceQuestions {
		mcqModel, err := mcq.ToModel()
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprintf(w, `{"error": "%v", "form": "%+v"}`, serverErrFormDecodingFailure, form)
			return
		}

		mcqModel.QuizID = quiz.ID

		mcqObj, err := repository.CreateMultipleChoiceQuestion(server.db, mcqModel)
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
			return
		}

		for _, ac := range mcq.AnswerChoices {
			acModel, err := ac.ToModel()
			if err != nil {
				server.logger.Warn().Err(err).Msg("")

				w.WriteHeader(http.StatusUnprocessableEntity)
				fmt.Fprintf(w, `{"error": "%v", "form": "%+v"}`, serverErrFormDecodingFailure, form)
				return
			}

			acModel.QuestionID = mcqObj.ID

			_, err = repository.CreateAnswerChoice(server.db, acModel)
			if err != nil {
				server.logger.Warn().Err(err).Msg("")

				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
				return
			}
		}
	}

	respBody := map[string]string{
		"teacher_link": quizModel.TeacherLink,
		"student_link": quizModel.StudentLink,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
	w.WriteHeader(http.StatusCreated)
	server.logger.Info().Msgf("New quiz page created: %d", quiz.ID)
}

// HandleReadQuiz is a handler for getting a single quiz page
func (server *Server) HandleReadQuiz(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Student tries to read a quiz page")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var quiz *model.Quiz

	if claims != nil {
		quiz, err = repository.ReadQuizByOwner(server.db, uint(id), claims["sub"].(string))
	} else {
		quiz, err = repository.ReadQuizWithNoOwner(server.db, uint(id))
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	dto := quiz.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleReadQuizByTeacherLink is a handler for getting a single quiz page
func (server *Server) HandleReadQuizByTeacherLink(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Student tries to read a quiz page")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	str := chi.URLParam(r, "str")
	if str == "" {
		server.logger.Info().Msgf("can not parse str: %v", str)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var quiz *model.Quiz
	var err error
	if claims != nil {
		quiz, err = repository.ReadQuizWithOwnerByTeacherLink(server.db, str, claims["sub"].(string))
	} else {
		quiz, err = repository.ReadQuizWithNoOwnerByTeacherLink(server.db, str)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	qss, err := repository.ListRelatedQuizSubmissions(server.db, quiz.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	dto := quiz.ToNestedDto(qss, model.OpenQuestions{}, model.TrueFalseQuestions{}, model.MultipleChoiceQuestions{})
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleReadQuizByStudentLink is a handler for getting a single quiz page
func (server *Server) HandleReadQuizByStudentLink(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "teacher" {
			server.logger.Warn().Err(errors.New("Teacher tries to read a quiz page for a student link")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	str := chi.URLParam(r, "str")
	if str == "" {
		server.logger.Info().Msgf("can not parse str: %v", str)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var quiz *model.Quiz
	var err error
	if claims != nil {
		quiz, err = repository.ReadQuizByOwnerByStudentLink(server.db, str)
	} else {
		quiz, err = repository.ReadQuizWithNoOwnerByStudentLink(server.db, str)
	}
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	oqs, err := repository.ListRelatedOpenQuestions(server.db, quiz.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	tfqs, err := repository.ListRelatedTrueFalseQuestions(server.db, quiz.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	mcqs, err := repository.ListRelatedMultipleChoiceQuestions(server.db, quiz.ID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	for _, mcq := range mcqs {
		answerChoices, err := repository.ListRelatedAnswerChoices(server.db, mcq.ID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				w.WriteHeader(http.StatusNotFound)
				return
			}

			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
			return
		}

		mcq.AnswerChoices = answerChoices
	}

	dto := quiz.ToNestedDto(nil, oqs, tfqs, mcqs)
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleUpdateQuiz is a handler for updating a quiz page
func (server *Server) HandleUpdateQuiz(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Student tries to read a quiz pages")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	form := &model.QuizForm{}
	if err := json.NewDecoder(r.Body).Decode(form); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	if err := server.validator.Struct(form); err != nil {
		handleValidationError(w, server.logger, err)
		return
	}

	quizModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	quizModel.ID = uint(id)

	if claims != nil {
		err = repository.UpdateQuizByOwner(server.db, quizModel, claims["sub"].(string))
	} else {
		err = repository.UpdateQuizWithNoOwner(server.db, quizModel)
	}

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataUpdateFailure)
		return
	}

	server.logger.Info().Msgf("Quiz updated: %d", id)
	w.WriteHeader(http.StatusAccepted)
}

// HandleDeleteQuiz is a handler for deleting a quiz page
func (server *Server) HandleDeleteQuiz(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to delete quiz pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "student" {
		server.logger.Warn().Err(errors.New("Student tries to delete quiz pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	if err := repository.DeleteQuizByOwner(server.db, uint(id), claims["sub"].(string)); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	server.logger.Info().Msgf("Quiz deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
