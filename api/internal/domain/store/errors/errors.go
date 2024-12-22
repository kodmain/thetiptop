package errors_domain_store

import (
	"net/http"

	"github.com/kodmain/thetiptop/api/internal/infrastructure/errors"
)

var (
	// Store errors
	ErrStoreNotFound = errors.New(http.StatusNotFound, "store.not_found")
	// Caisse errors
	ErrCaisseNotFound = errors.New(http.StatusNotFound, "caisse.not_found")
)
