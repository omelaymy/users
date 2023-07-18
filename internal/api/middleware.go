package api

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/omelaymy/users/internal/auth"
)

type MWManager struct {
	authUsecase auth.Usecase
}

func NewMWManager(
	authUsecase auth.Usecase,
) *MWManager {
	return &MWManager{
		authUsecase: authUsecase,
	}
}

func (mw *MWManager) BasicAuth() fiber.Handler {
	return basicauth.New(
		basicauth.Config{
			Authorizer: mw.authUsecase.Authentication,
		})
}

func (mw *MWManager) AdminAuth() fiber.Handler {
	return basicauth.New(basicauth.Config{
		Authorizer:   mw.authUsecase.AdminAuthorization,
		Unauthorized: Forbidden,
	})
}

func Forbidden(c *fiber.Ctx) error {
	c.Set(fiber.HeaderWWWAuthenticate, "Basic realm=admin")
	return c.SendStatus(fiber.StatusForbidden)
}
