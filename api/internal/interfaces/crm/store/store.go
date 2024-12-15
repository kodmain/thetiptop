package store

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	services "github.com/kodmain/thetiptop/api/internal/application/services/store"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
	storeRepository "github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
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

// @Tags      Store
// @Summary   Get caisse by store
// @Produce   application/json
// @Param     store_id path string true "Store ID" format(uuid)
// @Success   200 {object} nil "List of caisse"
// @Failure   400 {object} nil "Invalid store"
// @Failure   500 {object} nil "Internal server error"
// @Router    /store/{id} [get]
// @Id        store.GetStoreByID
func GetStoreByID(ctx *fiber.Ctx) error {
	StoreID := ctx.Params("id")
	if StoreID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON("Store ID is required")
	}

	dto := &transfert.Store{
		ID: &StoreID,
	}

	status, response := services.GetStoreByID(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			storeRepository.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		), dto,
	)

	return ctx.Status(status).JSON(response)
}
