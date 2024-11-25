package route

import (
	"encoding/json"
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
		if err.Error() == "no rows in result set" {
			common.ErrorResponse(w, "invalid authorization", http.StatusUnauthorized, "CheckAuthentication")
			return
		}
		common.ErrorResponse(w, err.Error(), http.StatusUnauthorized, "CheckAuthentication")
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
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest, "LoginAPI")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "LoginAPI")
		return
	}

	if body.Username == "" || body.Password == "" {
		common.ErrorResponse(w, "Invalid body", http.StatusBadRequest, "LoginAPI")
		return
	}

	userData, token,statusCode, err := service.LoginService(body)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "LoginAPI")
		return
	}
	common.SetTokenCookie(w, token)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(userData)
}

func (h *AuthApi) RegisterApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.RegisterSchema
	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest, "RegisterApi")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "RegisterApi")
		return
	}

	if body.Username == "" || body.Email == "" {
		common.ErrorResponse(w, "Invalid body", http.StatusBadRequest, "RegisterApi")
		return
	}
	registerData, statusCode, err := service.RegisterService(body)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "RegisterApi")
		return
	}

	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(registerData)
}

func (h *AuthApi) CompleteRegisterApi(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.VerifyOtpAndRegisterSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty request body", http.StatusBadRequest, "CompleteRegisterApi")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "CompleteRegisterApi")
		return
	}

	if body.Username == "" || body.Email == "" || body.Password == "" || body.Fullname == "" {
		common.ErrorResponse(w, "Invalid body", http.StatusBadRequest, "CompleteRegisterApi")
		return
	}

	userData, token, statusCode, err := service.VerifyOtpAndCompleteRegistration(body)

	if err != nil {
		if err.Error() == "invalid OTP" {
			common.ErrorResponse(w, "Invalid OTP", statusCode, "CompleteRegisterApi")
			return
		}
		common.ErrorResponse(w, err.Error(), statusCode, "CompleteRegisterApi")
		return
	}

	common.SetTokenCookie(w, token)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(userData)
}

func (h *AuthApi) LogoutApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	token, _ := common.CheckTokenValidity(r)

	userLogout, statusCode, err := service.LogoutService(userData.Id, token)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "LogoutApi")
	}

	if !userLogout {
		common.ErrorResponse(w, "User Not Found", statusCode, "LogoutApi")
	}

	w.WriteHeader(statusCode)
}

func (h *AuthApi) LogoutFromAllDeviceApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")

	userLogout, statusCode,err := service.LogoutFromAllDeviceService(userData.Id)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "LogoutFromAllDeviceApi")
	}

	if !userLogout {
		common.ErrorResponse(w, "User Not Found", statusCode, "LogoutFromAllDeviceApi")
	}

	w.WriteHeader(statusCode)
}

func (h *AuthApi) CreateOTPRequestAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.OtpRequestSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest, "CreateOTPRequestAPI")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "CreateOTPRequestAPI")
		return
	}

	statusCode, err := service.CreateOTPRequest(body.Email, body.RequestTypeCode)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "CreateOTPRequestAPI")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) VerifyOtpAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.VerifyOtpSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest, "VerifyPasswordChangeRequestAPI")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "VerifyPasswordChangeRequestAPI")
		return
	}

	statusCode, err := service.VerifyOTP(body.Email, body.Otp, body.RequestType)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "VerifyPasswordChangeRequestAPI")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) ForgetPasswordAPI(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var body types.PasswordChangeSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest, "ForgetPasswordAPI")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "ForgetPasswordAPI")
		return
	}

	statusCode, err := service.ForgetPasswordService(body.Email, body.Password)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "ForgetPasswordAPI")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) ResetPasswordRequestAPI(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	var body types.ResetPasswordSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError,"ResetPasswordRequestAPI")
		return
	}

	statusCode, err := service.ResetPasswordRequest(userData.Username, body.Password)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode,"ResetPasswordRequestAPI")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) ResetPasswordAPI(w http.ResponseWriter, r *http.Request, userData types.User,userToken string) {
	w.Header().Set("Content-Type", "application/json")
	var body types.PasswordChangeSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil || body.Email == "" {
		common.ErrorResponse(w, "Enter email", http.StatusBadRequest,"ResetPasswordAPI")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError,"ResetPasswordAPI")
		return
	}

	statusCode, err := service.ResetPasswordService(userData, userToken, body.Password)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode,"ResetPasswordAPI")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) DeleteUserApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	var body types.DeleteAccountSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Invalid Body", http.StatusBadRequest, "DeleteUserApi")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError, "DeleteUserApi")
		return
	}
	userDeleted, statusErr,err := service.DeleteAccountService(userData.Id, userData.Email, body.Otp)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusErr, "DeleteUserApi")
		return
	}

	if !userDeleted {
		common.ErrorResponse(w, "User Not Found", statusErr, "DeleteUserApi")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *AuthApi) ProfileChangeApi(w http.ResponseWriter, r *http.Request, userData types.User) {
	w.Header().Set("Content-Type", "application/json")
	var body types.ProfileSchema

	bodyErr := json.NewDecoder(r.Body).Decode(&body)

	if r.Body == nil {
		common.ErrorResponse(w, "Empty body", http.StatusBadRequest,"ProfileChangeApi")
		return
	}

	if bodyErr != nil {
		common.ErrorResponse(w, bodyErr.Error(), http.StatusInternalServerError,"ProfileChangeApi")
		return
	}

	res, statusCode, err := service.UpdateProfileService(body, userData)

	if err != nil {
		common.ErrorResponse(w, err.Error(), statusCode, "ProfileChangeApi")
		return
	}

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(res)
}
