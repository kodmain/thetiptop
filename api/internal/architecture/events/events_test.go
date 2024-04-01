package events_test

import (
	"reflect"
	"testing"

	"github.com/kodmain/thetiptop/api/internal/architecture/events"
)

func TestSubscribeAndNotify(t *testing.T) {
	// Définition d'un événement test et des données test
	var testEvent events.TYPE = 0
	testData := []any{"data1", "data2"}
	testData2 := []any{"data3", "data4"}

	// Variable pour stocker les données reçues par le subscriber
	var receivedData []any

	// Subscriber test qui stocke les données reçues
	testSubscriber := func(data ...any) {
		receivedData = append(receivedData, data...)
	}

	// Notifier avec testData
	events.Notify(testEvent, testData...)

	// Souscrire au testEvent
	events.Subscribe(testEvent, testSubscriber)

	// Notifier avec testData2
	events.Notify(testEvent, testData2...)

	expectedData := append(testData, testData2...)
	if !reflect.DeepEqual(receivedData, expectedData) {
		t.Errorf("Expected received data to be %v, got %v", expectedData, receivedData)
	}
}
