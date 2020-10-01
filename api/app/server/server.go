package server

import (
	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"

	"api/util/logger"
)

const (
	serverErrDataAccessFailure      = "data access failure"
	serverErrJsonCreationFailure    = "json creation failure"
	serverErrDataCreationFailure    = "data creation failure"
	serverErrFormDecodingFailure    = "form decoding failure"
	serverErrDataUpdateFailure      = "data update failure"
	serverErrFormErrResponseFailure = "form error response failure"
)

type Server struct {
	logger    *logger.Logger
	db        *gorm.DB
	validator *validator.Validate
}

func New(logger *logger.Logger, db *gorm.DB, validator *validator.Validate) *Server {
	return &Server{
		logger:    logger,
		db:        db,
		validator: validator,
	}
}

func (server *Server) Logger() *logger.Logger {
	return server.logger
}
