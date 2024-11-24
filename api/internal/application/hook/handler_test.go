package hook_test

import (
	"testing"
	"time"

	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/application/hook"
	"github.com/stretchr/testify/assert"
)

type Event string

func TestRegister(t *testing.T) {
	env.ForceTest()
	var ENAME hook.Event = "EVENT"
	var countHandler int = 0
	var countOnceHandler int = 0

	var handler hook.Handler = func(tags ...string) {
		countHandler++
	}

	var OnceHandler hook.OnceHandler = func(tags ...string) {
		countOnceHandler++
	}

	var OnceHandlerSync hook.OnceHandlerSync = func(tags ...string) {
		countOnceHandler++
	}

	var HandlerSync hook.HandlerSync = func(tags ...string) {
		countHandler++
	}

	// Enregistre l'événement avec le handler
	hook.Register(ENAME, handler)
	hook.Register(ENAME, OnceHandler)

	// Appelle l'événement
	hook.Call(ENAME)
	hook.Call(ENAME)
	// Attend que le handler ait le temps d'exécuter (les appels sont asynchrones)
	time.Sleep(1000 * time.Millisecond)

	// Vérifie que la valeur testValue est mise à jour par le handler
	assert.Equal(t, 2, countHandler)
	assert.Equal(t, 1, countOnceHandler)

	// Enregistre l'événement avec le handler
	hook.Register(ENAME, HandlerSync)
	hook.Register(ENAME, OnceHandlerSync)

	// Appelle l'événement
	hook.Call(ENAME)
	hook.Call(ENAME)
	// Attend que le handler ait le temps d'exécuter (les appels sont asynchrones)
	time.Sleep(1000 * time.Millisecond)

	// Vérifie que la valeur testValue est mise à jour par le handler
	assert.Equal(t, 6, countHandler)
	assert.Equal(t, 2, countOnceHandler)
}
