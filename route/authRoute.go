package route

import (
	"encoding/json"
	"log"
	"net/http"
	"todoGoApi/common"
	"todoGoApi/service"
	"todoGoApi/types"
)

type AuthApi struct{}

func (h *AuthApi) CheckAuthentication(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	userData, err := common.GetUserDataFromToken(r)

	if err != nil {
		log.Println(err.Error())
		log.Println(err.Error() == "no rows in result set")
		if err.Error() == "no rows in result set" {
			log.Println("if satisfied")
			common.ErrorResponse(w, "invalid authorization", http.StatusUnauthorized)
			return
		}
		common.ErrorResponse(w, err.Error(), http.StatusUnauthorized)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}
func (h *AuthApi) LoginAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.LoginSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	if body.Username == "" || body.Password == "" {
		common.ErrorResponse(w, "Invalid body", http.StatusBadRequest)
		return
	}

	userData, token, err := service.LoginService(body)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}
	common.SetTokenCookie(w, token)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(userData)
}

func (h *AuthApi) RegisterApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.RegisterSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	if body.Username == "" || body.Email == "" {
		common.ErrorResponse(w, "Invalid body", http.StatusBadRequest)
		return
	}
	log.Println("REGISTER BODY ",body)
	registerData, err := service.RegisterService(body)

	if err != nil {
		log.Println("REGISTER SERVICE ERROR ",err.Error())
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(registerData)
}

func (h *AuthApi) CompleteRegisterApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.VerifyOtpAndRegisterSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	if body.Username == "" || body.Email == "" || body.Password == "" || body.Fullname == "" {
		common.ErrorResponse(w, "Invalid body", http.StatusBadRequest)
		return
	}

	userData, token, err := service.VerifyOtpAndCompleteRegistration(body)

	if err != nil {
		if err.Error() == "invalid OTP" {
			common.ErrorResponse(w, "Invalid OTP", http.StatusBadRequest)
			return
		}
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
		return
	}
	
	common.SetTokenCookie(w, token)
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(userData)
}

func (h *AuthApi) LogoutApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	token, _ := common.CheckTokenValidity(r)

	userLogout, err := service.LogoutService(userData.Id, token)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	if !userLogout {
		common.ErrorResponse(w, "User Not Found", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) LogoutFromAllDeviceApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	userLogout, err := service.LogoutFromAllDeviceService(userData.Id)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	if !userLogout {
		common.ErrorResponse(w, "User Not Found", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) CreateOTPRequestAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.OtpRequestSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	status, err := service.CreateOTPRequest(body.Email, body.RequestTypeCode)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !status {
		common.ErrorResponse(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) VerifyPasswordChangeRequestAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.VerifyOtpSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	status, err := service.VerifyPasswordChangeOTP(body.Email, body.Otp, body.IsForgetPasswordReqest)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !status {
		common.ErrorResponse(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) ForgetPasswordAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.PasswordChangeSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	status, err := service.ChangeUserPasswordService(body.Email, body.Password, true)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !status {
		common.ErrorResponse(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) ResetPasswordAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.PasswordChangeSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest)
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError)
		return
	}

	status, err := service.ChangeUserPasswordService(body.Email, body.Password, false)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusBadRequest)
		return
	}

	if !status {
		common.ErrorResponse(w, "Something went wrong", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) DeleteUserApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	userDeleted, err := service.DeleteAccountService(userData.Id, userData.Email)

	if err != nil {
		common.ErrorResponse(w, err.Error(), http.StatusInternalServerError)
	}

	if !userDeleted {
		common.ErrorResponse(w, "User Not Found", http.StatusBadRequest)
	}

	w.WriteHeader(http.StatusNoContent)
}
