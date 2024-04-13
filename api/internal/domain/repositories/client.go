package repositories

import (
	"github.com/kodmain/thetiptop/api/internal/application/transfert"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
)

// ClientRepository Interface for managing client data storage operations.
//
// This interface defines methods for interacting with client data, enabling CRUD (Create, Read, Update, Delete) operations on clients.
type ClientRepository interface {
	// Create Creates a new client.
	//
	// This method adds a new client to the storage. It returns the created client with its assigned identifier.
	//
	// Parameters:
	// - client: *entities.Client The client to be created.
	//
	// Returns:
	// - *entities.Client: The created client.
	// - error: Error returned in case of problems during creation.
	Create(client *transfert.Client) (*entities.Client, error)
	Read(client *transfert.Client) (*entities.Client, error)
}
