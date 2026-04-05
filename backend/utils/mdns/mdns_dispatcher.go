package mdns

import (
	"sync"
	"transok/backend/consts"
)

// Handler uses the interface defined in the handlers package
type Handler interface {
	GetType() string
	Handle(payload consts.DiscoverPayload)
}

// Dispatcher is the mDNS message dispatcher
type Dispatcher struct {
	handlers map[string][]Handler
	mu       sync.RWMutex
}

var (
	defaultDispatcher *Dispatcher
	once              sync.Once
)

// GetDispatcher returns the singleton dispatcher instance
func GetDispatcher() *Dispatcher {
	once.Do(func() {
		defaultDispatcher = &Dispatcher{
			handlers: make(map[string][]Handler),
		}
	})
	return defaultDispatcher
}

// Subscribe is modified to receive the Handler interface
func (d *Dispatcher) Subscribe(handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	messageType := handler.GetType()
	if d.handlers[messageType] == nil {
		d.handlers[messageType] = make([]Handler, 0)
	}
	d.handlers[messageType] = append(d.handlers[messageType], handler)
}

// Unsubscribe unsubscribes from processing a specific type of message
func (d *Dispatcher) Unsubscribe(messageType string, handler Handler) {
	d.mu.Lock()
	defer d.mu.Unlock()

	if handlers, exists := d.handlers[messageType]; exists {
		for i, h := range handlers {
			if &h == &handler {
				d.handlers[messageType] = append(handlers[:i], handlers[i+1:]...)
				break
			}
		}
	}
}

// Dispatch dispatches messages to the corresponding handlers
func (d *Dispatcher) Dispatch(payload consts.DiscoverPayload) {
	d.mu.RLock()
	handlers, exists := d.handlers[payload.Type]
	d.mu.RUnlock()

	if !exists {
		return
	}

	// Call all registered handlers
	for _, handler := range handlers {
		go handler.Handle(payload)
	}
}
