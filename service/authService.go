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

func getUserDetail(email_username string) (types.User, error) {
	var userDetail types.User
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from get_user_detail($1)", email_username).Scan(&userDetail)
	if dbErr != nil {
		return types.User{}, dbErr
	}

	return userDetail, nil
}

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

func verifyUserPassword(loginData types.LoginSchema) (types.User, error) {
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
			return types.User{}, errors.New("username not found")
		}
		return types.User{}, dbErr
	}

	passwordMatched := common.CheckPasswordHash(loginData.Password, userPassword)

	if !passwordMatched {
		return types.User{}, errors.New("invalid password")
	}

	return userData, nil
}

func changeUserPasswordService(email string, password string, request_type int)(bool, error) {
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
		log.Println("?? checkRequestErr")
		return false, checkRequestErr
	}
	
	hashedPassword, hashErr := common.HashPassword(password)
	
	if hashErr != nil {
		if strings.Contains(hashErr.Error(), "Invalid id") {
			log.Println("?? hashErr")
			return false, errors.New("invalid id")
		}
		return false, hashErr
	}
	
	// change password
	
	_, updatePasswordErr := db.DB_CONN.Exec(context.Background(), "call update_password($1, $2)", email, hashedPassword)
	
	if updatePasswordErr != nil {
		log.Println("?? updatePasswordErr")
		return false, updatePasswordErr
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
	userData, loginErr := verifyUserPassword(loginData)

	if loginErr != nil {
		return types.User{}, "", loginErr
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

func ResetPasswordRequest(username string, password string) (bool, error) {
	data := types.LoginSchema{
		Username: username,
		Password: password,
	}
	userData, verifyPasswordErr := verifyUserPassword(data)

	if verifyPasswordErr != nil {
		return false, verifyPasswordErr
	}

	_, otpErr := CreateOTPRequest(userData.Email, PASSWORD_CHANGE)

	if otpErr != nil {
		return false, otpErr
	}

	return true, nil
}

// This service will send the otp
func CreateOTPRequest(email string, requestTypeCode int) (bool, error) {

	log.Println(email, "REQUEST TYPE : ", requestTypeCode)
	userExist, userExistErr := checkIfUserExistByEmailOrUsername(email)

	if userExistErr != nil {
		return false, errors.New("invalid email")
	}
	if !userExist {
		return false, errors.New("email not found")
	}

	existingOTP := common.CheckIfOtpAlreadyExist(email, requestTypeCode)
	doesOTPExist := (types.UserOTPDbStruct{}) != existingOTP
	var otp int

	if doesOTPExist {
		otp = existingOTP.Otp
	} else {
		var dbErr error
		otp, dbErr = common.CreateAndSaveOTP(email, requestTypeCode)
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

func VerifyOTP(email string, otp int, requestType int) (bool, error) {
	// verify otp
	var dbData types.UserOTPDbStruct
	dbErr := db.DB_CONN.QueryRow(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		email, otp, requestType).Scan(&dbData.Id, &dbData.Email, &dbData.Otp)

	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			return false, errors.New("invalid OTP")
		}
		return false, dbErr
	}

	// create request data
	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	_, createErr := db.DB_CONN.Exec(context.Background(), "CALL create_request($1,$2,$3)", email, requestType, dbDate)

	if createErr != nil {
		return false, createErr
	}

	return true, nil
}

func ForgetPasswordService(email string, password string) (bool, error) {
	userDetail, userDetailErr := getUserDetail(email)

	if userDetailErr != nil {
		return false, userDetailErr
	}

	_, changePasswordErr := changeUserPasswordService(email, password, FORGET_PASSWORD)
	
	if changePasswordErr != nil {
		return false, changePasswordErr
	}

	_, err := LogoutFromAllDeviceService(userDetail.Id)

	if err != nil {
		return false, err
	}

	return true, nil
}

func ResetPasswordService(userDetail types.User,currentToken string, password string) (bool, error)  {
	_, changePasswordErr := changeUserPasswordService(userDetail.Email, password, PASSWORD_CHANGE)

	if changePasswordErr != nil {
		log.Println("?? CHANGE PASSWORD ERR")
		return false, changePasswordErr
	}
	
	_, dbErr := db.DB_CONN.Exec(context.Background(), "call logout_from_all_device_except_one($1,$2)", userDetail.Id, currentToken)
	
	if dbErr != nil {
		log.Println(userDetail.Id, currentToken)
		if strings.Contains(dbErr.Error(), "Invalid id") {
			log.Println("?? LOGOUT ERR")
			return false, errors.New("invalid id")
		}
		return false, dbErr
	}
	return true, nil
}

// delete account
func DeleteAccountService(userId int, email string, otp int) (bool, error) {

	_, dbErr := db.DB_CONN.Exec(
		context.Background(), "select * from verify_user_otp($1,$2,$3)",
		email, otp, DELETE_ACCOUNT)

	if dbErr != nil {
		if dbErr.Error() == "no rows in result set" {
			return false, errors.New("invalid OTP")
		}
		return false, dbErr
	}

	_, deleteDbErr := db.DB_CONN.Exec(context.Background(), "call delete_user($1)", userId)

	if deleteDbErr != nil {
		if strings.Contains(deleteDbErr.Error(), "Invalid id") {
			return false, errors.New("invalid id")
		}
		return false, deleteDbErr
	}

	return true, nil
}

func UpdateProfileService(profileData types.ProfileSchema, userData types.User) (types.User, int, error) {
	res := userData

	if profileData.Username == userData.Username {
		return types.User{}, 400, errors.New("invalid request same username")
	}

	if profileData.Email == userData.Email {
		return types.User{}, 400, errors.New("invalid request same email")
	}

	if profileData.Fullname == userData.Fullname {
		return types.User{}, 400, errors.New("invalid request same fullname")
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
		return types.User{}, 500, existErr
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

		return types.User{}, 400, errors.New(existErrMsg)
	}

	if profileData.Email != "" {
		_, err := db.DB_CONN.Exec(context.Background(), "CALL update_email($1, $2)", userData.Id, profileData.Email)

		if err != nil {
			if strings.Contains(err.Error(), "Duplicate data error") {
				return types.User{}, 400, errors.New("email already exist")
			}
			return types.User{}, 500, err
		}

		res.Email = profileData.Email
	}

	if profileData.Username != "" {
		_, err := db.DB_CONN.Exec(context.Background(), "CALL update_username($1, $2)", userData.Id, profileData.Username)

		if err != nil {
			if strings.Contains(err.Error(), "Duplicate data error") {
				return types.User{}, 400, errors.New("username already exist")
			}
			return types.User{}, 500, err
		}

		res.Username = profileData.Username
	}

	if profileData.Fullname != "" {
		_, err := db.DB_CONN.Exec(context.Background(), "CALL update_fullname($1, $2)", userData.Id, profileData.Fullname)

		if err != nil {
			return types.User{}, 500, err
		}

		res.Fullname = profileData.Fullname
	}

	return res, 201, nil
}
