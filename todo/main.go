package main

import (
	"common/env"
	"common/queue"
	"common/uploader"
	"context"
	"fmt"
	"log"
	"net/http"
	"time"
	"todoGoApi/db"
	"todoGoApi/middleware"
	"todoGoApi/publisher"
	"todoGoApi/route"
)

func main() {
	envErr := env.LoadEnv()

	if envErr != nil {
		log.Fatalln("Failed to env file")
		return
	}

	imageUploader, err := uploader.InitiailizeImageWrappper()

	if err != nil{
		return
	}

	router := http.NewServeMux()

	todoRoutes := &route.TodoApi{}
	authRoutes := &route.AuthApi{}
	feedbackRoutes := &route.FeedbackApi{}

	// ḍatabase setup
	db.InitDatabase()
	defer db.CloseDatabase()

	// sqs setup
	ch := queue.InitializeQueue()
	defer ch.Close()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	publisher.InitEmailQueue(ch,ctx)

	// auth routes
	router.HandleFunc("GET /api/auth/checkauth", middleware.Logger(authRoutes.CheckAuthentication))
	router.HandleFunc("POST /api/auth/login", middleware.Logger(authRoutes.LoginAPI))
	router.HandleFunc("POST /api/auth/register", middleware.Logger(authRoutes.RegisterApi))
	router.HandleFunc("POST /api/auth/register/complete", middleware.Logger(authRoutes.CompleteRegisterApi))
	router.HandleFunc("PUT /api/auth/update-user", middleware.Logger(middleware.AuthRequired(authRoutes.ProfileChangeApi)))
	router.HandleFunc("DELETE /api/auth/delete-user", middleware.Logger(middleware.AuthRequired(authRoutes.DeleteUserApi)))
	router.HandleFunc("POST /api/auth/logout", middleware.Logger(middleware.AuthRequired(authRoutes.LogoutApi)))
	router.HandleFunc("POST /api/auth/logout-from-all-device", middleware.Logger(middleware.AuthRequired(authRoutes.LogoutFromAllDeviceApi)))
	router.HandleFunc("POST /api/auth/otp-request", middleware.Logger(authRoutes.CreateOTPRequestAPI))
	router.HandleFunc("POST /api/auth/verify-otp", middleware.Logger(authRoutes.VerifyOtpAPI))
	router.HandleFunc("POST /api/auth/reset-password-request", middleware.Logger(middleware.AuthRequired(authRoutes.ResetPasswordRequestAPI)))
	router.HandleFunc("POST /api/auth/reset-password", middleware.Logger(middleware.AuthRequiredReturnToken(authRoutes.ResetPasswordAPI)))
	router.HandleFunc("POST /api/auth/forget-password", middleware.Logger(authRoutes.ForgetPasswordAPI))
	// todo routes
	router.HandleFunc("GET /api/health", middleware.Logger(todoRoutes.HealthCheckAPI))
	router.HandleFunc("GET /api/todos", middleware.Logger(middleware.AuthRequired(todoRoutes.GetTodos)))
	router.HandleFunc("POST /api/add", middleware.Logger(middleware.AuthRequired(todoRoutes.AddTodo)))
	router.HandleFunc("PUT /api/update/{id}", middleware.Logger(middleware.AuthRequired(todoRoutes.UpdateTodo)))
	router.HandleFunc("DELETE /api/delete/{id}", middleware.Logger(middleware.AuthRequired(todoRoutes.DeleteTodo)))
	// feedback routes
	router.HandleFunc("POST /api/feedback/create", middleware.Logger(middleware.ImageUploaderUrl(feedbackRoutes.CreateFeedbackAPI, imageUploader)))

	server := http.Server{
		Addr:    ":8000",
		Handler: router,
	}

	fmt.Println("Server listening on port :8000")
	server.ListenAndServe()
}
