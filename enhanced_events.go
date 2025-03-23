package mediasoupgo

import (
	"fmt"
	"sync"
)

// EnhancedEventEmitter is a Go implementation of an event emitter.
type EnhancedEventEmitter struct {
	listeners map[string][]func(...interface{})
	mutex     sync.RWMutex
}

// NewEnhancedEventEmitter creates a new EnhancedEventEmitter.
func NewEnhancedEventEmitter() *EnhancedEventEmitter {
	return &EnhancedEventEmitter{
		listeners: make(map[string][]func(...interface{})),
	}
}

// Emit emits an event with the given name and arguments.
func (e *EnhancedEventEmitter) Emit(eventName string, args ...interface{}) bool {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	listeners, ok := e.listeners[eventName]
	if !ok {
		return false
	}

	for _, listener := range listeners {
		listener(args...)
	}
	return true
}

// SafeEmit emits an event safely, catching any panics and emitting a "listenererror" event if needed.
func (e *EnhancedEventEmitter) SafeEmit(eventName string, args ...interface{}) bool {
	defer func() {
		if r := recover(); r != nil {
			err, _ := r.(error)
			if err == nil {
				err = r.(error)
			}
			// Emit listenererror event
			e.Emit("listenererror", eventName, err)
		}
	}()

	listeners, ok := e.listeners[eventName]
	if !ok {
		return false
	}

	for _, listener := range listeners {
		listener(args...)
	}
	return true
}

// On adds a listener for the given event.
func (e *EnhancedEventEmitter) On(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.listeners[eventName] = append(e.listeners[eventName], listener)
	return e
}

// Off removes a listener for the given event.
func (e *EnhancedEventEmitter) Off(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	listeners, ok := e.listeners[eventName]
	if !ok {
		return e
	}

	for i, l := range listeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			e.listeners[eventName] = append(listeners[:i], listeners[i+1:]...)
			break
		}
	}
	return e
}

// AddListener adds a listener for the given event.
func (e *EnhancedEventEmitter) AddListener(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	return e.On(eventName, listener)
}

// PrependListener adds a listener to the beginning of the listeners array for the given event.
func (e *EnhancedEventEmitter) PrependListener(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.listeners[eventName] = append([]func(...interface{}){listener}, e.listeners[eventName]...)
	return e
}

// Once adds a one-time listener for the given event.
func (e *EnhancedEventEmitter) Once(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	onceListener := func(args ...interface{}) {
		e.Off(eventName, listener)
		listener(args...)
	}
	return e.On(eventName, onceListener)
}

// PrependOnceListener adds a one-time listener to the beginning of the listeners array for the given event.
func (e *EnhancedEventEmitter) PrependOnceListener(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	onceListener := func(args ...interface{}) {
		e.Off(eventName, listener)
		listener(args...)
	}
	return e.PrependListener(eventName, onceListener)
}

// RemoveListener removes a listener for the given event.
func (e *EnhancedEventEmitter) RemoveListener(eventName string, listener func(...interface{})) *EnhancedEventEmitter {
	return e.Off(eventName, listener)
}

// RemoveAllListeners removes all listeners for the given event, or all listeners if no event is specified.
func (e *EnhancedEventEmitter) RemoveAllListeners(eventName ...string) *EnhancedEventEmitter {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	if len(eventName) == 0 {
		e.listeners = make(map[string][]func(...interface{}))
	} else {
		delete(e.listeners, eventName[0])
	}
	return e
}

// ListenerCount returns the number of listeners for the given event.
func (e *EnhancedEventEmitter) ListenerCount(eventName string) int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return len(e.listeners[eventName])
}

// EnhancedOnce waits for a single event to be emitted and returns the arguments.
func EnhancedOnce(emitter *EnhancedEventEmitter, eventName string) ([]interface{}, error) {
	done := make(chan []interface{})
	var err error

	listener := func(args ...interface{}) {
		done <- args
	}

	emitter.Once(eventName, listener)

	args := <-done
	return args, err
}
