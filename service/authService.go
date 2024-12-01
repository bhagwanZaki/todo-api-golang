package service

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
	"todoGoApi/common"
	"todoGoApi/db"
	"todoGoApi/types"
)

var (
	PASSWORD_CHANGE = 0
	FORGET_PASSWORD = 1
	REGISTER        = 2
	DELETE_ACCOUNT  = 3
)

func getUserDetail(email_username string) (types.User, int, error) {
	var userDetail types.User
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from get_user_detail($1)", email_username).Scan(&userDetail)
	if dbErr != nil {
		common.Logger(dbErr.Error(), "getUserDetail")
		return types.User{}, http.StatusInternalServerError, errors.New("something went wrong")
	}

	return userDetail, http.StatusOK, nil
}

func checkIfUserExistByEmailOrUsername(email_username string) (bool, int, error) {
	var userCount int
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from check_if_user_exist_by_email_or_username($1)", email_username).Scan(&userCount)
	if dbErr != nil {
		common.Logger(dbErr.Error(), "checkIfUserExistByEmailOrUsername")
		return false, http.StatusInternalServerError, errors.New("something went wrong")
	}

	if userCount == 0 {
		return false, http.StatusBadRequest, nil
	}
	return true, http.StatusOK, nil
}

func checkIfUserExist(username string, email string) (bool, int, error) {
	var userCount int
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from check_if_user_exist($1,$2)", username, email).Scan(&userCount)
	if dbErr != nil {
		common.Logger(dbErr.Error(), "checkIfUserExist")
		return true, http.StatusInternalServerError, errors.New("something went wrong")
	}

	if userCount == 0 {
		return false, http.StatusOK, nil
	}
	return true, http.StatusOK, nil
}

func verifyUserPassword(loginData types.LoginSchema) (types.User, int, error) {
	var userData types.User
	var userPassword string

	dbErr := db.DB_CONN.QueryRow(
		context.Background(),
		"select * from login_user($1)", loginData.Username).Scan(
		&userData.Id,
		&userData.Username,
		&userData.Email,
		&userPassword,
		&userData.Fullname,
	)

	if dbErr != nil {
		common.Logger(dbErr.Error(), "verifyUserPassword")
		if dbErr.Error() == "no rows in result set" {
			return types.User{}, http.StatusBadRequest, errors.New("username not found")
		}
		return types.User{}, http.StatusInternalServerError, errors.New("something went wrong")
	}

	passwordMatched := common.CheckPasswordHash(loginData.Password, userPassword)

	if !passwordMatched {
		return types.User{}, http.StatusBadRequest, errors.New("invalid password")
	}

	return userData, http.StatusOK, nil
}

func changeUserPasswordService(email string, password string, request_type int) (bool, int, error) {
	currentTime := time.Now()
	var isRequestExist bool
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	checkRequestErr := db.DB_CONN.QueryRow(
		context.Background(),
		"call check_user_request_exist($1,$2,$3,$4)",
		email,
		request_type,
		dbDate,
		isRequestExist).Scan(&isRequestExist)

	if checkRequestErr != nil {
		common.Logger(checkRequestErr.Error(), "changeUserPasswordService")
		return false, http.StatusInternalServerError, errors.New("something went wrong")
	}
	if !isRequestExist {
		return false, http.StatusBadRequest, errors.New("no request found")
	}
	hashedPassword, hashErr := common.HashPassword(password)

	if hashErr != nil {
		common.Logger(hashErr.Error(), "changeUserPasswordService")
		return false, http.StatusInternalServerError, errors.New("internal server error")
	}

	// change password

	_, updatePasswordErr := db.DB_CONN.Exec(context.Background(), "call update_password($1, $2)", email, hashedPassword)

	if updatePasswordErr != nil {
		common.Logger(updatePasswordErr.Error(), "changeUserPasswordService")
		return false, http.StatusInternalServerError, errors.New("internal server error")
	}

	return true, http.StatusInternalServerError, nil
}

// PUBLIC FUNCTIONS
func LogoutService(userId int, token string) (bool, int, error) {

	_, dbErr := db.DB_CONN.Exec(context.Background(), "call logout($1,$2)", userId, token)

	if dbErr != nil {
		common.Logger(dbErr.Error(), "LogoutService")
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return false, http.StatusBadRequest, errors.New("invalid id")
		}
		return false, http.StatusInternalServerError, errors.New("something went wrong")
	}
	return true, http.StatusNoContent, nil
}

