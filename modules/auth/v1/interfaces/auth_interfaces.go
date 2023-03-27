package interfaces

import (
	"context"

	"github.com/fahmiaz411/bookla/modules/auth/v1/domain"
	"github.com/gofiber/fiber/v2"
)

type AuthUsecase interface {
	CheckRegisteredPhone(c *fiber.Ctx, req domain.CheckRegisteredPhoneRequest) (res domain.CheckRegisteredPhoneResponse, err error)
	CheckAvailableUsername(c *fiber.Ctx, req domain.CheckAvailableUsernameRequest) (res domain.CheckAvailableUsernameResponse, err error)
	ValidatePIN(c *fiber.Ctx, phone, pin string) (match bool, err error)
	Login(c *fiber.Ctx, phone string) (res domain.LoginResponse, err error)
	Register(c *fiber.Ctx, req domain.RegisterRequest) (res domain.RegisterResponse, err error)
}

type AuthRepoMysql interface {
	CheckRegisteredPhone(ctx context.Context, req domain.CheckRegisteredPhoneRequest) (res domain.CheckRegisteredPhoneResponse, err error)
	CheckAvailableUsername(ctx context.Context, req domain.CheckAvailableUsernameRequest) (res domain.CheckAvailableUsernameResponse, err error)
	ValidatePIN(ctx context.Context, phone, pin string) (match bool, err error)
	Login(ctx context.Context, phone string) (res domain.LoginResponse, err error)
	Register(ctx context.Context, req domain.RegisterRequest) (res domain.RegisterResponse, err error)
}