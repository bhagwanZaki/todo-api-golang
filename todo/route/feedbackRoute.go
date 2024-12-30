package route

import (
	"context"
	"net/http"
	"strings"
	"github.com/google/uuid"

	"common/response"
	"common/uploader"
	"todoGoApi/service"
	"todoGoApi/types"
)

type FeedbackApi struct{}

func (h *FeedbackApi) CreateFeedbackAPI(w http.ResponseWriter, r *http.Request,userData types.User, cld *uploader.ImageWrappper) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the file from the form data
	feedback := r.FormValue("feedback")
	if feedback == "" {
		response.ErrorResponse(w, "Unable to retrieve feedback from form", http.StatusBadRequest, "CreateFeedbackAPI")
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		response.ErrorResponse(w, "Unable to retrieve file from form", http.StatusBadRequest, "CreateFeedbackAPI")
		return
	}
	defer file.Close()
	
	if fileHeader.Size > 100 << 20{
		response.ErrorResponse(w, "File Size more than 200mb", http.StatusBadRequest, "CreateFeedbackAPI")
		return 
	}

	id := uuid.New()
	uniqueHash := strings.SplitN(id.String(), "-", 2)[0]
	fileName :=uniqueHash+fileHeader.Filename
	upload, err := cld.Upload(context.Background(), file,fileName)
	
	if err != nil{
		response.ErrorResponse(w,err.Error(), http.StatusBadRequest, "CreateFeedbackAPI")
		return
	}

	statusCode, err := service.CreateFeedbackService(userData, feedback, upload.SecureURL)

	if err != nil {
		response.ErrorResponse(w, err.Error(), statusCode, "CreateFeedbackAPI")
		return
	}

	w.WriteHeader(statusCode)
}