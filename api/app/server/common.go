package server

import (
	"api/util/logger"
	"api/util/validator"
	"encoding/json"
	"fmt"
	"net/http"
)

func handleValidationError(w http.ResponseWriter, l *logger.Logger, err error) {
	l.Warn().Err(err).Msg("")

	resp := validator.ToErrResponse(err)
	if resp == nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrFormErrResponseFailure)
		return
	}

	respBody, err := json.Marshal(resp)
	if err != nil {
		l.Warn().Err(err).Msg("")
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "%v"}`, serverErrJSONCreationFailure)
		return
	}

	w.WriteHeader(http.StatusUnprocessableEntity)
	w.Write(respBody)
	return
}
