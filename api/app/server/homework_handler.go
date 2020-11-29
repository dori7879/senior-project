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

// HandleListHomework is a handler that lists homework
func (server *Server) HandleListHomework(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to list homeworks")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "teacher" {
		server.logger.Warn().Err(errors.New("Teacher tries to list homeworks")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	homeworks, err := repository.ListHomeworksByOwnerEmail(server.db, claims["sub"].(string))
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	if homeworks == nil {
		fmt.Fprintf(w, "[]")
		return
	}

	dtos := homeworks.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleCreateHomework is a handler for creating a homework
func (server *Server) HandleCreateHomework(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		server.logger.Warn().Msgf("val: %v", val)
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "teacher" {
			server.logger.Warn().Err(errors.New("Registered teacher tries to create a homework")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	form := &model.HomeworkForm{}
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

	homeworkModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}
	server.logger.Warn().Msgf("claims: %v", claims)
	if claims != nil && claims["role"].(string) == "student" {
		user, err := repository.GetUserByEmail(server.db, claims["sub"].(string))
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
			return
		}
		server.logger.Warn().Err(err).Msgf("StudentID: %v", user.ID)
		homeworkModel.StudentID = user.ID
	}

	homework, err := repository.CreateHomework(server.db, homeworkModel)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
		return
	}

	server.logger.Info().Msgf("New homework  created: %d", homework.ID)
	w.WriteHeader(http.StatusCreated)
}

// HandleReadHomework is a handler for getting a single homework
func (server *Server) HandleReadHomework(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to read a homework")).Msg("")

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

	var homework *model.Homework

	if claims["role"].(string) == "teacher" {
		homework, err = repository.ReadHomeworkByIDandTeacher(server.db, uint(id), claims["sub"].(string))
	} else {
		homework, err = repository.ReadHomeworkByIDandOwner(server.db, uint(id), claims["sub"].(string))
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

	dto := homework.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleUpdateHomework is a handler for updating a homework
func (server *Server) HandleUpdateHomework(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to read a homework")).Msg("")

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

	form := &model.HomeworkForm{}
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

	homeworkModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	homeworkModel.ID = uint(id)

	if claims["role"].(string) == "teacher" {
		err = repository.UpdateHomeworkByTeacher(server.db, homeworkModel, claims["sub"].(string))
	} else {
		err = repository.UpdateHomeworkByOwner(server.db, homeworkModel, claims["sub"].(string))
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

	server.logger.Info().Msgf("Homework updated: %d", id)
	w.WriteHeader(http.StatusAccepted)
}

// HandleDeleteHomework is a handler for deleting a homework
func (server *Server) HandleDeleteHomework(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to delete a homework")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "student" {
		server.logger.Warn().Err(errors.New("Student tries to delete a homework")).Msg("")

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

	if err := repository.DeleteHomeworkByTeacher(server.db, uint(id), claims["sub"].(string)); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	server.logger.Info().Msgf("Homework  deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
