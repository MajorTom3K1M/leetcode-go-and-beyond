package emitter

import (
	"testing"
	"time"
)

func TestNewEventEmitter(t *testing.T) {
	e := NewEventEmitter()
	if e == nil {
		t.Fatal("emitter should not be nil")
	}
	if e.ListenerCount("") != 0 {
		t.Error("new emitter should have no listeners")
	}
}

func TestOnAndEmit(t *testing.T) {
	e := NewEventEmitter()

	received := []interface{}{}
	id := e.On("message", func(data interface{}) {
		received = append(received, data)
	})

	if id != "L-1" {
		t.Errorf("expected L-1, got %s", id)
	}

	// Emit event
	count := e.Emit("message", "hello")
	if count != 1 {
		t.Errorf("expected 1 listener called, got %d", count)
	}
	if len(received) != 1 || received[0] != "hello" {
		t.Error("callback should receive data")
	}

	// Emit again
	e.Emit("message", "world")
	if len(received) != 2 || received[1] != "world" {
		t.Error("callback should receive second emit")
	}

	// Emit non-existent event
	count = e.Emit("other", "data")
	if count != 0 {
		t.Error("non-existent event should have 0 listeners")
	}
}

func TestMultipleListeners(t *testing.T) {
	e := NewEventEmitter()

	calls := make([]int, 3)
	e.On("event", func(data interface{}) { calls[0]++ })
	e.On("event", func(data interface{}) { calls[1]++ })
	e.On("event", func(data interface{}) { calls[2]++ })

	count := e.Emit("event", nil)
	if count != 3 {
		t.Errorf("expected 3 listeners, got %d", count)
	}
	if calls[0] != 1 || calls[1] != 1 || calls[2] != 1 {
		t.Error("all callbacks should be called once")
	}
}

func TestOnce(t *testing.T) {
	e := NewEventEmitter()

	callCount := 0
	e.Once("event", func(data interface{}) {
		callCount++
	})

	if e.ListenerCount("event") != 1 {
		t.Error("should have 1 listener before emit")
	}

	e.Emit("event", nil)
	if callCount != 1 {
		t.Error("once listener should be called")
	}
	if e.ListenerCount("event") != 0 {
		t.Error("once listener should be removed after emit")
	}

	// Emit again - should not call
	e.Emit("event", nil)
	if callCount != 1 {
		t.Error("once listener should not be called again")
	}
}

func TestOnceMixedWithOn(t *testing.T) {
	e := NewEventEmitter()

	regularCalls := 0
	onceCalls := 0

	e.On("event", func(data interface{}) { regularCalls++ })
	e.Once("event", func(data interface{}) { onceCalls++ })
	e.On("event", func(data interface{}) { regularCalls++ })

	// First emit
	count := e.Emit("event", nil)
	if count != 3 {
		t.Errorf("first emit should call 3 listeners, got %d", count)
	}
	if regularCalls != 2 || onceCalls != 1 {
		t.Error("incorrect call counts after first emit")
	}

	// Second emit
	count = e.Emit("event", nil)
	if count != 2 {
		t.Errorf("second emit should call 2 listeners, got %d", count)
	}
	if regularCalls != 4 || onceCalls != 1 {
		t.Error("once listener should not be called on second emit")
	}
}

func TestOff(t *testing.T) {
	e := NewEventEmitter()

	callCount := 0
	id := e.On("event", func(data interface{}) {
		callCount++
	})

	e.Emit("event", nil)
	if callCount != 1 {
		t.Error("listener should be called before Off")
	}

	// Remove listener
	if !e.Off(id) {
		t.Error("Off should return true for existing listener")
	}

	e.Emit("event", nil)
	if callCount != 1 {
		t.Error("listener should not be called after Off")
	}

	// Remove non-existent
	if e.Off("L-999") {
		t.Error("Off should return false for non-existent listener")
	}
}

func TestOffAll(t *testing.T) {
	e := NewEventEmitter()

	e.On("event1", func(data interface{}) {})
	e.On("event1", func(data interface{}) {})
	e.On("event2", func(data interface{}) {})

	// Remove all for event1
	removed := e.OffAll("event1")
	if removed != 2 {
		t.Errorf("expected 2 removed, got %d", removed)
	}
	if e.ListenerCount("event1") != 0 {
		t.Error("event1 should have no listeners")
	}
	if e.ListenerCount("event2") != 1 {
		t.Error("event2 should still have 1 listener")
	}

	// Remove all for all events
	e.On("event1", func(data interface{}) {})
	e.On("event3", func(data interface{}) {})

	removed = e.OffAll("")
	if removed != 3 {
		t.Errorf("expected 3 removed, got %d", removed)
	}
	if e.ListenerCount("") != 0 {
		t.Error("all listeners should be removed")
	}
}

func TestListenerCount(t *testing.T) {
	e := NewEventEmitter()

	e.On("event1", func(data interface{}) {})
	e.On("event1", func(data interface{}) {})
	e.On("event2", func(data interface{}) {})

	if e.ListenerCount("event1") != 2 {
		t.Error("event1 should have 2 listeners")
	}
	if e.ListenerCount("event2") != 1 {
		t.Error("event2 should have 1 listener")
	}
	if e.ListenerCount("event3") != 0 {
		t.Error("event3 should have 0 listeners")
	}
	if e.ListenerCount("") != 3 {
		t.Error("total should be 3 listeners")
	}
}

