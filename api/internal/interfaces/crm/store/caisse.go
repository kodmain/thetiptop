package store

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/security"
	services "github.com/kodmain/thetiptop/api/internal/application/services/store"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	storeRepository "github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/store/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
)

// @Tags      Caisse
// @Summary   Get all caisse
// @Produce   application/json
// @Success   200 {object} entities.Caisse "List of caisse"
// @Failure   500 {object} nil "Internal server error"
// @Router    /caisse/{id} [get]
// @Id        store.GetCaisse
func GetCaisse(ctx *fiber.Ctx) error {
	clientID := ctx.Params("id")

	if clientID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON("Client ID is required")
	}

	dto := &transfert.Caisse{
		ID: &clientID,
	}

	status, response := services.GetCaisse(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			storeRepository.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		), dto,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags      Caisse
// @Summary   Get caisse by lieu
// @Produce   application/json
// @Param     store_id path string true "Store ID" format(uuid)
// @Success   200 {object} []entities.Caisse "List of caisse"
// @Failure   400 {object} nil "Invalid lieu"
// @Failure   500 {object} nil "Internal server error"
// @Router    /caisse/{lieu} [get]
// @Id        store.GetCaisseByStore
func GetCaisseByStore(ctx *fiber.Ctx) error {
	dto := &transfert.Caisse{}
	if err := ctx.BodyParser(dto); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err)
	}

	status, response := services.GetCaisseByStore(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			storeRepository.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		), dto,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags      Caisse
// @Accept	  multipart/form-data
// @Summary   Create a new caisse
// @Produce	  application/jsons
// @Param	  store_id formData string true "Store ID" format(uuid) default(440763b8-b8d9-4b36-9cc6-545a2c03071c)
// @Success   201 {object} entities.Caisse "Caisse created"
// @Failure   400 {object} nil "Invalid input"
// @Failure   500 {object} nil "Internal server error"
// @Router    /caisse [post]
// @Id        store.CreateCaisse
func CreateCaisse(ctx *fiber.Ctx) error {
	dtoCaisse := &transfert.Caisse{}
	if err := ctx.BodyParser(dtoCaisse); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	status, response := services.CreateCaisse(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			storeRepository.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		), dtoCaisse,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags      Caisse
// @Summary   Delete a caisse by ID
// @Produce   application/json
// @Param	  id   path	   string true "Client ID" format(uuid)
// @Success   204 {object} nil "Caisse deleted"
// @Failure   400 {object} nil "Invalid ID"
// @Failure   404 {object} nil "Caisse not found"
// @Failure   500 {object} nil "Internal server error"
// @Router    /caisse/{id} [delete]
// @Id        store.DeleteCaisse
func DeleteCaisse(ctx *fiber.Ctx) error {
	dtoCaisse := new(transfert.Caisse)
	if err := ctx.BodyParser(dtoCaisse); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	status, response := services.DeleteCaisse(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			storeRepository.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		), dtoCaisse,
	)

	return ctx.Status(status).JSON(response)
}

// @Tags      Caisse
// @Summary   Update a caisse by ID
// @Accept    application/json
// @Produce   application/json
// @Param	  id   path	   string true "Client ID" format(uuid)
// @Success   200 {object} entities.Caisse "Caisse updated"
// @Failure   400 {object} nil "Invalid input"
// @Failure   404 {object} nil "Caisse not found"
// @Failure   500 {object} nil "Internal server error"
// @Router    /caisse/{id} [put]
// @Id        store.UpdateCaisse
func UpdateCaisse(ctx *fiber.Ctx) error {
	dtoCaisse := new(transfert.Caisse)
	if err := ctx.BodyParser(dtoCaisse); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(err.Error())
	}

	status, response := services.UpdateCaisse(
		domain.Store(
			security.NewUserAccess(ctx.Locals("token")),
			storeRepository.NewStoreRepository(database.Get(config.GetString("services.store.database", config.DEFAULT))),
		), dtoCaisse,
	)

	return ctx.Status(status).JSON(response)
}
