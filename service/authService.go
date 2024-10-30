package service

import (
	"context"
	"errors"
	"fmt"
	"log"
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

func checkIfUserExistByEmailOrUsername(email_username string) (bool, error) {
	var userCount int
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from check_if_user_exist_by_email_or_username($1)", email_username).Scan(&userCount)
	if dbErr != nil {
		return true, dbErr
	}

	if userCount == 0 {
		return false, nil
	}
	return true, nil
}

func checkIfUserExist(username string, email string) (bool, error) {
	var userCount int
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from check_if_user_exist($1,$2)", username, email).Scan(&userCount)
	if dbErr != nil {
		return true, dbErr
	}

	if userCount == 0 {
		return false, nil
	}
	return true, nil
}

// PUBLIC FUNCTIONS
func LogoutService(userId int, token string) (bool, error) {

	_, dbErr := db.DB_CONN.Exec(context.Background(), "call logout($1,$2)", userId, token)

	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return false, errors.New("invalid id")
		}
		return false, dbErr
	}
	return true, nil
}

func LogoutFromAllDeviceService(userId int) (bool, error) {

	_, dbErr := db.DB_CONN.Exec(context.Background(), "call logout_from_all_device($1)", userId)

	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return false, errors.New("invalid id")
		}
		return false, dbErr
	}
	return true, nil
}

func LoginService(loginData types.LoginSchema) (types.User, string, error) {
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
		if dbErr.Error() == "no rows in result set" {
			return types.User{}, "", errors.New("username not found")
		}
		return types.User{}, "", dbErr
	}

	passwordMatched := common.CheckPasswordHash(loginData.Password, userPassword)

	if !passwordMatched {
		return types.User{}, "", errors.New("invalid password")
	}

	token, tokenErr := common.SaveTokenInDb(userData.Id)
	if tokenErr != nil {
		return types.User{}, "", tokenErr
	}

	return userData, token, nil
}

func RegisterService(registerData types.RegisterSchema) (types.UserRegisterStruct, error) {
	isUserExist, isUserExistErr := checkIfUserExist(registerData.Username, registerData.Email)

	if isUserExistErr != nil {
		log.Println("[checkIfOtpAlreadyExist] ERROR", isUserExistErr)
		return types.UserRegisterStruct{}, isUserExistErr
	}

	if isUserExist {
		return types.UserRegisterStruct{}, errors.New("username or email already exist")
	}

	existingOTP := common.CheckIfOtpAlreadyExist(registerData.Email, REGISTER)
	doesOTPExist := (types.UserOTPDbStruct{}) != existingOTP
	var otp int

	if doesOTPExist {
		otp = existingOTP.Otp
	} else {
		var dbErr error
		otp, dbErr = common.CreateAndSaveOTP(registerData.Email, REGISTER)
		if dbErr != nil {
			fmt.Println("DB ERROR ", dbErr.Error())
			return types.UserRegisterStruct{}, dbErr
		}
	}

	err := common.SendEmail(registerData.Email, "Email Verification", "Email Verifucation\nOTP : "+strconv.Itoa(otp))

	if err != nil {
		fmt.Println("SMTP Error", err.Error())
	}

	res := types.UserRegisterStruct{
		Username: registerData.Username,
		Email:    registerData.Email,
		Otp:      otp,
	}

	return res, nil
}

func VerifyOtpAndCompleteRegistration(data types.VerifyOtpAndRegisterSchema) (types.User, string, error) {
	var dbData types.UserOTPDbStruct

	dbErr := db.DB_CONN.QueryRow(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		data.Email, data.Otp, REGISTER).Scan(&dbData.Id, &dbData.Email, &dbData.Otp)

	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			return types.User{}, "", errors.New("invalid OTP")
		}
		log.Println("VERIFY DB ERR ", dbErr)
		return types.User{}, "", dbErr
	}

	var userId int
	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	password, hashErr := common.HashPassword(data.Password)

	if hashErr != nil {
		log.Println("HASH ERR ", hashErr)
		return types.User{}, "", errors.New("something went wrong")
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
		log.Println("CREATE USER DB ERR ", createErr)
		return types.User{}, "", createErr
	}

	token, tokenErr := common.SaveTokenInDb(userId)
	if tokenErr != nil {
		return types.User{}, "", tokenErr
	}

	return types.User{
		Id:       userId,
		Username: data.Username,
		Fullname: data.Fullname,
		Email:    data.Email,
	}, token, nil
}

