package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	StatusOk    = "OK"
	StatusError = "Error"
)

func OK() Response {
	return Response{Status: StatusOk}
}


func Error(msg string)Response{
return Response{Status: StatusError,Error: msg}
}

func ValidationErrors(errs validator.ValidationErrors) Response{

	var errMsg []string

	for _,err :=range errs{
		switch err.ActualTag(){
		case "required":
			errMsg = append(errMsg,fmt.Sprint("field %s is a required ",err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprint("field %s is not valid url Url",err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprint("field %s is not valid",err.Field()))
		}


	}
	return Response{
		Status: StatusError,
		Error: strings.Join(errMsg,", "),
	}
}