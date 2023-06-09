package delivery

import (
	"fmt"
	"net/http"
	"time"

	"github.com/fahmiaz411/bookla/helper/constant"
	"github.com/fahmiaz411/bookla/helper/encrypt"
	"github.com/fahmiaz411/bookla/helper/field"
	"github.com/fahmiaz411/bookla/helper/message"
	"github.com/fahmiaz411/bookla/helper/web"
	"github.com/fahmiaz411/bookla/middleware"
	"github.com/fahmiaz411/bookla/modules/auth/v1/domain"
	"github.com/fahmiaz411/bookla/modules/auth/v1/interfaces"
	"github.com/gofiber/fiber/v2"
)

type RESTHandler struct {
	Usecase interfaces.AuthUsecase
}

func NewRESTHandler(f fiber.Router, usecase interfaces.AuthUsecase) {
	handler := &RESTHandler{
		Usecase: usecase,
	}

	r := f.Group("/v1")

	r.Post("/auth/check-registered-phone", handler.CheckRegisteredPhone)
	r.Post("/auth/validate-otp", handler.ValidateOTP)
	r.Post("/auth/check-available-username", handler.CheckAvailableUsername)

	r.Post("/auth/login", handler.Login)
	r.Post("/auth/register", handler.Register)
}

func (h *RESTHandler) CheckRegisteredPhone(c *fiber.Ctx) error {
	var lang string = c.Get(constant.HeaderLanguage, "en")	

	req := domain.CheckRegisteredPhoneRequest{}
	c.BodyParser(&req)

	if req.Phone == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.Phone),
			},
		})
	}

	res, err := h.Usecase.CheckRegisteredPhone(c, req)
	if err != nil {
		return nil
	}

	return c.JSON(web.BaseResponse{
		Success: true,
		Meta: web.BaseMeta{
			Code:    http.StatusOK,
			Message: message.Local[lang].Success,
		},
		Data: res,
	})
}

func (h *RESTHandler) ValidateOTP(c *fiber.Ctx) error {
	var lang string = c.Get(constant.HeaderLanguage, "en")	

	req := domain.ValidateOTPRequest{}
	c.BodyParser(&req)

	serverOtp := domain.ServerOTP{}
	encrypt.DeToken(req.EncryptedServerOTP, &serverOtp)

	if serverOtp.ExpiredAt.Unix() < time.Now().UTC().Unix() {
		fmt.Println(serverOtp.ExpiredAt, time.Now())
		return c.Status(http.StatusGone).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusGone,
				Message: message.Local[lang].OTPExpired,
			},
		})
	}

	if !encrypt.Compare(serverOtp.OTP, req.ClientOTP) {
		fmt.Println(serverOtp.OTP)
		return c.Status(http.StatusAccepted).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusAccepted,
				Message: message.Local[lang].OTPNotMatch,
			},
		})
	}

	sessionExpiredAt := time.Now().UTC().Add(1 * time.Hour)
	sessionToken, _ := encrypt.Token(middleware.SessionToken{
		Phone: serverOtp.Phone,
		Step: middleware.SessionStepRegister,
		ExpiredAt: sessionExpiredAt,
	})

	return c.JSON(web.BaseResponse{
		Success: true,
		Meta: web.BaseMeta{
			Code:    http.StatusOK,
			Message: message.Local[lang].Success,
		},
		Data: domain.ValidateOTPResponse{
			Phone: serverOtp.Phone,
			SessionToken: sessionToken,
			ExpiredAt: sessionExpiredAt,
		},
	})
}

func (h *RESTHandler) CheckAvailableUsername(c *fiber.Ctx) error {
	var lang string = c.Get(constant.HeaderLanguage, "en")	

	req := domain.CheckAvailableUsernameRequest{}
	c.BodyParser(&req)

	if req.Username == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.Username),
			},
		})
	}

	res, err := h.Usecase.CheckAvailableUsername(c, req)
	if err != nil {
		return nil
	}

	return c.JSON(web.BaseResponse{
		Success: true,
		Meta: web.BaseMeta{
			Code:    http.StatusOK,
			Message: message.Local[lang].Success,
		},
		Data: res,
	})
}

func (h *RESTHandler) Login(c *fiber.Ctx) error {
	var lang string = c.Get(constant.HeaderLanguage, "en")	

	req := domain.LoginRequest{}
	c.BodyParser(&req)

	if req.SessionToken == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.SessionToken),
			},
		})
	} else if req.PIN == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.PIN),
			},
		})
	}

	sessionBody := middleware.SessionToken{}
	err := encrypt.DeToken(req.SessionToken, &sessionBody)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusUnauthorized,
				Message: message.Local[lang].InvalidSessionToken,
			},
		})
	}

	var isMatchPin bool
	isMatchPin, err = h.Usecase.ValidatePIN(c, sessionBody.Phone, req.PIN)
	if err != nil {
		return nil
	}

	if !isMatchPin {
		return c.Status(http.StatusUnauthorized).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusUnauthorized,
				Message: message.Local[lang].PinNotMatch,
			},
		})
	}

	res, err := h.Usecase.Login(c, sessionBody.Phone)
	if err != nil {
		return nil
	}

	return c.JSON(web.BaseResponse{
		Success: true,
		Meta: web.BaseMeta{
			Code:    http.StatusOK,
			Message: message.Local[lang].Success,
		},
		Data: res,
	})
}

func (h *RESTHandler) Register(c *fiber.Ctx) error {
	var lang string = c.Get(constant.HeaderLanguage, "en")	
	
	req := domain.RegisterRequest{}
	c.BodyParser(&req)

	if req.SessionToken == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.SessionToken),
			},
		})
	} else if req.Fullname == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.Fullname),
			},
		})
	} else if req.PIN == constant.EmptyString {
		return c.Status(http.StatusBadRequest).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusBadRequest,
				Message: message.LocalFunc[lang].Required(field.PIN),
			},
		})
	}

	sessionBody := middleware.SessionToken{}
	err := encrypt.DeToken(req.SessionToken, &sessionBody)
	if err != nil {
		return c.Status(http.StatusUnauthorized).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusUnauthorized,
				Message: message.Local[lang].InvalidSessionToken,
			},
		})
	} else if sessionBody.ExpiredAt.Unix() < time.Now().UTC().Unix() {
		return c.Status(http.StatusUnauthorized).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusUnauthorized,
				Message: message.Local[lang].SessionTokenExpired,
			},
		})
	}

	req.Phone = sessionBody.Phone

	if req.Username != constant.EmptyString {
		checkUsername, err := h.Usecase.CheckAvailableUsername(c, domain.CheckAvailableUsernameRequest{
			Username: req.Username,
		})
		if err != nil {
			return c.Status(http.StatusInternalServerError).JSON(web.BaseResponse{
				Success: false,
				Meta: web.BaseMeta{
					Code:    http.StatusInternalServerError,
					Message: err.Error(),
				},
			})
		}
		if !checkUsername.Available {
			return c.Status(http.StatusConflict).JSON(web.BaseResponse{
				Success: false,
				Meta: web.BaseMeta{
					Code:    http.StatusConflict,
					Message: message.LocalFunc[lang].AlreadyExist(field.Username, req.Username),
				},
			})
		}
	}

	res, err := h.Usecase.Register(c, req)
	if err != nil {
		return nil
	}

	return c.JSON(web.BaseResponse{
		Success: true,
		Meta: web.BaseMeta{
			Code:    http.StatusOK,
			Message: message.Local[lang].Success,
		},
		Data: res,
	})
}