// This service will send the otp
func CreateOTPRequest(email string, requestTypeCode int) (bool, error) {

	userExist, userExistErr := checkIfUserExistByEmailOrUsername(email)

	if userExistErr != nil {
		return false, errors.New("invalid email")
	}
	if !userExist {
		return false, errors.New("email not found")
	}

	existingOTP := common.CheckIfOtpAlreadyExist(email, REGISTER)
	doesOTPExist := (types.UserOTPDbStruct{}) != existingOTP
	var otp int

	if doesOTPExist {
		otp = existingOTP.Otp
	} else {
		var dbErr error
		otp, dbErr = common.CreateAndSaveOTP(email, REGISTER)
		if dbErr != nil {
			fmt.Println("DB ERROR ", dbErr.Error())
			return false, dbErr
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
		return false, errors.New("invalid request")
	}

	// TODO: MOVE THIS INTO RABBITMQ
	err := common.SendEmail(email, subject, msgBody)

	if err != nil {
		log.Println("SMTP Error", err.Error())
	}

	return true, nil
}

// This will verify the otp and create a request

func VerifyPasswordChangeOTP(email string, otp int, is_forget_password_request bool) (bool, error) {
	// verify otp
	var dbData types.UserOTPDbStruct

	requestType := PASSWORD_CHANGE

	if is_forget_password_request {
		requestType = FORGET_PASSWORD
	}
	dbErr := db.DB_CONN.QueryRow(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		email, otp, requestType).Scan(&dbData.Id, &dbData.Email, &dbData.Otp)

	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			return false, errors.New("invalid OTP")
		}
		return false, dbErr
	}

	// delete user otp

	_, err := db.DB_CONN.Exec(context.Background(), "CALL delete_user_top($1,$2)", dbData.Id, email)

	if err != nil {
		return false, err
	}

	// create request data
	dbRequestType := "password_change"

	if is_forget_password_request {
		dbRequestType = "forget_password"
	}

	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	_, createErr := db.DB_CONN.Exec(context.Background(), "CALL create_request($1,$2,$3)", email, dbRequestType, dbDate)

	if createErr != nil {
		return false, createErr
	}

	return true, nil
}

// This will change the password
func ChangeUserPasswordService(email string, password string, is_forget_password_request bool) (bool, error) {
	// check for request
	var isRequestExist bool

	requestType := PASSWORD_CHANGE

	if is_forget_password_request {
		requestType = FORGET_PASSWORD
	}

	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	checkRequestErr := db.DB_CONN.QueryRow(
		context.Background(),
		"select * from check_user_request_exist($1,$2,$3)",
		email,
		requestType,
		dbDate,
	).Scan(&isRequestExist)

	if checkRequestErr != nil {
		return false, checkRequestErr
	}

	hashedPassword, hashErr := common.HashPassword(password)

	if hashErr != nil {
		if strings.Contains(hashErr.Error(), "Invalid id") {
			return false, errors.New("invalid id")
		}
		return false, hashErr
	}

	// change password

	_, updatePasswordErr := db.DB_CONN.Exec(context.Background(), "call update_password($1, $2)", email, hashedPassword)

	if updatePasswordErr != nil {
		return false, updatePasswordErr
	}

	return true, nil
}

// delete account
func DeleteAccountService(userId int, email string) (bool, error) {
	var isRequestExist bool

	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	checkRequestErr := db.DB_CONN.QueryRow(
		context.Background(),
		"select * from check_user_request_exist($1,$2,$3)",
		email,
		DELETE_ACCOUNT,
		dbDate,
	).Scan(&isRequestExist)

	if checkRequestErr != nil {
		return false, checkRequestErr
	}

	_, dbErr := db.DB_CONN.Exec(context.Background(), "call delete_user($1)", userId)

	if dbErr != nil {
		if strings.Contains(dbErr.Error(), "Invalid id") {
			return false, errors.New("invalid id")
		}
		return false, dbErr
	}

	return true, nil
}
