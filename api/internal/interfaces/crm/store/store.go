package store

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	services "github.com/kodmain/thetiptop/api/internal/application/services/store"
	"github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/store/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

// @Tags		Store
// @Accept		multipart/form-data
// @Summary		List all store.
// @Produce		application/json
// @Success		200	{object}	nil "list of store"
// @Router		/store [get]
// @Id			store.List
func List(ctx *fiber.Ctx) error {
	status, response := services.ListStores(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			repositories.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		),
	)

	return ctx.Status(status).JSON(response)
}
