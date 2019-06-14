package util

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
)

func CheckId(w http.ResponseWriter, idString string) (uint64, error) {

	id, err := strconv.ParseUint(idString, 10, 64)
	if err != nil {
		SendErr(w, http.StatusBadRequest, errors.New(idString+" is not integer"))
	}

	return id, err

}

func jsonError(message string) string {

	var err = struct {
		ErrorMessage string `json:"error_message"`
	}{message}

	outputJson, _ := json.Marshal(err)
	return string(outputJson)

}

func SendErr(w http.ResponseWriter, status int, err error) {

	w.WriteHeader(status)
	fmt.Fprint(w, jsonError(err.Error()))

}

func SendNoConnectionWithDb(w http.ResponseWriter) {

	w.WriteHeader(http.StatusServiceUnavailable)
	fmt.Fprint(w, jsonError("No database connection"))

}
