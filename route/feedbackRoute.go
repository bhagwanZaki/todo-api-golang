package route

import (
	"log"
	"net/http"
	"todoGoApi/common"
	"todoGoApi/types"
	"context"
	"strings"
	"todoGoApi/service"
	"github.com/google/uuid"
)

type FeedbackApi struct{}

func (h *FeedbackApi) CreateFeedbackAPI(w http.ResponseWriter, r *http.Request,userData types.User, cld *common.ImageWrappper) {
	w.Header().Set("Content-Type", "application/json")

	// Retrieve the file from the form data
	feedback := r.FormValue("feedback")
	if feedback == "" {
		common.ErrorResponse(w, "Unable to retrieve feedback from form", http.StatusBadRequest, "CreateFeedbackAPI")
		return
	}

	file, fileHeader, err := r.FormFile("image")
	if err != nil {
		common.ErrorResponse(w, "Unable to retrieve file from form", http.StatusBadRequest, "CreateFeedbackAPI")
		return
	}
	defer file.Close()
	
	if fileHeader.Size > 100 << 20{
		common.ErrorResponse(w, "File Size more than 200mb", http.StatusBadRequest, "CreateFeedbackAPI")
		return 
	}

	id := uuid.New()
	uniqueHash := strings.SplitN(id.String(), "-", 2)[0]
	fileName :=uniqueHash+fileHeader.Filename
	upload, err := cld.Upload(context.Background(), file,fileName)
	
	if err != nil{
		common.ErrorResponse(w,err.Error(), http.StatusBadRequest, "CreateFeedbackAPI")
		return
	}

	statusCode, err := service.CreateFeedbackService(userData, feedback, upload.SecureURL)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "CreateFeedbackAPI")
		return
	}

	w.WriteHeader(statusCode)
}