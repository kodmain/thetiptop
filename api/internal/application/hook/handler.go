package hook

import "sync"

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

func Register(event Event, handler Handler) {
	mu.Lock()
	handlers[event] = append(handlers[event], handler)
	mu.Unlock()
}

func Call(event Event) {
	mu.Lock()
	for _, handler := range handlers[event] {
		go handler.call()
	}
	mu.Unlock()
}
