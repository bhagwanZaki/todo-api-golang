package common

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"log"
	random "math/rand/v2"
	"net/http"
	"time"
	"todoGoApi/db"
	"todoGoApi/types"
	"todoGoApi/utils"

	"golang.org/x/crypto/bcrypt"
)

func CheckIfOtpAlreadyExist(email string, requestType int) types.UserOTPDbStruct {
	var existingOTP types.UserOTPDbStruct
	//
	checkErr := db.DB_CONN.QueryRow(
		context.Background(), "select * from check_if_otp_already_exist($1,$2)",
		email, requestType).Scan(&existingOTP.Id, &existingOTP.Email, &existingOTP.Otp)

	if checkErr != nil && checkErr.Error() != "no rows in result set" {
		log.Println("[checkIfOtpAlreadyExist] ERROR", checkErr)
		return types.UserOTPDbStruct{}
	}

	log.Println("existingOTP ", existingOTP)
	return existingOTP
}

func CreateAndSaveOTP(email string, request_type int) (int, error) {
	otp := random.IntN(900000)

	if otp < 100000 {
		otp = otp + 100000
	}

	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	fmt.Println(otp)
	_, dbErr := db.DB_CONN.Exec(
		context.Background(),
		"call create_user_otp($1,$2,$3,$4)",
		email,
		otp,
		request_type,
		dbDate,
	)

	if dbErr != nil {
		return 0, dbErr
	}

	return otp, nil
}

func GetToken(token string) (types.TokenStruct, error) {
	var tokenData types.TokenStruct
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from check_auth($1)", token).Scan(
		&tokenData.Id,
		&tokenData.Token,
		&tokenData.Digest,
	)

	if dbErr != nil {
		return types.TokenStruct{}, dbErr
	}

	return tokenData, nil
}

func generateBytes(length int) ([]byte, error) {
	bytes := make([]byte, length)
	_, err := rand.Read(bytes)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// hashToken hashes the given token using SHA-256.
func hashToken(token string) string {
	hash := sha256.New()
	hash.Write([]byte(token))
	return hex.EncodeToString(hash.Sum(nil))
}

// createToken generates the token, token key, and digest.
func CreateToken() (string, string, string, error) {
	bytes, err := generateBytes(utils.AUTH_TOKEN_LENGTH / 2)
	if err != nil {
		return "", "", "", err
	}
	token := hex.EncodeToString(bytes)
	digest := hashToken(token)
	tokenKey := token[:utils.TOKEN_KEY_LENGTH]

	return token, tokenKey, digest, nil
}

func SaveTokenInDb(id int) (string, error) {
	token, dbToken, digest, tokenErr := CreateToken()

	if tokenErr != nil {
		return "", errors.New("token issue")
	}
	currentTime := time.Now()
	dbDate := fmt.Sprintf("%d-%d-%d", currentTime.Year(), currentTime.Month(), currentTime.Day())

	_, tokenDbErr := db.DB_CONN.Exec(context.Background(), "CALL create_token($1,$2,$3,$4)", id, dbToken, digest, dbDate)

	if tokenDbErr != nil {
		return "", tokenDbErr
	}

	return token, nil
}

func SetTokenCookie(w http.ResponseWriter, token string) {
	cookie := &http.Cookie{
		Name:     "token",
		Value:    token,
		MaxAge:   0,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}

	http.SetCookie(w, cookie)
}

func CheckTokenValidity(r *http.Request) (string, error) {
	cookie, err := r.Cookie("token")

	if err != nil {
		return "", err
	}

	token := cookie.Value

	if len(token) != utils.AUTH_TOKEN_LENGTH {
		return "", errors.New("invalid token")
	}

	dbToken, dbErr := GetToken(token[:utils.TOKEN_KEY_LENGTH])

	if dbErr != nil {
		return "", dbErr
	}

	hashedDigest := hashToken(token)

	if hmac.Equal([]byte(dbToken.Digest), []byte(hashedDigest)) {
		return token[:utils.TOKEN_KEY_LENGTH], nil
	}

	return "", errors.New("invalid token")
}

func GetUserDataFromToken(r *http.Request) (types.User, error) {
	token, err := CheckTokenValidity(r)

	if err != nil {
		return types.User{}, err
	}

	if token == "" {
		return types.User{}, errors.New("invalid token")
	}

	var userData types.User
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from get_user_data_from_token($1)", token).Scan(
		&userData.Id,
		&userData.Username,
		&userData.Email,
		&userData.Fullname,
	)

	if dbErr != nil {
		return types.User{}, dbErr
	}

	return userData, nil

}
func GetUserDataWithTokenFromToken(r *http.Request) (types.User, string, error) {
	token, err := CheckTokenValidity(r)

	if err != nil {
		return types.User{}, "", err
	}

	if token == "" {
		return types.User{}, "", errors.New("invalid token")
	}

	var userData types.User
	dbErr := db.DB_CONN.QueryRow(context.Background(), "select * from get_user_data_from_token($1)", token).Scan(
		&userData.Id,
		&userData.Username,
		&userData.Email,
		&userData.Fullname,
	)

	if dbErr != nil {
		return types.User{}, "", dbErr
	}

	return userData, token, nil
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
