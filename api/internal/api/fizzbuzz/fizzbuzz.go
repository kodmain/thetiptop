package fizzbuzz

import (
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/domain/services"
)

const start int = 1
const max int = 100

// It returns a JSON response with the fizzbuzz result
// @Summary      Run FizzBuzz
// @Description  Executes the FizzBuzz algorithm based on the provided parameters.
// @Tags         FizzBuzz
// @Accept       */*
// @Produce      json
// @Param        int1   path int    true  "First integer to replace with str1"
// @Param        int2   path int    true  "Second integer to replace with str2"
// @Param        limit  path int    true  "The upper limit for the FizzBuzz sequence"
// @Param        str1   path string true  "String to replace multiples of int1"
// @Param        str2   path string true  "String to replace multiples of int2"
// @Success      200    {object}     []string "A list of strings representing the FizzBuzz sequence"
// @Failure      400    {object}     object   "Bad Request - invalid input parameters"
// @Router       /fizzbuzz/:int1/:int2/:limit/:str1/:str2 [get]
// @Id           metrics.Counter => fizzbuzz.FizzBuzz
func FizzBuzz(c *fiber.Ctx) error {
	var errors map[string]string = make(map[string]string)

	int1, err := strconv.Atoi(c.Params("int1"))
	if err != nil {
		errors["int1"] = "invalid int1"
	}
	int2, err := strconv.Atoi(c.Params("int2"))
	if err != nil {
		errors["int2"] = "invalid int2"
	}
	limit, err := strconv.Atoi(c.Params("limit"))
	if err != nil {
		errors["limit"] = "invalid limit"
	}

	if limit < start || limit > max {
		errors["limit_range"] = fmt.Sprintf("limit must be between %d and %d", start, max)
	}

	if int1 <= 0 || int2 <= 0 || limit <= 0 {
		errors["no_zero"] = "int1, int2 and limit must be greater than 0"
	}

	if int1 == int2 {
		errors["different"] = "int1 and int2 must be different"
	}

	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": errors})
	}

	str1 := c.Params("str1")
	str2 := c.Params("str2")

	fizzBuzzService := services.NewFizzBuzzService()
	result := fizzBuzzService.RunFizzBuzz(int1, int2, limit, str1, str2)

	return c.JSON(result)
}
