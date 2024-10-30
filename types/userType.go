package types

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
}

type TokenStruct struct {
	Id     int
	Token  string
	Digest string
}

type UserOTPDbStruct struct {
	Id    int    `json:"id"`
	Email string `json:"email"`
	Otp   int    `json:"otp"`
}

type UserOTPStruct struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Otp      int    `json:"otp"`
}

type UserRegisterStruct struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Otp      int    `json:"otp"`
}

// SCHEMAs

type LoginSchema struct {
	Username string
	Password string
}

type RegisterSchema struct {
	Username string
	Email    string
}

type VerifyOtpAndRegisterSchema struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Fullname string `json:"fullname"`
	Password string `json:"password"`
	Otp      int    `json:"otp"`
}

type OtpRequestSchema struct {
	Email           string `json:"email"`
	RequestTypeCode int    `json:"requestType"`
}

type VerifyOtpSchema struct {
	Email                  string `json:"email"`
	Otp                    int    `json:"otp"`
	IsForgetPasswordReqest bool   `json:"isForgetPassword"`
}

type PasswordChangeSchema struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
