package utils

import (
	"fmt"
	"net/http"
)

func HandleError(err error, msg string, res http.ResponseWriter, status int) {
	fmt.Println(msg, ": ", err)
	http.Error(res, msg, status)
}
