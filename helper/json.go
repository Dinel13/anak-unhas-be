package helper

import (
	"encoding/json"
	"net/http"
)

func ReadJson(request *http.Request, result interface{}) error {
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	if err != nil {
		return err
	}
	return nil
}

func WriteJson(w http.ResponseWriter, status int, data interface{}, wrap string) error {
	wraper := make(map[string]interface{})

	wraper[wrap] = data

	js, err := json.Marshal(wraper)
	PanicIfError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)

	return nil
}

func WriteJsonError(w http.ResponseWriter, err error, status ...int) {
	statusCode := http.StatusBadRequest
	if len(status) > 0 {
		statusCode = status[0]
	}

	logErorr(err)

	type jsonError struct {
		Message string `json:"message"`
	}

	theError := jsonError{
		Message: err.Error(),
	}

	WriteJson(w, statusCode, theError, "error")
}