func LogoutFromAllDeviceService(userId int) (bool, int, error) {

	_, dbErr := db.DB_CONN.Exec(context.Background(), "call logout_from_all_device($1)", userId)

	if dbErr != nil {
		common.Logger(dbErr.Error(), "LogoutFromAllDeviceService")
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return false, http.StatusBadRequest, errors.New("invalid id")
		}
		return false, http.StatusInternalServerError, errors.New("something went wrong")
	}
	return true, http.StatusNoContent, nil
}

func LoginService(loginData types.LoginSchema) (types.User, string, int, error) {
	userData, statusCode, verifyPasswordErr := verifyUserPassword(loginData)

	if verifyPasswordErr != nil {
		return types.User{}, "", statusCode, verifyPasswordErr
	}

	token, tokenErr := common.SaveTokenInDb(userData.Id)
	if tokenErr != nil {
		return types.User{}, "", http.StatusInternalServerError, tokenErr
	}

	return userData, token, http.StatusCreated, nil
}

func RegisterService(registerData types.RegisterSchema) (types.UserRegisterStruct, int, error) {
	isUserExist, statusCode, isUserExistErr := checkIfUserExist(registerData.Username, registerData.Email)

	if isUserExistErr != nil {
		return types.UserRegisterStruct{}, statusCode, isUserExistErr
	}

	if isUserExist {
		return types.UserRegisterStruct{}, http.StatusBadRequest, errors.New("username or email already exist")
	}

	existingOTP, otpErr := common.CheckIfOtpAlreadyExist(registerData.Email, REGISTER)

	if otpErr != nil {
		return types.UserRegisterStruct{}, http.StatusInternalServerError, errors.New("something went wrong")
	}
	doesOTPExist := (types.UserOTPDbStruct{}) != existingOTP
	var otp int

	if doesOTPExist {
		otp = existingOTP.Otp
	} else {
		var dbErr error
		otp, dbErr = common.CreateAndSaveOTP(registerData.Email, REGISTER)
		if dbErr != nil {
			return types.UserRegisterStruct{}, http.StatusInternalServerError, errors.New("something went wrong")
		}
	}

	err := common.SendEmail(registerData.Email, "Email Verification", "Email Verifucation\nOTP : "+strconv.Itoa(otp))

	if err != nil {
		common.Logger("SMTP ERROR : "+err.Error(), "RegisterService")
	}

	res := types.UserRegisterStruct{
		Username: registerData.Username,
		Email:    registerData.Email,
		Otp:      otp,
	}

	return res, http.StatusCreated, nil
}

func VerifyOtpAndCompleteRegistration(data types.VerifyOtpAndRegisterSchema) (types.User, string, int, error) {
	var dbData types.UserOTPDbStruct

	dbErr := db.DB_CONN.QueryRow(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		data.Email, data.Otp, REGISTER).Scan(&dbData.Id, &dbData.Email, &dbData.Otp)

	if dbErr != nil {
		common.Logger(dbErr.Error(), "VerifyOtpAndCompleteRegistration")
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return types.User{}, "", http.StatusBadRequest, errors.New("invalid OTP")
		}
		return types.User{}, "", http.StatusInternalServerError, errors.New("something went wrong")
	}

	var userId int
	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	password, hashErr := common.HashPassword(data.Password)

	if hashErr != nil {
		common.Logger(hashErr.Error(), "VerifyOtpAndCompleteRegistration")
		return types.User{}, "", http.StatusInternalServerError, errors.New("something went wrong")
	}

	createErr := db.DB_CONN.QueryRow(
		context.Background(),
		"call create_user($1,$2,$3,$4,$5,$6)",
		data.Username,
		data.Fullname,
		data.Email,
		password,
		dbDate,
		userId,
	).Scan(&userId)

	if createErr != nil {
		common.Logger("CREATE USER DB ERR : "+createErr.Error(), "VerifyOtpAndCompleteRegistration")
		return types.User{}, "", http.StatusInternalServerError, errors.New("unable to create user. Try again")
	}

	token, tokenErr := common.SaveTokenInDb(userId)
	if tokenErr != nil {
		return types.User{}, "", http.StatusInternalServerError, errors.New("something went wrong")
	}

	return types.User{
		Id:       userId,
		Username: data.Username,
		Fullname: data.Fullname,
		Email:    data.Email,
	}, token, http.StatusCreated, nil
}

