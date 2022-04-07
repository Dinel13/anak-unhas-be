package exception

import (
	"github.com/dinel13/anak-unhas-be/helper"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func ErrorHandler(w http.ResponseWriter, r *http.Request, e interface{})  {
	if notFoundError(w,e) {
		return
	}

	if validationError(w, e){
		return
	}

	if badRequestError(w,e) {
		return
	}
	internalServerError(w,e)
}

func badRequestError(w http.ResponseWriter, e interface{}) bool {
	exception, ok := e.(BadRequestError)
	if ok {
		helper.WriteResJson(w, http.StatusBadRequest, exception.Error)
		return  true
	}
	return false
}

func internalServerError(w http.ResponseWriter, e interface{}) {
	helper.WriteResJson(w, http.StatusInternalServerError, e)
}

func validationError(w http.ResponseWriter, e interface{}) bool {
	exception, ok := e.(validator.ValidationErrors)
	if ok {
		helper.WriteResJson(w, http.StatusNotFound, exception.Error )
		return true
	}
	return false
}

func notFoundError(w http.ResponseWriter, e interface{}) bool {
	exception, ok := e.(NotFoundError)
	if ok {
		helper.WriteResJson(w, http.StatusNotFound, exception.Error )
		return true
	}
	return false
}
