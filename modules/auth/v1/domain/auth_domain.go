package domain

import "time"

// Check Registered Phone

type CheckRegisteredPhoneRequest struct {
	Phone string `json:"phone"`
}

type PhoneRegisteredInfo struct {
	ID int64 `json:"id"`
	StepRegister int `json:"step_register"`
	SessionToken string `json:"session_token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type PhoneUnregisteredInfo struct {
	EncryptedServerOTP string `json:"encrypted_server_otp"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CheckRegisteredPhoneResponse struct {
	Phone string `json:"phone"`
	Registered bool `json:"registered"`
	Info any `json:"info"`
}

type ServerOTP struct {
	Phone string `json:"phone"`
	ExpiredAt time.Time `json:"expired_at"`
	OTP string `json:"otp"`
}

// Registered User

type LoginRequest struct {
	SessionToken string `json:"session_token"`
	PIN string `json:"pin"`
}

type LoginResponse struct {
	ID int64 `json:"id"`
	Phone string `json:"phone"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	StepRegister int `json:"step_register"`
	AccessToken string `json:"access_token"`
}

// Unregistered User

type ValidateOTPRequest struct {
	EncryptedServerOTP string `json:"encrypted_server_otp"`
	ClientOTP string `json:"client_otp"`
}

type ValidateOTPResponse struct {
	Phone string `json:"phone"`
	SessionToken string `json:"session_token"`
	ExpiredAt time.Time `json:"expired_at"`
}

type CheckAvailableUsernameRequest struct {
	Username string `json:"username"`
}

type CheckAvailableUsernameResponse struct {
	Username string `json:"username"`
	Available bool `json:"available"`
}

type RegisterRequest struct {
	SessionToken string `json:"session_token"`
	Phone string `json:"-"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	PIN string `json:"pin"`
}

type RegisterResponse struct {
	ID int64 `json:"id"`
	Phone string `json:"phone"`
	Username string `json:"username"`
	Fullname string `json:"fullname"`
	StepRegister int `json:"step_register"`
	AccessToken string `json:"access_token"`
}