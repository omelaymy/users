package delivery

import (
	"errors"
	"strings"

	ut "github.com/go-playground/universal-translator"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"github.com/omelaymy/users/internal/api"
	apiErrors "github.com/omelaymy/users/internal/api/http/errors"
	"github.com/omelaymy/users/internal/users"
)

type Handlers struct {
	usersUsecase     users.Usecase
	validate         *validator.Validate
	errorsTranslator ut.Translator
}

func NewHandlers(
	usersUsecase users.Usecase,
	validate *validator.Validate,
	errorsTranslator ut.Translator,

) *Handlers {
	return &Handlers{
		usersUsecase:     usersUsecase,
		validate:         validate,
		errorsTranslator: errorsTranslator,
	}
}

// @Summary Get User Information
// @Description Get information about a specific user
// @Tags Users
// @Produce json
// @Param id path string true "User ID"
// @Security BasicAuth
// @Success 200 {object} api.UserResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /v1/users/{id} [get]
func (h *Handlers) GetUserHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				apiErrors.InvalidId,
			)
		}

		user, err := h.usersUsecase.GetUser(id)
		if err != nil {
			code := fiber.StatusInternalServerError
			if errors.Is(err, users.UserNotFoundError) {
				code = fiber.StatusNotFound
			}
			return fiber.NewError(code, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(api.UserResponse{
			Id:       user.Id,
			Email:    user.Email,
			Username: user.Username,
			Admin:    user.Admin,
		})
	}
}

// @Summary Update User
// @Description Update a user with the provided information (requires admin access)
// @Tags Users
// @Param id path string true "User ID"
// @Param user body api.UserRequest true "User object to update"
// @Security BasicAuth
// @Success 200 {object} api.SuccessResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 404 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /v1/users/{id} [put]
func (h *Handlers) UpdateUserHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				apiErrors.InvalidId,
			)
		}

		var user api.UserRequest
		if err = c.BodyParser(&user); err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				apiErrors.InvalidRequestBodyError,
				err.Error(),
			)
		}

		if err = h.validate.StructCtx(c.Context(), &user); err != nil {
			errs := err.(validator.ValidationErrors)
			return fiber.NewError(
				fiber.StatusBadRequest, formattingValidatorErrors(h.errorsTranslator, errs),
			)
		}

		err = h.usersUsecase.UpdateUser(&users.User{
			Id:       id,
			Email:    user.Email,
			Username: user.Username,
			Password: user.Password,
			Admin:    user.Admin,
		})
		if err != nil {
			code := fiber.StatusInternalServerError
			if errors.Is(err, users.UserNotFoundError) {
				code = fiber.StatusNotFound
			}
			if errors.Is(err, users.UserAlreadyExistsError) {
				code = fiber.StatusBadRequest
			}
			return fiber.NewError(code, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(
			api.SuccessResponse{
				Success: true,
			},
		)
	}
}

// @Summary Create User
// @Description Create a new user with the provided information (requires admin access)
// @Tags Users
// @Accept json
// @Produce json
// @Param user body api.UserRequest true "User object to create"
// @Security BasicAuth
// @Success 200 {object} api.UserIdResponse
// @Failure 400 {object} api.ErrorResponse
// @Failure 500 {object} api.ErrorResponse
// @Router /v1/users [post]
func (h *Handlers) CreateUserHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		var user api.UserRequest
		if err := c.BodyParser(&user); err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				apiErrors.InvalidRequestBodyError,
				err.Error(),
			)
		}

		if err := h.validate.StructCtx(c.Context(), &user); err != nil {
			errs := err.(validator.ValidationErrors)
			return fiber.NewError(
				fiber.StatusBadRequest, formattingValidatorErrors(h.errorsTranslator, errs),
			)
		}

		id, err := h.usersUsecase.CreateUser(&users.User{
			Email:    user.Email,
			Username: user.Username,
			Password: user.Password,
			Admin:    user.Admin,
		})
		if err != nil {
			code := fiber.StatusInternalServerError
			if errors.Is(err, users.UserAlreadyExistsError) {
				code = fiber.StatusBadRequest
			}
			return fiber.NewError(code, err.Error())
		}

		return c.Status(fiber.StatusOK).JSON(api.UserIdResponse{
			Id: id,
		})
	}
}

// @Summary Get Users
// @Description Get a list of all users
// @Tags Users
// @Produce json
// @Security BasicAuth
// @Success 200 {array} api.UserResponse
// @Router /v1/users [get]
func (h *Handlers) GetUsersHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(h.usersUsecase.GetUsers())
	}
}

// @Summary Delete User
// @Description Delete a user by ID (requires admin access)
// @Tags Users
// @Param id path string true "User ID"
// @Security BasicAuth
// @Success 200 {object} api.SuccessResponse "User deleted successfully"
// @Failure 400 {object} api.ErrorResponse
// @Router /v1/users/{id} [delete]
func (h *Handlers) DeleteUserHandler() fiber.Handler {
	return func(c *fiber.Ctx) error {
		id, err := uuid.Parse(c.Params("id"))
		if err != nil {
			return fiber.NewError(
				fiber.StatusBadRequest,
				apiErrors.InvalidId,
			)
		}

		h.usersUsecase.DeleteUser(id)

		return c.Status(fiber.StatusOK).JSON(
			api.SuccessResponse{
				Success: true,
			},
		)
	}
}

func formattingValidatorErrors(tr ut.Translator, errs validator.ValidationErrors) string {
	var sb strings.Builder
	for i, e := range errs {
		sb.WriteString(e.Translate(tr))
		if i != len(errs)-1 {
			sb.WriteString("\n")
		}
	}
	return sb.String()
}
