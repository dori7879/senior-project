package server

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"

	"api/util/auth"
	"api/util/logger"
)

const (
	serverErrDataAccessFailure             = "data access failure"
	serverErrDataAccessUnauthorized        = "data access unauthorized"
	serverErrJSONCreationFailure           = "json creation failure"
	serverErrDataCreationFailure           = "data creation failure"
	serverErrFormDecodingFailure           = "form decoding failure"
	serverErrDataUpdateFailure             = "data update failure"
	serverErrFormErrResponseFailure        = "form error response failure"
	serverErrInvalidCredentials            = "invalid credentials"
	serverErrCompareHashPasswordFailure    = "hash and password compare failure"
	serverErrRefreshTokenValidationFailure = "refresh token validation failure"
	serverErrTokenGenerationFailure        = "token generation failure"
)

// Server is a struct that has various other structs that all are available in handlers
type Server struct {
	logger    *logger.Logger
	db        *gorm.DB
	validator *validator.Validate
	jwtUtils  *auth.JwtUtils
}

// New is a server constructor
func New(logger *logger.Logger, db *gorm.DB, validator *validator.Validate, jwtUtils *auth.JwtUtils) *Server {
	return &Server{
		logger:    logger,
		db:        db,
		validator: validator,
		jwtUtils:  jwtUtils,
	}
}

// Logger is a function that returns the logger instance
func (server *Server) Logger() *logger.Logger {
	return server.logger
}
