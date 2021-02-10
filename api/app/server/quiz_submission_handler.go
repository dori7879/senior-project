package server

import (
	"api/app/router/middleware"
	"api/model"
	"api/repository"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

// HandleListQuizSubmission is a handler that lists quiz submission
func (server *Server) HandleListQuizSubmission(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to list quiz submissions")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "teacher" {
		server.logger.Warn().Err(errors.New("Teacher tries to list quiz submissions")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	quizSubmissions, err := repository.ListQuizSubmissionsByOwnerEmail(server.db, claims["sub"].(string))
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	if quizSubmissions == nil {
		fmt.Fprintf(w, "[]")
		return
	}

	dtos := quizSubmissions.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleCreateQuizSubmission is a handler for creating a quiz submission
func (server *Server) HandleCreateQuizSubmission(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
	}

	if claims != nil && claims["role"].(string) == "teacher" {
		server.logger.Warn().Err(errors.New("Registered teacher tries to create a quiz submission")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	form := &model.QuizSubmissionForm{}
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

	qSubmissionModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	if claims != nil && claims["role"].(string) == "student" {
		user, err := repository.GetUserByEmail(server.db, claims["sub"].(string))
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
			return
		}

		qSubmissionModel.StudentID = user.ID
		qSubmissionModel.StudentFullname = user.FirstName + " " + user.LastName
	} else {
		qs, err := repository.ReadQuiz(server.db, qSubmissionModel.QuizID)
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
			return
		}

		if qs.Mode == "registered" {
			server.logger.Info().Msg("Guest tries to submit a quiz submission to a registered-only Quiz")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
	}

	quizSubmission, err := repository.CreateQuizSubmission(server.db, qSubmissionModel)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
		return
	}

	server.logger.Info().Msgf("New quiz submission  created: %d", quizSubmission.ID)
	w.WriteHeader(http.StatusCreated)
}

// HandleReadQuizSubmission is a handler for getting a single quiz submission
func (server *Server) HandleReadQuizSubmission(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to read a quiz submission")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}
	claims := val.(jwt.MapClaims)

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	var quizSubmission *model.QuizSubmission

	if claims["role"].(string) == "teacher" {
		quizSubmission, err = repository.ReadQuizSubmissionByIDandTeacher(server.db, uint(id), claims["sub"].(string))
	} else {
		quizSubmission, err = repository.ReadQuizSubmissionByIDandOwner(server.db, uint(id), claims["sub"].(string))
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

	dto := quizSubmission.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleUpdateQuizSubmission is a handler for updating a quiz submission
func (server *Server) HandleUpdateQuizSubmission(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
	}

	id, err := strconv.ParseUint(chi.URLParam(r, "id"), 0, 64)
	if err != nil || id == 0 {
		server.logger.Info().Msgf("can not parse ID: %v", id)

		w.WriteHeader(http.StatusUnprocessableEntity)
		return
	}

	form := &model.QuizSubmissionForm{}
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

	qSubmissionModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	qSubmissionModel.ID = uint(id)

	if claims != nil && claims["role"].(string) == "teacher" {
		err = repository.UpdateQuizSubmissionByTeacher(server.db, qSubmissionModel, claims["sub"].(string))
	} else if claims != nil && claims["role"].(string) == "student" {
		err = repository.UpdateQuizSubmissionByOwner(server.db, qSubmissionModel, claims["sub"].(string))
	} else {
		err = repository.UpdateQuizSubmission(server.db, qSubmissionModel)
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

	server.logger.Info().Msgf("Quiz Submission updated: %d", id)
	w.WriteHeader(http.StatusAccepted)
}

// HandleDeleteQuizSubmission is a handler for deleting a quiz submission
func (server *Server) HandleDeleteQuizSubmission(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to delete a quiz submission")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "student" {
		server.logger.Warn().Err(errors.New("Student tries to delete a quiz submission")).Msg("")

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

	if err := repository.DeleteQuizSubmissionByTeacher(server.db, uint(id), claims["sub"].(string)); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	server.logger.Info().Msgf("Quiz Submission deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
