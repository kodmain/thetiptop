package user

import (
	"github.com/gofiber/fiber/v2"

	"github.com/kodmain/thetiptop/api/config"
	"github.com/kodmain/thetiptop/api/internal/application/services"
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/user/repositories"
	domain "github.com/kodmain/thetiptop/api/internal/domain/user/services"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/database"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/providers/mail"
)

// @Tags		Employee
// @Accept		multipart/form-data
// @Summary		Register a employee.
// @Produce		application/json
// @Param		email		formData	string	true	"Email address" format(email) default(user-thetiptop@yopmail.com)
// @Param		password	formData	string	true	"Password" default(Aa1@azetyuiop)
// @Success		201	{object}	nil "Employee created"
// @Failure		400	{object}	nil "Invalid email or password"
// @Failure		409	{object}	nil "Employee already exists"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/employee/register [post]
// @Id			user.RegisterEmployee
func RegisterEmployee(c *fiber.Ctx) error {
	dtoCredential := &transfert.Credential{}
	if err := c.BodyParser(dtoCredential); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	dtoEmployee := &transfert.Employee{}
	if err := c.BodyParser(dtoEmployee); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.RegisterEmployee(
		domain.User(
			repositories.NewUserRepository(database.Get(config.Get("services.employee.database", "default").(string))),
			mail.Get(),
		), dtoCredential, dtoEmployee,
	)

	return c.Status(status).JSON(response)
}

// @Tags		Employee
// @Accept		multipart/form-data
// @Summary		Update a employee.
// @Produce		application/json
// @Param		id			formData	string	true	"Employee ID" format(uuid)
// @Param		newsletter	formData	bool	true	"Newsletter" default(false)
// @Success		204	{object}	nil "Password updated"
// @Failure		400	{object}	nil "Invalid email, password or token"
// @Failure		404	{object}	nil "Employee not found"
// @Failure		409	{object}	nil "Employee already validated"
// @Failure		410	{object}	nil "Token expired"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/employee [put]
// @Id			user.UpdateEmployee
func UpdateEmployee(c *fiber.Ctx) error {
	dtoEmployee := &transfert.Employee{}
	if err := c.BodyParser(dtoEmployee); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	status, response := services.UpdateEmployee(
		domain.User(
			repositories.NewUserRepository(database.Get(config.Get("services.employee.database", "default").(string))),
			mail.Get(),
		), dtoEmployee,
	)

	return c.Status(status).JSON(response)
}

// @Tags		Employee
// @Accept		multipart/form-data
// @Summary		Get a employee by ID.
// @Produce		application/json
// @Param		id			path		string	true	"Employee ID" format(uuid)
// @Success		200	{object}	nil "Employee details"
// @Failure		400	{object}	nil "Invalid employee ID"
// @Failure		404	{object}	nil "Employee not found"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/employee/{id} [get]
// @Id			user.GetEmployee
func GetEmployee(c *fiber.Ctx) error {
	EmployeeID := c.Params("id")

	if EmployeeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee ID is required")
	}

	dtoEmployee := &transfert.Employee{
		ID: &EmployeeID,
	}

	status, response := services.GetEmployee(
		domain.User(
			repositories.NewUserRepository(database.Get(config.Get("services.employee.database", "default").(string))),
			mail.Get(),
		), dtoEmployee,
	)

	return c.Status(status).JSON(response)
}

// @Tags		Employee
// @Summary		Delete a client by ID.
// @Produce		application/json
// @Param		id			path		string	true	"Employee ID" format(uuid)
// @Success		204	{object}	nil "Employee deleted"
// @Failure		400	{object}	nil "Invalid employee ID"
// @Failure		404	{object}	nil "Employee not found"
// @Failure		500	{object}	nil "Internal server error"
// @Router		/employee/{id} [delete]
// @Id			user.DeleteEmployee
func DeleteEmployee(c *fiber.Ctx) error {
	EmployeeID := c.Params("id")

	if EmployeeID == "" {
		return fiber.NewError(fiber.StatusBadRequest, "Employee ID is required")
	}

	dtoEmployee := &transfert.Employee{
		ID: &EmployeeID,
	}

	status, response := services.DeleteEmployee(
		domain.User(
			repositories.NewUserRepository(database.Get(config.Get("services.employee.database", "default").(string))),
			mail.Get(),
		), dtoEmployee,
	)

	return c.Status(status).JSON(response)
}
