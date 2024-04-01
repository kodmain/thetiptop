package repositories

import (
	"github.com/google/uuid"
	"github.com/kodmain/thetiptop/api/internal/domain/entities"
)

// ClientRepository Interface for managing client data storage operations.
//
// This interface defines methods for interacting with client data, enabling CRUD (Create, Read, Update, Delete) operations on clients.
type ClientRepository interface {

	// GetClient Retrieves a client by their UUID identifier.
	//
	// This method is used to obtain the details of a specific client using their unique identifier.
	// It returns an error if the client is not found or if a database error occurs.
	//
	// Parameters:
	// - id: uuid.UUID The UUID identifier of the client to retrieve.
	//
	// Returns:
	// - *entities.Client: The found client, or nil in case of error.
	// - error: Error returned in case of problems during retrieval.
	GetClient(id uuid.UUID) (*entities.Client, error)

	// GetClients Retrieves a list of clients based on a filter.
	//
	// This method returns a list of clients that match the criteria specified in the filter.
	// The filter is a map of key-value pairs where keys are field names and values are filtering criteria.
	//
	// Parameters:
	// - filter: map[string]any The filter criteria for retrieving clients.
	//
	// Returns:
	// - []*entities.Client: A slice of clients matching the filter criteria.
	// - error: Error returned in case of problems during retrieval.
	GetClients(filter map[string]any) ([]*entities.Client, error)

	// CreateClient Creates a new client.
	//
	// This method adds a new client to the storage. It returns the created client with its assigned identifier.
	//
	// Parameters:
	// - client: *entities.Client The client to be created.
	//
	// Returns:
	// - *entities.Client: The created client.
	// - error: Error returned in case of problems during creation.
	CreateClient(client *entities.Client) (*entities.Client, error)

	// UpdateClient Updates an existing client.
	//
	// This method modifies the details of an existing client. It returns the updated client information.
	//
	// Parameters:
	// - client: *entities.Client The client information to be updated.
	//
	// Returns:
	// - *entities.Client: The updated client.
	// - error: Error returned in case of problems during the update.
	UpdateClient(client *entities.Client) (*entities.Client, error)

	// DeleteClient Deletes a client by their UUID identifier.
	//
	// This method removes a client from the storage using their UUID. It returns an error if the deletion process fails.
	//
	// Parameters:
	// - id: uuid.UUID The UUID identifier of the client to be deleted.
	//
	// Returns:
	// - error: Error returned in case of problems during deletion.
	DeleteClient(id uuid.UUID) error
}
