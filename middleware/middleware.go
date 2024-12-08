package middleware

import (
	"log"
	"net/http"
	"time"
	"todoGoApi/common"
	"todoGoApi/types"
)

const (
	Reset  = "\033[0m"
	Yellow = "\033[33m"
	Green  = "\033[32m"
	Blue   = "\033[34m"
	Red    = "\033[31m"
	Orange = "\033[38;5;214m" // ANSI code for orange (a close approximation)
)

// ColorMethod returns the colored method string based on HTTP method
func ColorMethod(method string) string {
	switch method {
	case "POST":
		return Green + method + Reset
	case "GET":
		return Blue + method + Reset
	case "DELETE":
		return Orange + method + Reset
	case "UPDATE":
		return Yellow + method + Reset
	default:
		return method
	}
}

func Logger(next func(http.ResponseWriter, *http.Request)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		timer := time.Now()
		next(w, r)
		log.Println(ColorMethod(r.Method), r.URL.Path, Yellow, time.Since(timer), Reset)
	}
}

func AuthRequired(next func(http.ResponseWriter, *http.Request, types.User)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, userDataErr := common.GetUserDataFromToken(r)

		if userDataErr != nil {
			common.ErrorResponse(w, userDataErr.Error(), http.StatusUnauthorized, "AuthRequired")
			log.Println(Yellow+"UnAuthorized"+Reset, ColorMethod(r.Method), r.URL.Path)
			return
		}

		next(w, r, userData)
	}
}

func AuthRequiredReturnToken(next func(http.ResponseWriter, *http.Request, types.User, string)) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, token, userDataErr := common.GetUserDataWithTokenFromToken(r)

		if userDataErr != nil {
			common.ErrorResponse(w, userDataErr.Error(), http.StatusUnauthorized, "AuthRequired")
			log.Println(Yellow+"UnAuthorized"+Reset, ColorMethod(r.Method), r.URL.Path)
			return
		}

		next(w, r, userData, token)
	}
}

func ImageUploaderUrl(next func(http.ResponseWriter, *http.Request, types.User, *common.ImageWrappper), imageUploader *common.ImageWrappper) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		userData, userDataErr := common.GetUserDataFromToken(r)

		if userDataErr != nil {
			common.ErrorResponse(w, userDataErr.Error(), http.StatusUnauthorized, "AuthRequired")
			log.Println(Yellow+"UnAuthorized"+Reset, ColorMethod(r.Method), r.URL.Path)
			return
		}

		next(w, r, userData, imageUploader)
	}
}
