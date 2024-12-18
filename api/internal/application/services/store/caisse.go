package services

import (
	"github.com/gofiber/fiber/v2"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/application/validator"
	"github.com/kodmain/thetiptop/api/internal/domain/store/services"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/data"
)

func GetCaisse(store services.StoreServiceInterface, dtoCaisse *transfert.Caisse) (int, any) {
	if err := dtoCaisse.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	caisses, err := store.GetCaisse(dtoCaisse)
	if err != nil {
		return err.Code(), err
	}
	return fiber.StatusOK, caisses
}

func CreateCaisse(store services.StoreServiceInterface, dtoCaisse *transfert.Caisse) (int, any) {
	if err := dtoCaisse.Check(data.Validator{
		"store_id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	caisse, err := store.CreateCaisse(dtoCaisse)
	if err != nil {
		return err.Code(), err
	}
	return fiber.StatusCreated, caisse
}

func DeleteCaisse(store services.StoreServiceInterface, dtoCaisse *transfert.Caisse) (int, any) {
	if err := dtoCaisse.Check(data.Validator{
		"id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	err := store.DeleteCaisse(dtoCaisse)
	if err != nil {
		return err.Code(), err
	}
	return fiber.StatusNoContent, nil
}

func UpdateCaisse(store services.StoreServiceInterface, dtoCaisse *transfert.Caisse) (int, any) {
	if err := dtoCaisse.Check(data.Validator{
		"id":       {validator.Required, validator.ID},
		"store_id": {validator.Required, validator.ID},
	}); err != nil {
		return err.Code(), err
	}

	updatedCaisse, err := store.UpdateCaisse(dtoCaisse)
	if err != nil {
		return err.Code(), err
	}
	return fiber.StatusOK, updatedCaisse
}
