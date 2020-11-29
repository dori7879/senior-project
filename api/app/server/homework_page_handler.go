package server

import (
	"api/app/router/middleware"
	"api/model"
	"api/repository"
	"api/util/linkgeneration"
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
	"github.com/jinzhu/gorm"
)

// HandleListHomeworkPage is a handler that lists homework pages
func (server *Server) HandleListHomeworkPage(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to list homework pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "student" {
		server.logger.Warn().Err(errors.New("Student tries to list homework pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	homeworkPages, err := repository.ListHomeworkPagesByOwner(server.db, claims["sub"].(string))
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	if homeworkPages == nil {
		fmt.Fprintf(w, "[]")
		return
	}

	dtos := homeworkPages.ToDto()
	if err := json.NewEncoder(w).Encode(dtos); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleCreateHomeworkPage is a handler for creating a homework page
func (server *Server) HandleCreateHomeworkPage(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Registered student tries to create a homework page")).Msg("")

			w.WriteHeader(http.StatusUnauthorized)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
			return
		}
	}

	form := &model.HomeworkPageForm{}
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

	homeworkPageModel, err := form.ToModel()
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

	homeworkPageModel.StudentLink = studentRandomString
	homeworkPageModel.TeacherLink = teacherRandomString

	if claims != nil && claims["role"].(string) == "teacher" {
		user, err := repository.GetUserByEmail(server.db, claims["sub"].(string))
		if err != nil {
			server.logger.Warn().Err(err).Msg("")

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
			return
		}
		server.logger.Warn().Err(err).Msgf("TeacherID: %v", &user.ID)
		homeworkPageModel.TeacherID = user.ID
	}

	homeworkPage, err := repository.CreateHomeworkPage(server.db, homeworkPageModel)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
		return
	}

	respBody := map[string]string{
		"teacher_link": homeworkPageModel.TeacherLink,
		"student_link": homeworkPageModel.StudentLink,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
	w.WriteHeader(http.StatusCreated)
	server.logger.Info().Msgf("New homework page created: %d", homeworkPage.ID)
}

// HandleReadHomeworkPage is a handler for getting a single homework page
func (server *Server) HandleReadHomeworkPage(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Student tries to read a homework pages")).Msg("")

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

	var homeworkPage *model.HomeworkPage

	if claims != nil {
		homeworkPage, err = repository.ReadHomeworkPageByOwner(server.db, uint(id), claims["sub"].(string))
	} else {
		homeworkPage, err = repository.ReadHomeworkPageWithNoOwner(server.db, uint(id))
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

	dto := homeworkPage.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleReadHomeworkPageByTeacherLink is a handler for getting a single homework page
func (server *Server) HandleReadHomeworkPageByTeacherLink(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Student tries to read a homework pages")).Msg("")

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

	var homeworkPage *model.HomeworkPage
	var err error
	if claims != nil {
		homeworkPage, err = repository.ReadHomeworkPageWithOwnerByTeacherLink(server.db, str, claims["sub"].(string))
	} else {
		homeworkPage, err = repository.ReadHomeworkPageWithNoOwnerByTeacherLink(server.db, str)
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

	hws, err := repository.ListRelatedHomeworks(server.db, homeworkPage.ID)
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

	dto := homeworkPage.ToNestedDto(hws)
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleReadHomeworkPageByStudentLink is a handler for getting a single homework page
func (server *Server) HandleReadHomeworkPageByStudentLink(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "teacher" {
			server.logger.Warn().Err(errors.New("Teacher tries to read a homework page for a student link")).Msg("")

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

	var homeworkPage *model.HomeworkPage
	var err error
	if claims != nil {
		homeworkPage, err = repository.ReadHomeworkPageByOwnerByStudentLink(server.db, str)
	} else {
		homeworkPage, err = repository.ReadHomeworkPageWithNoOwnerByStudentLink(server.db, str)
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

	dto := homeworkPage.ToDto()
	if err := json.NewEncoder(w).Encode(dto); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
}

// HandleUpdateHomeworkPage is a handler for updating a homework page
func (server *Server) HandleUpdateHomeworkPage(w http.ResponseWriter, r *http.Request) {
	var claims jwt.MapClaims
	if val := r.Context().Value(middleware.CtxKeyJWTClaims); val != nil {
		claims = val.(jwt.MapClaims)
		if claims["role"].(string) == "student" {
			server.logger.Warn().Err(errors.New("Student tries to read a homework pages")).Msg("")

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

	form := &model.HomeworkPageForm{}
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

	homeworkPageModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	homeworkPageModel.ID = uint(id)

	if claims != nil {
		err = repository.UpdateHomeworkPageByOwner(server.db, homeworkPageModel, claims["sub"].(string))
	} else {
		err = repository.UpdateHomeworkPageWithNoOwner(server.db, homeworkPageModel)
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

// HandleDeleteHomeworkPage is a handler for deleting a homework page
func (server *Server) HandleDeleteHomeworkPage(w http.ResponseWriter, r *http.Request) {
	var val interface{}
	if val = r.Context().Value(middleware.CtxKeyJWTClaims); val == nil {
		server.logger.Warn().Err(errors.New("Guest tries to delete homework pages")).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessUnauthorized)
		return
	}

	claims := val.(jwt.MapClaims)

	if claims["role"].(string) == "student" {
		server.logger.Warn().Err(errors.New("Student tries to delete homework pages")).Msg("")

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

	if err := repository.DeleteHomeworkPageByOwner(server.db, uint(id), claims["sub"].(string)); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	server.logger.Info().Msgf("Homework page deleted: %d", id)
	w.WriteHeader(http.StatusAccepted)
}