func ResetPasswordRequest(username string, password string) (int, error) {
	data := types.LoginSchema{
		Username: username,
		Password: password,
	}
	userData, statusCode, verifyPasswordErr := verifyUserPassword(data)

	if verifyPasswordErr != nil {
		return statusCode, verifyPasswordErr
	}

	statusCode,otpErr := CreateOTPRequest(userData.Email, PASSWORD_CHANGE)
 
	if otpErr != nil {
		return statusCode,otpErr
	}

	return http.StatusCreated, nil
}

// This service will send the otp
func CreateOTPRequest(email string, requestTypeCode int) (int, error) {

	log.Println(email, "REQUEST TYPE : ", requestTypeCode)
	userExist, statusCode, userExistErr := checkIfUserExistByEmailOrUsername(email)

	if userExistErr != nil {
		return statusCode, userExistErr
	}

	if !userExist {
		return statusCode, errors.New("email not found")
	}
	
	existingOTP, otpErr := common.CheckIfOtpAlreadyExist(email, requestTypeCode)
	
	if otpErr != nil {
		return http.StatusInternalServerError, errors.New("something went wrong")
	}
	doesOTPExist := (types.UserOTPDbStruct{}) != existingOTP
	var otp int

	if doesOTPExist {
		otp = existingOTP.Otp
	} else {
		var dbErr error
		otp, dbErr = common.CreateAndSaveOTP(email, requestTypeCode)
		if dbErr != nil {
			return http.StatusInternalServerError, errors.New("something went wrong")
		}
	}

	var subject string
	var msgBody string
	switch requestTypeCode {
	case 0:
		subject = "Password Change Request"
		msgBody = "A password reset/change request has been initiated for your account\nTo complete the process, use the following one-time password (OTP):\nOTP : " + strconv.Itoa(otp) + "\nIf you didnt request this change, please disregard this email.\nStay secure!"
	case 1:
		subject = "Forget Password Request"
		msgBody = "A password reset/change request has been initiated for your account\nTo complete the process, use the following one-time password (OTP):\nOTP : " + strconv.Itoa(otp) + "\nIf you didnt request this change, please disregard this email.\nStay secure!"
	case 3:
		subject = "Delete Account Request"
		msgBody = "Account delete request has been initiated from your account.\nIf you did not make this request, please update your password.\n\nTo proceed with the account deletion, please use the following one-time password (OTP):\n\nOTP: " + strconv.Itoa(otp) + "\nIf you have initiated this request, enter the OTP to confirm the account deletion process.\nKeep in mind that once the account is deleted, all data associated with it will be permanently removed."
	default:
		return http.StatusBadRequest,errors.New("invalid request")
	}

	// TODO: MOVE THIS INTO RABBITMQ
	err := common.SendEmail(email, subject, msgBody)

	if err != nil {
		common.Logger("SMTP ERROR : "+err.Error(), "CreateOTPRequeset")
	}

	return http.StatusCreated, nil
}

// This will verify the otp and create a request

func VerifyOTP(email string, otp int, requestType int) (int,error) {
	// verify otp
	var dbData types.UserOTPDbStruct
	dbErr := db.DB_CONN.QueryRow(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		email, otp, requestType).Scan(&dbData.Id, &dbData.Email, &dbData.Otp)

	if dbErr != nil {
		common.Logger(dbErr.Error(),"VerifyOTP")
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return http.StatusBadRequest, errors.New("invalid OTP")
		}
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	// create request data
	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	_, createErr := db.DB_CONN.Exec(context.Background(), "CALL create_request($1,$2,$3)", email, requestType, dbDate)
	
	if createErr != nil {
		common.Logger(createErr.Error(),"VerifyOTP")
		return http.StatusInternalServerError, errors.New("something went wrong")
	}

	return http.StatusOK, nil
}

func ForgetPasswordService(email string, password string) (int, error) {
	userDetail, statusCode, userDetailErr := getUserDetail(email)

	if userDetailErr != nil {
		return statusCode, userDetailErr
	}

	_, statusErr, changePasswordErr := changeUserPasswordService(email, password, FORGET_PASSWORD)

	if changePasswordErr != nil {
		return statusErr, changePasswordErr
	}

	_, logoutCode, err := LogoutFromAllDeviceService(userDetail.Id)

	if err != nil {
		return logoutCode, err
	}

	return http.StatusNoContent, nil
}

