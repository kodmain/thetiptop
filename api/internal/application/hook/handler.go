package hook

import (
	"sync"

	"github.com/kodmain/thetiptop/api/env"
)

type Handler func()

func (Handler) IsOnce() bool { return false }
func (h Handler) call()      { h() }

type OnceHandler func()

func (OnceHandler) IsOnce() bool { return true }
func (h OnceHandler) call()      { h() }

type HookHandler interface {
	IsOnce() bool
	call()
}

var (
	handlers = map[Event][]HookHandler{}
	mu       sync.Mutex
)

func Register(event Event, handler HookHandler) {
	mu.Lock()
	handlers[event] = append(handlers[event], handler)
	mu.Unlock()
}

func Call(event Event) {
	if !env.IsTest() {
		mu.Lock()
		defer mu.Unlock()

		// Liste temporaire pour les gestionnaires restants
		var remainingHandlers []HookHandler

		for _, handler := range handlers[event] {
			go handler.call()
			if !handler.IsOnce() {
				remainingHandlers = append(remainingHandlers, handler)
			}
		}

		// Remplacer les gestionnaires par la liste mise à jour
		handlers[event] = remainingHandlers
	}
}
