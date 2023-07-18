package di

import (
	"fmt"
	"strings"

	"github.com/go-playground/locales/en"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/omelaymy/users/config"
	"github.com/omelaymy/users/internal/api"
	"github.com/omelaymy/users/internal/api/http/delivery"
	"github.com/omelaymy/users/internal/api/http/errors"
	"github.com/omelaymy/users/pkg/db/inmemory"
	"github.com/omelaymy/users/pkg/flags"
	"github.com/omelaymy/users/pkg/logger"
	"github.com/omelaymy/users/pkg/secure"
	"github.com/rs/zerolog"
	"github.com/samber/do"

	ut "github.com/go-playground/universal-translator"
	authRepo "github.com/omelaymy/users/internal/auth/repository"
	authUsecase "github.com/omelaymy/users/internal/auth/usecase"
	usersRepo "github.com/omelaymy/users/internal/users/repository"
	usersUsecase "github.com/omelaymy/users/internal/users/usecase"
)

func NewFlags(*do.Injector) (*flags.Flags, error) {
	flags, err := flags.New()
	if err != nil {
		return nil, fmt.Errorf("failed to create flags: %w", err)
	}

	return flags, nil
}

func NewConfig(i *do.Injector) (*config.Config, error) {
	allFlags := do.MustInvoke[*flags.Flags](i)

	viperConfig, err := config.LoadConfig(*allFlags.ConfigFile)
	if err != nil {
		return nil, fmt.Errorf("load config error: %w", err)
	}

	cfg, err := config.ParseConfig(viperConfig)
	if err != nil {
		return nil, fmt.Errorf("parse config error: %w", err)
	}

	cfg.BaseAdmin.Password, err = secure.HashPassword(cfg.BaseAdmin.Password)
	if err != nil {
		return nil, fmt.Errorf("hashing admin password error: %w", err)
	}

	return cfg, nil
}

func NewLogger(i *do.Injector) (*zerolog.Logger, error) {
	cfg := do.MustInvoke[*config.Config](i)
	logger := logger.NewLogger(cfg)

	return &logger, nil
}

func NewInMemoryDatabase(i *do.Injector) (*inmemory.InMemoryDatabase, error) {
	cfg := do.MustInvoke[*config.Config](i)

	db := inmemory.NewInMemoryDatabase()
	_, err := db.InsertUser(inmemory.User{
		Email:    cfg.BaseAdmin.Email,
		Username: cfg.BaseAdmin.Username,
		Password: cfg.BaseAdmin.Password,
		Admin:    true,
	})
	if err != nil {
		return nil, err
	}

	return db, nil

}

func NewUsers(i *do.Injector) (*usersUsecase.Users, error) {
	return usersUsecase.NewUsers(
		do.MustInvoke[*config.Config](i),
		do.MustInvoke[*usersRepo.UsersRepository](i),
	), nil
}

func NewUsersRepository(i *do.Injector) (*usersRepo.UsersRepository, error) {
	return usersRepo.NewUsersRepository(
		do.MustInvoke[*inmemory.InMemoryDatabase](i),
		do.MustInvoke[*zerolog.Logger](i),
	), nil
}

func NewAuth(i *do.Injector) (*authUsecase.Auth, error) {
	return authUsecase.NewAuth(
		do.MustInvoke[*authRepo.AuthRepository](i),
	), nil
}

func NewAuthRepository(i *do.Injector) (*authRepo.AuthRepository, error) {
	return authRepo.NewAuthRepository(
		do.MustInvoke[*inmemory.InMemoryDatabase](i),
		do.MustInvoke[*zerolog.Logger](i),
	), nil
}

func NewMWManager(i *do.Injector) (*api.MWManager, error) {
	return api.NewMWManager(
		do.MustInvoke[*authUsecase.Auth](i),
	), nil
}

func NewHandlers(i *do.Injector) (*delivery.Handlers, error) {
	return delivery.NewHandlers(
		do.MustInvoke[*usersUsecase.Users](i),
		do.MustInvoke[*validator.Validate](i),
		do.MustInvoke[ut.Translator](i),
	), nil
}

func NewRoutes(i *do.Injector) (*delivery.Routes, error) {
	return delivery.NewRoutes(
		do.MustInvoke[*delivery.Handlers](i),
		do.MustInvoke[*api.MWManager](i),
		do.MustInvoke[*fiber.App](i),
	), nil
}

func NewHttpErrorHandler(i *do.Injector) (*errors.HttpErrorHandler, error) {
	return errors.NewHttpErrorHandler(
		do.MustInvoke[*zerolog.Logger](i),
	), nil
}

func NewFiberApp(i *do.Injector) (*fiber.App, error) {
	h := do.MustInvoke[*errors.HttpErrorHandler](i)

	app := fiber.New(fiber.Config{ErrorHandler: h.Handler})
	app.Use(cors.New(cors.Config{
		AllowOrigins: "*",
		AllowHeaders: "*",
		AllowMethods: strings.Join([]string{
			fiber.MethodGet,
			fiber.MethodPost,
			fiber.MethodPut,
			fiber.MethodDelete,
		}, ","),
	}))

	return app, nil
}

func NewTranslator(i *do.Injector) (ut.Translator, error) {
	cfg := do.MustInvoke[*config.Config](i)

	enTranslator := en.New()
	uni := ut.New(enTranslator, enTranslator)
	translator, _ := uni.GetTranslator(cfg.System.DefaultLocale)

	return translator, nil
}

func NewValidate(i *do.Injector) (*validator.Validate, error) {
	translator := do.MustInvoke[ut.Translator](i)

	v := validator.New()
	err := v.RegisterTranslation("required", translator, func(ut ut.Translator) error {
		return ut.Add("required", "{0} must have a value!", true)
	}, func(ut ut.Translator, fe validator.FieldError) string {
		t, _ := ut.T("required", strings.ToLower(fe.Field()))
		return t
	})
	if err != nil {
		return nil, err
	}

	return v, nil
}
