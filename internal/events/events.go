// Package events provides simple EventEmmiter support for Go Programming Language
package events

import (
	"fmt"
	"log"
	"log/slog"
	"reflect"
	"slices"
	"sync"
	"sync/atomic"
)

const (
	DefaultMaxListeners = 0
	EnableWarning       = false
)

type (
	EventName       string
	Listener[T any] func(arg T)
	Events[T any]   map[EventName][]Listener[T]
)

type EventEmmiter[T any] interface {
	AddListener(EventName, ...Listener[T])
	Emit(EventName, T)
	EventNames() []EventName
	GetMaxListeners() int
	ListenerCount(EventName) int
	Listeners(EventName) []Listener[T]
	On(EventName, ...Listener[T])
	Once(EventName, ...Listener[T])
	RemoveAllListeners(EventName) bool
	RemoveListener(EventName, Listener[T]) bool
	Clear()
	SetMaxListeners(int)
	Len() int
}

func (e Events[T]) CopyTo(emmiter EventEmmiter[T]) {
	for evt, listeners := range e {
		if len(listeners) > 0 {
			emmiter.AddListener(evt, listeners...)
		}
	}
}

func New[T any]() EventEmmiter[T] {
	return &emmiter[T]{maxListeners: DefaultMaxListeners, evtListeners: Events[T]{}}
}

var _ EventEmmiter[any] = &emmiter[any]{}

type emmiter[T any] struct {
	maxListeners int
	evtListeners Events[T]
	mu           sync.RWMutex
}

func (e *emmiter[T]) AddListener(evt EventName, listener ...Listener[T]) {
	if len(listener) == 0 {
		return
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	if e.evtListeners == nil {
		e.evtListeners = Events[T]{}
	}

	listeners := e.evtListeners[evt]

	if e.maxListeners > 0 && len(listeners) == e.maxListeners {
		if EnableWarning {
			log.Printf(`(events) warning: possible EventEmitter memory '
                    leak detected. %d listeners added. '
                    Use emitter.SetMaxListeners(n int) to increase limit.`, len(listeners))
		}
		return
	}

	if listeners == nil {
		listeners = make([]Listener[T], e.maxListeners)
	}

	e.evtListeners[evt] = append(listeners, listener...)
}

func (e *emmiter[T]) Emit(evt EventName, data T) {
	slog.Debug("emit", slog.String("event", string(evt)), slog.Any("value", data))
	e.mu.RLock()
	if e.evtListeners == nil {
		e.mu.RUnlock()
		return // has no listeners to emit/speak yet
	}
	listeners := e.evtListeners[evt]
	list := slices.Clone(listeners)
	e.mu.RUnlock()
	// prevent deadlock caused by two level lock
	if len(list) > 0 {
		fmt.Println(len(list))
		for _, l := range list {
			if l != nil {
				l(data)
			}
		}
	}
}

func (e *emmiter[T]) EventNames() []EventName {
	e.mu.RLock()
	defer e.mu.RUnlock()

	if e.evtListeners == nil || e.Len() == 0 {
		return nil
	}

	names := make([]EventName, 0, e.Len())
	i := 0
	for k := range e.evtListeners {
		names = append(names, k)
		i++
	}
	return names
}

func (e *emmiter[T]) GetMaxListeners() int {
	return e.maxListeners
}

func (e *emmiter[T]) ListenerCount(evt EventName) (count int) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	evtListeners := e.evtListeners[evt]

	return e.listenerCount(evtListeners)
}

func (e *emmiter[T]) listenerCount(evtListeners []Listener[T]) (count int) {
	for _, l := range evtListeners {
		if l == nil {
			continue
		}
		count++
	}

	return
}

func (e *emmiter[T]) Listeners(evt EventName) []Listener[T] {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.evtListeners == nil {
		return nil
	}
	var listeners []Listener[T]
	if evtListeners := e.evtListeners[evt]; evtListeners != nil {
		// do not pass any inactive/removed listeners(nil)
		for _, l := range evtListeners {
			if l == nil {
				continue
			}

			listeners = append(listeners, l)
		}

		if len(listeners) > 0 {
			return listeners
		}
	}

	return nil
}

func (e *emmiter[T]) On(evt EventName, listener ...Listener[T]) {
	e.AddListener(evt, listener...)
}

type oneTimelistener[T any] struct {
	evt        EventName
	emitter    *emmiter[T]
	listener   Listener[T]
	fired      int32
	executeRef Listener[T]
}

func (l *oneTimelistener[T]) execute(vals T) {
	if atomic.CompareAndSwapInt32(&l.fired, 0, 1) {
		l.listener(vals)
		go l.emitter.RemoveListener(l.evt, l.executeRef)
	}
}

func (e *emmiter[T]) Once(evt EventName, listener ...Listener[T]) {
	if len(listener) == 0 {
		return
	}

	var modifiedListeners []Listener[T]

	for _, listener := range listener {
		oneTime := &oneTimelistener[T]{
			evt:      evt,
			emitter:  e,
			listener: listener,
		}
		oneTime.executeRef = oneTime.execute
		modifiedListeners = append(modifiedListeners, oneTime.executeRef)
	}
	e.AddListener(evt, modifiedListeners...)
}

func (e *emmiter[T]) RemoveAllListeners(evt EventName) bool {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.evtListeners == nil {
		return false // has nothing to remove
	}

	if listeners, ok := e.evtListeners[evt]; ok {
		count := e.listenerCount(listeners)
		delete(e.evtListeners, evt)
		return count > 0
	}

	return false
}

func (e *emmiter[T]) RemoveListener(evt EventName, listener Listener[T]) bool {
	if e.evtListeners == nil {
		return false
	}

	if listener == nil {
		return false
	}

	e.mu.Lock()
	defer e.mu.Unlock()

	listeners := e.evtListeners[evt]

	if listeners == nil {
		return false
	}

	idx := -1
	listenerPointer := reflect.ValueOf(listener).Pointer()

	for index, item := range listeners {
		itemPointer := reflect.ValueOf(item).Pointer()
		if itemPointer == listenerPointer {
			idx = index
			break
		}
	}

	if idx < 0 {
		return false
	}

	var modifiedListeners []Listener[T] = nil

	if len(listeners) > 1 {
		modifiedListeners = append(listeners[:idx], listeners[idx+1:]...)
	}

	e.evtListeners[evt] = modifiedListeners

	return true
}

func (e *emmiter[T]) Clear() {
	e.mu.Lock()
	defer e.mu.Unlock()
	e.evtListeners = Events[T]{}
}

func (e *emmiter[T]) SetMaxListeners(n int) {
	if n < 0 {
		if EnableWarning {
			log.Printf("(events) warning: MaxListeners must be positive number, tried to set: %d", n)
			return
		}
	}
	e.maxListeners = n
}

func (e *emmiter[T]) Len() int {
	e.mu.RLock()
	defer e.mu.RUnlock()
	if e.evtListeners == nil {
		return 0
	}
	return len(e.evtListeners)
}
