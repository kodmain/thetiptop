package data

import (
	"reflect"
)

// UpdateEntityWithDto met à jour une entité en utilisant un DTO
//
// Parameters:
// - entity: interface{} : Entité à mettre à jour
// - dto: interface{} : DTO contenant les nouvelles valeurs
//
// Returns:
// - error: erreur si quelque chose ne va pas
func UpdateEntityWithDto(entity interface{}, dto interface{}) {
	// Récupère les valeurs des structs entity et DTO
	entityVal := reflect.ValueOf(entity).Elem()
	dtoVal := reflect.ValueOf(dto).Elem()

	// Parcourt les champs du DTO pour mettre à jour les champs correspondants dans l'entité
	for i := 0; i < dtoVal.NumField(); i++ {
		dtoField := dtoVal.Field(i)
		entityField := entityVal.FieldByName(dtoVal.Type().Field(i).Name)

		// Vérifie que le champ de l'entité existe et est assignable
		if entityField.IsValid() && entityField.CanSet() {
			// Si le champ du DTO est un pointeur et non nil
			if dtoField.Kind() == reflect.Ptr && !dtoField.IsNil() {
				// Si le champ correspondant dans l'entité est aussi un pointeur
				if entityField.Kind() == reflect.Ptr {
					// Compare les valeurs pointées, et non les pointeurs eux-mêmes
					if !reflect.DeepEqual(entityField.Elem().Interface(), dtoField.Elem().Interface()) {
						entityField.Set(dtoField) // Assigner directement le pointeur si les valeurs sont différentes
					}
				} else {
					// Si le champ de l'entité n'est pas un pointeur, dé-référencer le pointeur du DTO et comparer
					if !reflect.DeepEqual(entityField.Interface(), dtoField.Elem().Interface()) {
						entityField.Set(dtoField.Elem()) // Assigner la valeur pointée si elle est différente
					}
				}
			} else if dtoField.Kind() != reflect.Ptr { // Si ce n'est pas un pointeur, on compare les valeurs directement
				if !reflect.DeepEqual(entityField.Interface(), dtoField.Interface()) {
					entityField.Set(dtoField) // Assigner la valeur si elle est différente
				}
			}
		}
	}
}