func ResetPasswordService(userDetail types.User, currentToken string, password string) (int, error) {
	_, statusCode, changePasswordErr := changeUserPasswordService(userDetail.Email, password, PASSWORD_CHANGE)

	if changePasswordErr != nil {
		common.Logger(changePasswordErr.Error(), "ResetPasswordService")
		return statusCode, changePasswordErr
	}

	_, dbErr := db.DB_CONN.Exec(context.Background(), "call logout_from_all_device_except_one($1,$2)", userDetail.Id, currentToken)

	if dbErr != nil {
		common.Logger(dbErr.Error(), "ResetPasswordService")
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return http.StatusBadRequest, errors.New("invalid id")
		}
		return http.StatusInternalServerError, errors.New("something went wrong")
	}
	return http.StatusNoContent, nil
}

// delete account
func DeleteAccountService(userId int, email string, otp int) (bool, int, error) {

	_, dbErr := db.DB_CONN.Exec(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		email, otp, DELETE_ACCOUNT)

	if dbErr != nil {
		common.Logger(dbErr.Error(), "DeleteAccountService")
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return false, http.StatusBadRequest, errors.New("invalid OTP")
		}
		return false, http.StatusInternalServerError, errors.New("something went wrong")
	}

	_, deleteDbErr := db.DB_CONN.Exec(context.Background(), "call delete_user($1)", userId)

	if deleteDbErr != nil {
		common.Logger(deleteDbErr.Error(), "DeleteAccountService")
		if strings.Contains(deleteDbErr.Error(), "Invalid id") {
			return false, http.StatusBadRequest, errors.New("invalid OTP")
		}
		return false, http.StatusInternalServerError, errors.New("something went wrong")
	}

	return true, http.StatusOK, nil
}

func UpdateProfileService(profileData types.ProfileSchema, userData types.User) (types.User, int, error) {
	res := userData

	if profileData.Username == userData.Username {
		return types.User{}, http.StatusBadRequest, errors.New("invalid request same username")
	}

	if profileData.Email == userData.Email {
		return types.User{}, http.StatusBadRequest, errors.New("invalid request same email")
	}

	if profileData.Fullname == userData.Fullname {
		return types.User{}, http.StatusBadRequest, errors.New("invalid request same fullname")
	}

	// /check if ussername or email aleardy exist
	var existStatusCode int

	existErr := db.DB_CONN.QueryRow(
		context.Background(),
		"select * from check_user_exists($1, $2)",
		profileData.Username,
		profileData.Email,
	).Scan(&existStatusCode)

	if existErr != nil {
		common.Logger(existErr.Error(),"UpdateProfileService")
		return types.User{}, http.StatusInternalServerError, errors.New("something went wrong")
	}

	if existStatusCode != 0 {
		var existErrMsg string

		if existStatusCode == 1 {
			existErrMsg = "username already exist"
		} else if existStatusCode == 2 {
			existErrMsg = "email already exist"
		} else if existStatusCode == 3 {
			existErrMsg = "username and email already exist"
		}

		return types.User{}, http.StatusBadRequest, errors.New(existErrMsg)
	}

	if profileData.Email != "" {
		_, err := db.DB_CONN.Exec(context.Background(), "CALL update_email($1, $2)", userData.Id, profileData.Email)
		
		if err != nil {
			common.Logger(err.Error(),"UpdateProfileService")
			if strings.Contains(err.Error(), "Duplicate data error") {
				return types.User{}, http.StatusBadRequest, errors.New("email already exist")
			}
			return types.User{}, http.StatusInternalServerError, errors.New("something went wrong")
		}

		res.Email = profileData.Email
	}

	if profileData.Username != "" {
		_, err := db.DB_CONN.Exec(context.Background(), "CALL update_username($1, $2)", userData.Id, profileData.Username)

		if err != nil {
			common.Logger(err.Error(),"UpdateProfileService")
			if strings.Contains(err.Error(), "Duplicate data error") {
				return types.User{}, http.StatusBadRequest, errors.New("username already exist")
			}
			return types.User{}, http.StatusInternalServerError, errors.New("something went wrong")
		}

		res.Username = profileData.Username
	}

	if profileData.Fullname != "" {
		_, err := db.DB_CONN.Exec(context.Background(), "CALL update_fullname($1, $2)", userData.Id, profileData.Fullname)

		if err != nil {
			common.Logger(err.Error(),"UpdateProfileService")
			return types.User{}, http.StatusInternalServerError, errors.New("something went wrong")
		}

		res.Fullname = profileData.Fullname
	}

	return res, http.StatusCreated, nil
}
