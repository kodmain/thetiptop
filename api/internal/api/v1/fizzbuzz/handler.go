// Package fizzbuzz implementation all handler for FizzBuzz API
package fizzbuzz

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/kodmain/thetiptop/api/internal/model/fizzbuzz"
)

// FizzBuzz generates a FizzBuzz sequence based on the values of `Int1`, `Int2`, `Limit`, `Str1`, and `Str2`
// in the `fizzbuzz.Request` object stored in the request context `c`. The generated sequence is returned as
// a JSON response.
//
// The `Int1` and `Int2` values are used as divisors to determine if a given number is divisible by either one
// of them or both. If a number is divisible by both, the strings `Str1` and `Str2` are concatenated, otherwise,
// only the appropriate string is appended to the result slice. If the number is not divisible by either one of
// the integers, the number is appended as a string to the result slice.
//
// The generated FizzBuzz sequence is returned as a JSON response in the same order as it was generated.
// @Summary		Return FizzBuzz result.
// @Description	Returns a list of strings with numbers from 1 to limit, where: \n all multiples of int1 are replaced by str1, all multiples of int2 are replaced by str2, all multiples of int1 and int2 are replaced by str1str2.
// @Tags		FizzBuzz
// @Accept		*/*
// @Produce		json
// @Help		Name    ?        type    required   description
// @Param       int1    path     int     true       "Give the first number"  default(3)
// @Param       int2    path     int     true       "Give the second number" default(5)
// @Param       limit   path     int     true       "Give limit of fizzbuzz" default(100)
// @Param       str1    path     string  true       "Give the first word"    default(fizz)
// @Param       str2    path     string  true       "Give the second word"   default(buzz)
// @Failure     400  	{array}  string
// @Success		200		{array}  string
// @Router		/api/v1/fizzbuzz/:int1/:int2/:limit/:str1/:str2 [get]
func FizzBuzz(c *fiber.Ctx) error {
	var fbr fizzbuzz.Request = c.Context().UserValue("fizzbuzz.Request").(fizzbuzz.Request)
	var result []string

	for i := 1; i <= fbr.Limit; i++ {
		if i%fbr.Int1 == 0 && i%fbr.Int2 == 0 {
			result = append(result, fbr.Str1+fbr.Str2)
		} else if i%fbr.Int1 == 0 {
			result = append(result, fbr.Str1)
		} else if i%fbr.Int2 == 0 {
			result = append(result, fbr.Str2)
		} else {
			result = append(result, strconv.Itoa(i))
		}
	}

	return c.JSON(result)
}

// FizzBuzzControls takes in 5 parameters from a Fiber HTTP Context object: int1, int2, limit, str1, and str2.
// The function first checks if int1, int2, and limit are valid integers. If any of these parameters are not valid integers,
// it will return a response with status code 400 and an error message.
// If int1 or int2 are equal to 0, or if limit is equal to 0, the function will return a response with status code 400 and an error message.
// If all the parameters are valid, it will create a new FizzBuzz request object with these parameters and store it in the context's
// UserValue field under the key "fizzbuzz.Request". It then calls the Next() function to pass the context object to the next middleware or handler.
// If an error occurs at any point, it will return a response with status code 400 and an error message.
func FizzBuzzControls(c *fiber.Ctx) error {
	var errors = map[string]string{}

	int1, err := strconv.Atoi(c.Params("int1"))
	if err != nil {
		errors["int1"] = "int1 parameter is not an integer"
	}

	int2, err := strconv.Atoi(c.Params("int2"))
	if err != nil {
		errors["int2"] = "int2 parameter is not an integer"
	}

	limit, err := strconv.Atoi(c.Params("limit"))
	if err != nil {
		errors["limit"] = "limit parameter is not an integer"
	}

	if int1 == 0 || int2 == 0 {
		errors["int_greater"] = "int1 and int2 parameters must be greater than 0"
	}

	if limit == 0 {
		errors["limit_greater"] = "limit parameters must be greater than 0"
	}

	if len(errors) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(errors)
	}

	str1 := c.Params("str1")
	str2 := c.Params("str2")

	fbr := fizzbuzz.Request{
		Int1:  int1,
		Int2:  int2,
		Limit: limit,
		Str1:  str1,
		Str2:  str2,
	}

	c.Context().SetUserValue("fizzbuzz.Request", fbr)

	return c.Next()
}
