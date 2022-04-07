package helper

import (
	"encoding/json"
	"net/http"
)

func ReadReqJson(request *http.Request, result interface{}){
	decoder := json.NewDecoder(request.Body)
	err := decoder.Decode(result)
	PanicIfError(err)
}

func WriteResJson(w http.ResponseWriter, status int, data interface{}) {
	wrapper := make(map[string]interface{})
	wrapper["data"] = data

	js, err := json.Marshal(wrapper)
	PanicIfError(err)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}
