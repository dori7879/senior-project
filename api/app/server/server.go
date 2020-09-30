package server

import (
	"api/util/logger"
)

type Server struct {
	logger *logger.Logger
}

func New(logger *logger.Logger) *Server {
	return &Server{logger: logger}
}

func (server *Server) Logger() *logger.Logger {
	return server.logger
}
