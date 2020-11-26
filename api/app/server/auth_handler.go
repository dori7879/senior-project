package server

import (
	"api/model"
	"api/repository"
	"api/util/auth"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

// HandleSignUp is a handler that registers new users
func (server *Server) HandleSignUp(w http.ResponseWriter, r *http.Request) {
	form := &model.UserForm{}
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

	userModel, err := form.ToModel()
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnprocessableEntity)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormDecodingFailure)
		return
	}

	user, err := repository.CreateUser(server.db, userModel)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataCreationFailure)
		return
	}

	server.logger.Info().Msgf("New user page created: %d", user.ID)
	w.WriteHeader(http.StatusCreated)
}

// HandleLogin is a handler that authenticates users
func (server *Server) HandleLogin(w http.ResponseWriter, r *http.Request) {
	form := &model.UserForm{}
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

	user, err := repository.GetUserByEmail(server.db, form.Email)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	err = bcrypt.CompareHashAndPassword(user.PasswordHash, []byte(form.Password))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusUnauthorized)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrInvalidCredentials)
		return
	} else if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrCompareHashPasswordFailure)
		return
	}

	var userRole string

	_, err = repository.GetStudent(server.db, user.ID)
	if err != nil {
		userRole = "teacher"
	} else {
		userRole = "student"
	}

	aToken, rToken, err := server.jwtUtils.GenerateTokenPair(user.Email, userRole)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrTokenGenerationFailure)
		return
	}

	respBody := map[string]string{
		"access_token":  aToken,
		"refresh_token": rToken,
		"role":          userRole,
		"first_name":    user.FirstName,
		"last_name":     user.LastName,
		"email":         user.Email,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// HandleTokenRefresh is a handler that refreshes both tokens
func (server *Server) HandleTokenRefresh(w http.ResponseWriter, r *http.Request) {
	form := &model.RefreshTokenForm{}
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

	valid, email, keyID, err := server.jwtUtils.ValidateToken(form.RefreshToken)
	if err != nil {
		// Respond with 500 status code
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrRefreshTokenValidationFailure)
		return
	}

	if !valid || keyID != auth.RefreshTokenSecretKey {
		// Response contains {"error": "not valid"}
		server.logger.Warn().Err(errors.New("Either refresh token is invalid or access token was passed instead")).Msg("")
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrRefreshTokenValidationFailure)
		return
	}

	user, err := repository.GetUserByEmail(server.db, email)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrDataAccessFailure)
		return
	}

	var userRole string

	_, err = repository.GetStudent(server.db, user.ID)
	if err != nil {
		userRole = "teacher"
	} else {
		userRole = "student"
	}

	aToken, rToken, err := server.jwtUtils.GenerateTokenPair(user.Email, userRole)
	if err != nil {
		server.logger.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrTokenGenerationFailure)
		return
	}

	respBody := map[string]string{
		"access_token":  aToken,
		"refresh_token": rToken,
		"role":          userRole,
	}

	if err := json.NewEncoder(w).Encode(respBody); err != nil {
		server.logger.Warn().Err(err).Msg("")

		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}
	w.WriteHeader(http.StatusOK)
}
