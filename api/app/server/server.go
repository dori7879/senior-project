package server

import (
	"github.com/jinzhu/gorm"

	"api/util/logger"
)

const (
	serverErrDataAccessFailure   = "data access failure"
	serverErrJsonCreationFailure = "json creation failure"
)

type Server struct {
	logger *logger.Logger
	db     *gorm.DB
}

func New(logger *logger.Logger, db *gorm.DB) *Server {
	return &Server{
		logger: logger,
		db:     db,
	}
}

func (server *Server) Logger() *logger.Logger {
	return server.logger
}
