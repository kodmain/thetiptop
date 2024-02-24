// Package fizzbuzz grouping the data models linked to the fizzbuzz API
package fizzbuzz

// Stats is a struct that represents the statistics for a particular Request.
type Stats struct {
	// Request is the Request object associated with these statistics.
	Request Request
	// Hits is the number of times the Request has been made.
	Hits uint
}
