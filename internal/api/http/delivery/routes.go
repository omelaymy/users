package delivery

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/swagger"
	"github.com/omelaymy/users/internal/api"
)

type Routes struct {
	h      *Handlers
	mw     *api.MWManager
	router fiber.Router
}

func NewRoutes(
	h *Handlers,
	mw *api.MWManager,
	router fiber.Router,
) *Routes {
	return &Routes{
		h:      h,
		mw:     mw,
		router: router,
	}
}

func (r *Routes) RegisterRoutes() {
	r.router.Get("/docs/*", swagger.HandlerDefault)

	api := r.router.Group("/api")
	v1 := api.Group("/v1")
	users := v1.Group("/users").Use(r.mw.BasicAuth())

	users.Get("", r.h.GetUsersHandler())
	users.Get("/:id<guid>", r.h.GetUserHandler())
	users.Post("", r.mw.AdminAuth(), r.h.CreateUserHandler())
	users.Put("/:id<guid>", r.mw.AdminAuth(), r.h.UpdateUserHandler())
	users.Delete("/:id<guid>", r.mw.AdminAuth(), r.h.DeleteUserHandler())
}
