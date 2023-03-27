package mysql

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
	"time"

	"github.com/fahmiaz411/bookla/helper/constant"
	"github.com/fahmiaz411/bookla/helper/encrypt"
	"github.com/fahmiaz411/bookla/middleware"
	"github.com/fahmiaz411/bookla/modules/auth/v1/domain"
	"github.com/fahmiaz411/bookla/modules/auth/v1/interfaces"
)

type MysqlRepository struct {
	Conn *sql.DB
}

func NewMysqlRepository(Conn *sql.DB) interfaces.AuthRepoMysql {
	return &MysqlRepository{
		Conn: Conn,
	}
}

func (m *MysqlRepository) CheckRegisteredPhone(ctx context.Context, req domain.CheckRegisteredPhoneRequest) (res domain.CheckRegisteredPhoneResponse, err error) {
	var stmt *sql.Stmt
	stmt, err = m.Conn.PrepareContext(ctx, `
		SELECT
			id,
			step_register
		FROM users
		WHERE phone = ?
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	req.Phone = strings.Replace(req.Phone, "+", "", 1)

	if req.Phone[:1] == "0" {
		req.Phone = "62" + req.Phone[1:]
	} else if req.Phone[:2] != "62" {
		req.Phone = "62" + req.Phone
	}

	var rows *sql.Rows
	rows, err = stmt.QueryContext(ctx, req.Phone)
	if err != nil {
		return
	}

	var foundId int64
	var stepRegister int

	for rows.Next() {
		if err = rows.Scan(
			&foundId,
			&stepRegister,
		); err != nil {
			return
		}
	}

	res.Phone = req.Phone

	sessionExpiredAt := time.Now().UTC().Add(1 * time.Hour)
	sessionToken, _ := encrypt.Token(middleware.SessionToken{
		Phone: req.Phone,
		Step: middleware.SessionStepValidatePIN,
		ExpiredAt: sessionExpiredAt,
	})

	otp := encrypt.RandInt(4)
	otpExpiredAt := time.Now().UTC().Add(5 * time.Minute)
	encOTP, _ := encrypt.Token(domain.ServerOTP{
		Phone: req.Phone,
		ExpiredAt: otpExpiredAt,
		OTP: encrypt.Hash(otp),
	})
	fmt.Println(otp)

	if foundId != int64(constant.ZeroValue) {
		res.Registered = true
		res.Info = domain.PhoneRegisteredInfo{
			ID: foundId,
			StepRegister: stepRegister,
			SessionToken: sessionToken,
			ExpiredAt: sessionExpiredAt,
		}
	} else {
		res.Info = domain.PhoneUnregisteredInfo{
			EncryptedServerOTP: encOTP,
			ExpiredAt: otpExpiredAt,
		}
	}

	return
}

func (m *MysqlRepository) CheckAvailableUsername(ctx context.Context, req domain.CheckAvailableUsernameRequest) (res domain.CheckAvailableUsernameResponse, err error) {
	var stmt *sql.Stmt
	stmt, err = m.Conn.PrepareContext(ctx, `
		SELECT
			id
		FROM users
		WHERE username = ?
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.QueryContext(ctx, req.Username)
	if err != nil {
		return
	}

	var foundId int64

	for rows.Next() {
		if err = rows.Scan(
			&foundId,
		); err != nil {
			return
		}
	}

	res.Username = req.Username

	if foundId == int64(constant.ZeroValue) {
		res.Available = true
	}

	return
}

func (m *MysqlRepository) ValidatePIN(ctx context.Context, phone, pin string) (match bool, err error) {
	var stmt *sql.Stmt
	stmt, err = m.Conn.PrepareContext(ctx, `
		SELECT
			pin
		FROM users
		WHERE phone = ?
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.QueryContext(ctx, phone)
	if err != nil {
		return
	}

	var hashedPin string

	if rows.Next() {
		if err = rows.Scan(
			&hashedPin,
		); err != nil {
			return
		}
	}

	if hashedPin == constant.EmptyString {
		err = fmt.Errorf("User not found")
	}

	match = encrypt.Compare(hashedPin, pin)

	return
}

func (m *MysqlRepository) Login(ctx context.Context, phone string) (res domain.LoginResponse, err error) {
	var stmt *sql.Stmt
	stmt, err = m.Conn.PrepareContext(ctx, `
		SELECT
			id,
			COALESCE(username, ''),
			fullname,
			step_register			
		FROM users
		WHERE phone = ?
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	var rows *sql.Rows
	rows, err = stmt.QueryContext(ctx, phone)
	if err != nil {
		return
	}

	if rows.Next() {
		if err = rows.Scan(
			&res.ID,
			&res.Username,
			&res.Fullname,
			&res.StepRegister,
		); err != nil {
			return
		}
	}

	res.Phone = phone
	res.AccessToken, _ = encrypt.Token(middleware.AccessToken{
		ID: res.ID,
		Phone: phone,
		Username: res.Username,
	})

	return
}

func (m *MysqlRepository) Register(ctx context.Context, req domain.RegisterRequest) (res domain.RegisterResponse, err error) {
	var stmt *sql.Stmt
	stmt, err = m.Conn.PrepareContext(ctx, `
		INSERT INTO users (
			phone,
			username,
			fullname,
			pin,
			step_register
		) VALUES (
			?,
			?,
			?,
			?,
			?
		)
	`)
	if err != nil {
		return
	}
	defer stmt.Close()

	var (
		username any
		stepRegister int = 1
	)

	if req.Username != constant.EmptyString {
		username = req.Username
	}

	insValues := []any{
		req.Phone,
		username,
		req.Fullname,
		encrypt.Hash(req.PIN),
		stepRegister,
	}

	var resultIns sql.Result
	resultIns, err = stmt.ExecContext(ctx, insValues...)
	
	lastInsertId, _ := resultIns.LastInsertId()

	res.ID = lastInsertId
	res.Phone = req.Phone
	res.Username = req.Username
	res.Fullname = req.Fullname
	res.StepRegister = stepRegister
	res.AccessToken, _ = encrypt.Token(middleware.AccessToken{
		ID: res.ID,
		Phone: req.Phone,
		Username: res.Username,
	})

	return
}