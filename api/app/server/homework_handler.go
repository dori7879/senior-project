package server

import (
	"api/repository"
	"encoding/json"
	"fmt"
	"net/http"
)

func (server *Server) HandleListHomework(w http.ResponseWriter, r *http.Request) {
	homeworks, err := repository.ListHomeworks(server.db)
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
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJsonCreationFailure)
		return
	}
}

func (server *Server) HandleCreateHomework(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusCreated)
}

func (server *Server) HandleReadHomework(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("{}"))
}

func (server *Server) HandleUpdateHomework(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}

func (server *Server) HandleDeleteHomework(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusAccepted)
}
