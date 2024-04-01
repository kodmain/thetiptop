package events

import "sync"

var (
	histories   map[TYPE][]any          = make(map[TYPE][]any)
	subscribers map[TYPE][]func(...any) = make(map[TYPE][]func(...any))
	mu          sync.RWMutex
)

// Notify notifie les subscribers de l'événement et stocke l'historique.
func Notify(event TYPE, data ...any) {
	mu.Lock()
	histories[event] = append(histories[event], data...)
	mu.Unlock()

	mu.RLock()
	if subscriber, ok := subscribers[event]; ok {
		for _, action := range subscriber {
			action(data...)
		}
	}
	mu.RUnlock()
}

// Subscribe ajoute un nouveau subscriber à un événement.
func Subscribe(event TYPE, subscriber func(...any)) {
	mu.Lock()
	defer mu.Unlock()

	if _, ok := subscribers[event]; !ok {
		subscribers[event] = []func(...any){}
	}

	subscribers[event] = append(subscribers[event], subscriber)
	if history, ok := histories[event]; ok {
		subscriber(history...)
	}
}
