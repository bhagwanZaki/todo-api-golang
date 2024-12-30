package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"common/logger"
	"todoGoApi/db"
	"todoGoApi/types"
)

func CreateFeedbackService(userData types.User , feedBack string, imageUrl string) (int, error){
	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	
	_,err := db.DB_CONN.Exec(
		context.Background(),
		"call create_feedback($1, $2, $3, $4, $5, $6, $7)",
		userData.Id,
		userData.Username,
		userData.Email,
		feedBack,
		imageUrl,
		"todo",
		dbDate,
	)

	if err != nil {
		logger.Logger(err.Error(), "CreateFeedbackService")
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	return http.StatusNoContent ,nil
}