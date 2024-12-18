package hook

import (
	"sync"

	"github.com/kodmain/thetiptop/api/env"
	"github.com/kodmain/thetiptop/api/internal/infrastructure/observability/logger"
)

type Handler func(tags ...string)

func (h Handler) once() bool          { return false }
func (h Handler) call(tags ...string) { h(tags...) }
func (h Handler) sync() bool          { return false }

type OnceHandler func(tags ...string)

func (h OnceHandler) once() bool          { return true }
func (h OnceHandler) call(tags ...string) { h(tags...) }
func (h OnceHandler) sync() bool          { return false }

type HandlerSync func(tags ...string)

func (h HandlerSync) once() bool          { return false }
func (h HandlerSync) call(tags ...string) { h(tags...) }
func (h HandlerSync) sync() bool          { return true }

type OnceHandlerSync func(tags ...string)

func (h OnceHandlerSync) once() bool          { return true }
func (h OnceHandlerSync) call(tags ...string) { h(tags...) }
func (h OnceHandlerSync) sync() bool          { return true }

type HookHandler interface {
	sync() bool
	once() bool
	call(tags ...string)
}

var (
	handlers = map[Event][]HookHandler{}
	history  = map[Event][]string{}
	mu       sync.Mutex
)

func Register(event Event, handler HookHandler) {
	mu.Lock()
	handlers[event] = append(handlers[event], handler)
	mu.Unlock()
	if _, ok := history[event]; ok {
		Call(event, history[event]...)
	}
}

func Call(event Event, tags ...string) {
	logger.Warnf("hook called for event %s with tags %v", event, tags)
	if !env.IsTest() {
		mu.Lock()
		defer mu.Unlock()

		history[event] = tags

		var remainingHandlers []HookHandler

		for _, handler := range handlers[event] {
			if !handler.sync() {
				go handler.call(tags...)
			} else {
				handler.call(tags...)
			}

			if !handler.once() {
				remainingHandlers = append(remainingHandlers, handler)
			}
		}

		// Remplacer les gestionnaires par la liste mise à jour
		handlers[event] = remainingHandlers
	}
}
