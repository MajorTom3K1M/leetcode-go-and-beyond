package emitter

import (
	"fmt"
	"slices"
	"time"
)

type Listener struct {
	ID        string
	EventName string
	Callback  func(data interface{})
	Once      bool // if true, remove after first call
	CreatedAt time.Time
}

type EventLog struct {
	EventName string
	Data      interface{}
	Timestamp time.Time
	Listeners int // number of listeners that received this event
}

type EventEmitter struct {
	listeners   map[string][]*Listener
	eventLogs   []EventLog
	listenerSeq int
	// Add fields as needed
}

func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners:   make(map[string][]*Listener),
		eventLogs:   make([]EventLog, 0),
		listenerSeq: 0,
	}
}

// On registers a listener for an event.
// Returns the listener ID (format: "L-{sequential number}").
// The same callback can be registered multiple times.
func (e *EventEmitter) On(eventName string, callback func(data interface{})) string {
	e.listenerSeq++
	listenerID := fmt.Sprintf("L-%d", e.listenerSeq)

	newListener := &Listener{
		ID:        listenerID,
		EventName: eventName,
		Callback:  callback,
		Once:      false,
		CreatedAt: time.Now(),
	}
	e.listeners[eventName] = append(e.listeners[eventName], newListener)

	return listenerID
}

// Once registers a listener that will be removed after it fires once.
// Returns the listener ID.
func (e *EventEmitter) Once(eventName string, callback func(data interface{})) string {
	e.listenerSeq++
	listenerID := fmt.Sprintf("L-%d", e.listenerSeq)

	newListener := &Listener{
		ID:        listenerID,
		EventName: eventName,
		Callback:  callback,
		Once:      true,
		CreatedAt: time.Now(),
	}
	e.listeners[eventName] = append(e.listeners[eventName], newListener)

	return listenerID
}

// Off removes a listener by ID.
// Returns true if the listener was found and removed.
func (e *EventEmitter) Off(listenerID string) bool {
	for key, listeners := range e.listeners {
		for i, listener := range listeners {
			if listener.ID == listenerID {
				listeners = append(listeners[:i], listeners[i+1:]...)

				if len(listeners) == 0 {
					delete(e.listeners, key)
				} else {
					e.listeners[key] = listeners
				}

				return true
			}
		}
	}

	return false
}

// OffAll removes all listeners for an event.
// Returns the number of listeners removed.
// If eventName is empty, removes ALL listeners for ALL events.
func (e *EventEmitter) OffAll(eventName string) int {
	removed := 0
	if eventName == "" {
		for event, listeners := range e.listeners {
			removed += len(listeners)
			delete(e.listeners, event)
		}

		return removed
	}

	if _, exists := e.listeners[eventName]; !exists {
		return 0
	}

	removed = len(e.listeners[eventName])

	delete(e.listeners, eventName)

	return removed
}

// Emit triggers an event with the given data.
// Calls all registered listeners for that event.
// Returns the number of listeners called.
// Logs the event in eventLogs.
// "Once" listeners should be removed after being called.
func (e *EventEmitter) Emit(eventName string, data interface{}) int {
	listeners, exists := e.listeners[eventName]
	if !exists {
		e.eventLogs = append(e.eventLogs, EventLog{
			EventName: eventName,
			Data:      data,
			Timestamp: time.Now(),
			Listeners: 0,
		})

		return 0
	}

	remaining := []*Listener{}
	for _, l := range listeners {
		l.Callback(data)
		if !l.Once {
			remaining = append(remaining, l)
		}
	}
	e.listeners[eventName] = remaining

	e.eventLogs = append(e.eventLogs, EventLog{
		EventName: eventName,
		Data:      data,
		Timestamp: time.Now(),
		Listeners: len(listeners),
	})

	return len(listeners)
}

// ListenerCount returns the number of listeners for an event.
// If eventName is empty, returns total listeners for all events.
func (e *EventEmitter) ListenerCount(eventName string) int {
	total := 0

	if eventName == "" {
		for _, listeners := range e.listeners {
			total += len(listeners)
		}

		return total
	}

	if listeners, exists := e.listeners[eventName]; exists {
		return len(listeners)
	}

	return 0
}

// EventNames returns all event names that have at least one listener.
func (e *EventEmitter) EventNames() []string {
	result := []string{}
	for eventName, listeners := range e.listeners {
		if len(listeners) > 0 {
			result = append(result, eventName)
		}
	}
	return result
}

// GetListeners returns all listeners for an event.
// Returns empty slice if no listeners.
func (e *EventEmitter) GetListeners(eventName string) []*Listener {
	if listeners, exists := e.listeners[eventName]; exists {
		return listeners
	}
	return []*Listener{}
}

// GetEventLogs returns all event logs, sorted by timestamp ascending.
func (e *EventEmitter) GetEventLogs() []EventLog {
	result := make([]EventLog, len(e.eventLogs))
	copy(result, e.eventLogs)
	slices.SortFunc(result, func(a, b EventLog) int {
		return a.Timestamp.Compare(b.Timestamp)
	})

	return result
}

// GetEventLogsByName returns all logs for a specific event name.
// Sorted by timestamp ascending.
func (e *EventEmitter) GetEventLogsByName(eventName string) []EventLog {
	result := []EventLog{}
	for _, eventLog := range e.eventLogs {
		if eventLog.EventName == eventName {
			result = append(result, eventLog)
		}
	}
	return result
}

// ClearEventLogs removes all event logs.
// Returns the number of logs cleared.
func (e *EventEmitter) ClearEventLogs() int {
	removed := len(e.eventLogs)
	e.eventLogs = make([]EventLog, 0)
	return removed
}
