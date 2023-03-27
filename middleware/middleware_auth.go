package middleware

import (
	"time"
)

type SessionToken struct {
	Phone string `json:"phone"`
	Step	  string	`json:"step"`
	ExpiredAt time.Time `json:"expired_at"`
}

type AccessToken struct {
	ID int64 `json:"id"`
	Phone string `json:"phone"`
	Username string `json:"username"`
}

const (
	SessionStepRegister = "register"
	SessionStepValidatePIN = "validate-pin"
)