func TestEventNames(t *testing.T) {
	e := NewEventEmitter()

	e.On("zebra", func(data interface{}) {})
	e.On("apple", func(data interface{}) {})
	e.On("apple", func(data interface{}) {})

	names := e.EventNames()
	if len(names) != 2 {
		t.Errorf("expected 2 event names, got %d", len(names))
	}

	nameMap := make(map[string]bool)
	for _, n := range names {
		nameMap[n] = true
	}
	if !nameMap["zebra"] || !nameMap["apple"] {
		t.Error("should have zebra and apple events")
	}
}

func TestGetListeners(t *testing.T) {
	e := NewEventEmitter()

	id1 := e.On("event", func(data interface{}) {})
	id2 := e.Once("event", func(data interface{}) {})

	listeners := e.GetListeners("event")
	if len(listeners) != 2 {
		t.Fatalf("expected 2 listeners, got %d", len(listeners))
	}

	// Check listener properties
	idMap := make(map[string]*Listener)
	for _, l := range listeners {
		idMap[l.ID] = l
	}

	if idMap[id1].Once != false {
		t.Error("first listener should not be Once")
	}
	if idMap[id2].Once != true {
		t.Error("second listener should be Once")
	}
	if idMap[id1].EventName != "event" {
		t.Error("listener should have correct EventName")
	}

	// Non-existent event
	if len(e.GetListeners("other")) != 0 {
		t.Error("non-existent event should return empty slice")
	}
}

func TestEventLogs(t *testing.T) {
	e := NewEventEmitter()

	e.On("event1", func(data interface{}) {})
	e.On("event1", func(data interface{}) {})

	time.Sleep(10 * time.Millisecond)
	e.Emit("event1", "data1")

	time.Sleep(10 * time.Millisecond)
	e.Emit("event2", "data2") // no listeners

	time.Sleep(10 * time.Millisecond)
	e.Emit("event1", "data3")

	logs := e.GetEventLogs()
	if len(logs) != 3 {
		t.Fatalf("expected 3 logs, got %d", len(logs))
	}

	// Check order (ascending by timestamp)
	if logs[0].EventName != "event1" || logs[0].Data != "data1" {
		t.Error("first log incorrect")
	}
	if logs[0].Listeners != 2 {
		t.Errorf("first log should have 2 listeners, got %d", logs[0].Listeners)
	}
	if logs[1].EventName != "event2" || logs[1].Listeners != 0 {
		t.Error("second log incorrect")
	}
	if logs[2].EventName != "event1" || logs[2].Data != "data3" {
		t.Error("third log incorrect")
	}

	// Check timestamps are ascending
	for i := 1; i < len(logs); i++ {
		if logs[i].Timestamp.Before(logs[i-1].Timestamp) {
			t.Error("logs should be sorted by timestamp ascending")
		}
	}
}

func TestGetEventLogsByName(t *testing.T) {
	e := NewEventEmitter()

	e.On("event1", func(data interface{}) {})

	e.Emit("event1", "a")
	e.Emit("event2", "b")
	e.Emit("event1", "c")
	e.Emit("event2", "d")

	logs := e.GetEventLogsByName("event1")
	if len(logs) != 2 {
		t.Errorf("expected 2 logs for event1, got %d", len(logs))
	}
	if logs[0].Data != "a" || logs[1].Data != "c" {
		t.Error("logs should be filtered correctly")
	}

	// Non-existent event
	if len(e.GetEventLogsByName("other")) != 0 {
		t.Error("non-existent event should return empty slice")
	}
}

func TestClearEventLogs(t *testing.T) {
	e := NewEventEmitter()

	e.Emit("event", "data")
	e.Emit("event", "data")
	e.Emit("event", "data")

	cleared := e.ClearEventLogs()
	if cleared != 3 {
		t.Errorf("expected 3 cleared, got %d", cleared)
	}
	if len(e.GetEventLogs()) != 0 {
		t.Error("logs should be empty after clear")
	}
}

func TestSequentialListenerIDs(t *testing.T) {
	e := NewEventEmitter()

	id1 := e.On("a", func(data interface{}) {})
	id2 := e.Once("b", func(data interface{}) {})
	id3 := e.On("c", func(data interface{}) {})

	if id1 != "L-1" || id2 != "L-2" || id3 != "L-3" {
		t.Errorf("IDs should be sequential: got %s, %s, %s", id1, id2, id3)
	}
}

func TestEmitWithNilData(t *testing.T) {
	e := NewEventEmitter()

	var received interface{} = "not nil"
	e.On("event", func(data interface{}) {
		received = data
	})

	e.Emit("event", nil)
	if received != nil {
		t.Error("should receive nil data")
	}
}

func TestEmitOrder(t *testing.T) {
	e := NewEventEmitter()

	order := []int{}
	e.On("event", func(data interface{}) { order = append(order, 1) })
	e.On("event", func(data interface{}) { order = append(order, 2) })
	e.On("event", func(data interface{}) { order = append(order, 3) })

	e.Emit("event", nil)

	if len(order) != 3 {
		t.Fatal("all listeners should be called")
	}
	if order[0] != 1 || order[1] != 2 || order[2] != 3 {
		t.Error("listeners should be called in registration order")
	}
}

func TestListenerCreatedAt(t *testing.T) {
	e := NewEventEmitter()

	before := time.Now()
	e.On("event", func(data interface{}) {})
	after := time.Now()

	listeners := e.GetListeners("event")
	if len(listeners) != 1 {
		t.Fatal("should have 1 listener")
	}

	createdAt := listeners[0].CreatedAt
	if createdAt.Before(before) || createdAt.After(after) {
		t.Error("CreatedAt should be between before and after")
	}
}
