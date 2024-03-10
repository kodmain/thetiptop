// Package fizzbuzz grouping the data models linked to the fizzbuzz API
package fizzbuzz

// Request is a struct that represents a FizzBuzz request.
type Request struct {
	// Int1 is the first integer to be used as a divisor.
	Int1 int
	// Int2 is the second integer to be used as a divisor.
	Int2 int
	// Limit is the maximum number of iterations to perform.
	Limit int
	// Str1 is the string to print when the current iteration is divisible by Int1.
	Str1 string
	// Str2 is the string to print when the current iteration is divisible by Int2.
	Str2 string
}
