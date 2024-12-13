package events

import (
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	transfert "github.com/kodmain/thetiptop/api/internal/application/transfert/crm"
	"github.com/kodmain/thetiptop/api/internal/domain/store/entities"
	"github.com/kodmain/thetiptop/api/internal/domain/store/repositories"
)

// ConvertEntityToTransfer Converts an entities.Store to transfert.Store
//
// Parameters:
// - store: *entities.Store The store entity to convert.
//
// Returns:
// - *transfert.Store: The converted store for transfer.
func ConvertEntityToTransfer(store *entities.Store) *transfert.Store {
	return &transfert.Store{
		Label:    store.Label,
		IsOnline: store.IsOnline,
	}
}

// ConvertTransferToEntity Converts a transfert.Store to entities.Store
//
// Parameters:
// - store: *transfert.Store The store transfer object to convert.
//
// Returns:
// - *entities.Store: The converted store entity.
func ConvertTransferToEntity(store *transfert.Store) *entities.Store {
	return &entities.Store{
		Label:    store.Label,
		IsOnline: store.IsOnline,
	}
}

// CreateStores Ensures stores in the database match the desired state
func CreateStores(repo repositories.StoreRepositoryInterface) {
	desiredStores := []*transfert.Store{
		{Label: aws.String("DigitalStore"), IsOnline: aws.Bool(true)},
		{Label: aws.String("PhysicalStore"), IsOnline: aws.Bool(false)},
	}

	// Retrieve existing stores from the database
	existingEntityStores, err := repo.ReadStores(&transfert.Store{})
	if err != nil {
		panic(fmt.Sprintf("Failed to read stores: %v", err))
	}

	// Convert existing entities to transfer objects
	var existingStores []*transfert.Store
	for _, entityStore := range existingEntityStores {
		existingStores = append(existingStores, ConvertEntityToTransfer(entityStore))
	}

	// Map existing stores by label for quick lookup
	existingStoreMap := make(map[string]*transfert.Store)
	for _, store := range existingStores {
		existingStoreMap[*store.Label] = store
	}

	// Prepare stores to add and remove
	var storesToAdd []*transfert.Store
	var storesToRemove []*transfert.Store

	// Determine stores to add
	for _, desiredStore := range desiredStores {
		if _, exists := existingStoreMap[*desiredStore.Label]; !exists {
			storesToAdd = append(storesToAdd, desiredStore)
		}
	}

	// Determine stores to remove
	desiredStoreMap := make(map[string]bool)
	for _, desiredStore := range desiredStores {
		desiredStoreMap[*desiredStore.Label] = true
	}
	for _, existingStore := range existingStores {
		if !desiredStoreMap[*existingStore.Label] {
			storesToRemove = append(storesToRemove, existingStore)
		}
	}

	// Create stores that are missing
	if len(storesToAdd) > 0 {
		if err := repo.CreateStores(storesToAdd); err != nil {
			panic(fmt.Sprintf("Failed to insert stores: %v", err))
		}
		fmt.Printf("%d stores were added\n", len(storesToAdd))

		// Retrieve existing stores from the database
		existingEntityStores, err = repo.ReadStores(&transfert.Store{})
		if err != nil {
			panic(fmt.Sprintf("Failed to read stores: %v", err))
		}

		// Convert existing entities to transfer objects
		for _, entityStore := range existingEntityStores {
			entityStore.Caisses = []*entities.Caisse{
				{StoreID: &entityStore.ID},
				{StoreID: &entityStore.ID},
				{StoreID: &entityStore.ID},
				{StoreID: &entityStore.ID},
			}
		}

		repo.UpdateStores(existingEntityStores)

		fmt.Println("Caisse synchronization completed")
	}

	// Remove stores that are extra
	if len(storesToRemove) > 0 {
		if err := repo.DeleteStores(storesToRemove); err != nil {
			panic(fmt.Sprintf("Failed to delete stores: %v", err))
		}
		fmt.Printf("%d stores were removed\n", len(storesToRemove))
	}

	fmt.Println("Store synchronization completed")
}
