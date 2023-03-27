package usecase

import (
	"context"
	"net/http"
	"time"

	"github.com/fahmiaz411/bookla/helper/web"
	"github.com/fahmiaz411/bookla/modules/auth/v1/domain"
	"github.com/fahmiaz411/bookla/modules/auth/v1/interfaces"
	"github.com/fahmiaz411/bookla/modules/auth/v1/repository"
	"github.com/gofiber/fiber/v2"
)

type Usecase struct {
	repo           *repository.Repository
	contentTimeout time.Duration
}

func NewUsecase(repo *repository.Repository, timeout time.Duration) interfaces.AuthUsecase {
	return &Usecase{
		repo:           repo,
		contentTimeout: timeout,
	}
}

func (u *Usecase) CheckRegisteredPhone(c *fiber.Ctx, req domain.CheckRegisteredPhoneRequest) (res domain.CheckRegisteredPhoneResponse, err error) {
	ctx, cancel := context.WithTimeout(c.UserContext(), u.contentTimeout)
	defer cancel()

	res, err = u.repo.MySQL.CheckRegisteredPhone(ctx, req)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return
}

func (u *Usecase) CheckAvailableUsername(c *fiber.Ctx, req domain.CheckAvailableUsernameRequest) (res domain.CheckAvailableUsernameResponse, err error) {
	ctx, cancel := context.WithTimeout(c.UserContext(), u.contentTimeout)
	defer cancel()

	res, err = u.repo.MySQL.CheckAvailableUsername(ctx, req)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return
}

func (u *Usecase) ValidatePIN(c *fiber.Ctx, phone, pin string) (match bool, err error) {
	ctx, cancel := context.WithTimeout(c.UserContext(), u.contentTimeout)
	defer cancel()

	match, err = u.repo.MySQL.ValidatePIN(ctx, phone, pin)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return
}

func (u *Usecase) Login(c *fiber.Ctx, phone string) (res domain.LoginResponse, err error) {
	ctx, cancel := context.WithTimeout(c.UserContext(), u.contentTimeout)
	defer cancel()

	res, err = u.repo.MySQL.Login(ctx, phone)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return
}

func (u *Usecase) Register(c *fiber.Ctx, req domain.RegisterRequest) (res domain.RegisterResponse, err error) {
	ctx, cancel := context.WithTimeout(c.UserContext(), u.contentTimeout)
	defer cancel()

	res, err = u.repo.MySQL.Register(ctx, req)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(web.BaseResponse{
			Success: false,
			Meta: web.BaseMeta{
				Code:    http.StatusInternalServerError,
				Message: err.Error(),
			},
		})
	}

	return
}