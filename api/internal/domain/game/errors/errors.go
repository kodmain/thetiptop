package errors_domain_game

import (
	"net/http"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

var (
	// Ticket errors
	ErrTicketNotFound = errors.New(http.StatusNotFound, "ticket.not_found")
)
