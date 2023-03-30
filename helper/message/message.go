package message

import (
	"fmt"
	"strings"
)

type Message struct{
	Success string

	InvalidRequestBody string
	InvalidSessionToken string
	SessionTokenExpired string
	InvalidAccessToken string

	OTPExpired string
	OTPNotMatch string
	PinNotMatch string
}

type MessageFunc struct {
	NotFound func(name, property, value string) string
	CanotNull func(property string) string
	InvalidId func(name string) string
	ShoudMatchEnum func(property string, enums []string) string
	Required func(field string) string
	AlreadyExist func(field, value string) string
}

var Local = map[string]Message{
	"en": {
		Success: "Success",
		InvalidRequestBody: "Invalid Request Body",
		InvalidSessionToken: "Invalid Session Token",
		SessionTokenExpired: "Session Token Expired",
		InvalidAccessToken: "Invalid Access Token",
		OTPExpired: "OTP expired",
		OTPNotMatch: "OTP not match",
		PinNotMatch: "PIN not match",
	},
	"id": {
		Success: "Berhasil",
		InvalidRequestBody: "'Request Body' tidak valid",
		InvalidSessionToken: "Token Sesi tidak valid",
		SessionTokenExpired: "Token Sesi kadaluarsa",
		InvalidAccessToken: "Token Akses tidak valid",
		OTPExpired: "OTP kadaluarsa",
		OTPNotMatch: "OTP tidak cocok",
		PinNotMatch: "PIN tidak cocok",
	},
}

var LocalFunc = map[string]MessageFunc{
	"en": {
		NotFound: func(name, property, value string) string {
			return fmt.Sprintf("%s with %s %s Not Found", name, property, value)
		},
		CanotNull: func(property string) string {
			return fmt.Sprintf("%s cannot be null", property)
		},
		InvalidId: func(name string) string {
			return fmt.Sprintf("Invalid %s Id", name)
		},
		ShoudMatchEnum: func(property string, enums []string) string {
			return fmt.Sprintf("Field '%s' should match one of: %s", property, strings.Join(enums, ", "))
		},
		Required: func(field string) string {
			return fmt.Sprintf("Field '%s' required", field)
		},
		AlreadyExist: func(field, value string) string {
			FieldToTitle(field)
		
			return fmt.Sprintf("%s with value '%s' already exist", FieldToTitle(field), value)
		},
	},
	"id": {
		NotFound: func(name, property, value string) string {
			return fmt.Sprintf("%s with %s %s Not Found", name, property, value)
		},
		CanotNull: func(property string) string {
			return fmt.Sprintf("%s tidak boleh kosong", property)
		},
		InvalidId: func(name string) string {
			return fmt.Sprintf("Id %s tidak valid", name)
		},
		ShoudMatchEnum: func(property string, enums []string) string {
			return fmt.Sprintf("Kolom '%s' harus berisi: %s", property, strings.Join(enums, ", "))
		},
		Required: func(field string) string {
			return fmt.Sprintf("Kolom '%s' diperlukan", field)
		},
		AlreadyExist: func(field, value string) string {
			FieldToTitle(field)
		
			return fmt.Sprintf("%s dengan data '%s' sudah ada", FieldToTitle(field), value)
		},
	},
}

func FieldToTitle(field string) string {
	var formattedField string
	for _, f := range strings.Split(field, "_") {
		formattedField += strings.Title(f) + " "
	}

	return strings.TrimSpace(formattedField)
}