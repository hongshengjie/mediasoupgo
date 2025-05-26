package events

import (
	"fmt"
	"testing"
	"time"
)

var testEvents = Events[[]any]{
	"user_created": []Listener[[]any]{
		func(payload []any) {
			fmt.Printf("A new User just created!\n")
		},
		func(payload []any) {
			fmt.Printf("A new User just created, *from second event listener\n")
		},
	},
	"user_joined": []Listener[[]any]{func(payload []any) {
		user := payload[0].(string)
		room := payload[1].(string)
		fmt.Printf("%s joined to room: %s\n", user, room)
	}},
	"user_left": []Listener[[]any]{func(payload []any) {
		user := payload[0].(string)
		room := payload[1].(string)
		fmt.Printf("%s left from the room: %s\n", user, room)
	}},
}

func createUser(user string) {
	e.Emit("user_created", []any{user})
}

func joinUserTo(user string, room string) {
	e.Emit("user_joined", []any{user, room})
}

func leaveFromRoom(user string, room string) {
	e.Emit("user_left", []any{user, room})
}

var e = New[[]any]()

func ExampleEvents() {

	// regiter our events to the default event emmiter
	for evt, listeners := range testEvents {
		e.On(evt, listeners...)
	}

	user := "user1"
	room := "room1"

	createUser(user)
	joinUserTo(user, room)
	leaveFromRoom(user, room)

	// Output:
	// A new User just created!
	// A new User just created, *from second event listener
	// user1 joined to room: room1
	// user1 left from the room: room1
}

func TestEvents(t *testing.T) {
	e := New[[]any]()
	expectedPayload := "this is my payload"

	e.On("my_event", func(payload []any) {

		if s, ok := payload[0].(string); !ok {
			t.Fatalf("Payload is not the correct type, got: %#v", payload)
		} else if s != expectedPayload {
			t.Fatalf("Eexpected %s, got: %s", expectedPayload, s)
		}
	})

	e.Emit("my_event", []any{expectedPayload})
	if e.Len() != 1 {
		t.Fatalf("Length of the events is: %d, while expecting: %d", e.Len(), 1)
	}

	if e.Len() != 1 {
		t.Fatalf("Length of the listeners is: %d, while expecting: %d", e.ListenerCount("my_event"), 1)
	}

	e.RemoveAllListeners("my_event")
	if e.Len() != 0 {
		t.Fatalf("Length of the events is: %d, while expecting: %d", e.Len(), 0)
	}

	if e.Len() != 0 {
		t.Fatalf("Length of the listeners is: %d, while expecting: %d", e.ListenerCount("my_event"), 0)
	}
}

func TestEventsOnce(t *testing.T) {
	// on default
	Clear()

	var count = 0
	Once("my_event", func(payload any) {
		if count > 0 {
			t.Fatalf("Once's listener fired more than one time! count: %d", count)
		}
		if payload.(string) != "foo" {
			t.Fatalf("Once's listener payload is incorrect: %+v", payload)
		}
		count++
	})

	if l := ListenerCount("my_event"); l != 1 {
		t.Fatalf("Real  event's listeners should be: %d but has: %d", 1, l)
	}

	if l := len(Listeners("my_event")); l != 1 {
		t.Fatalf("Real  event's listeners (from Listeners) should be: %d but has: %d", 1, l)
	}

	for i := 0; i < 10; i++ {
		Emit("my_event", "foo")
	}

	time.Sleep(10 * time.Millisecond)

	if l := ListenerCount("my_event"); l > 0 {
		t.Fatalf("Real event's listeners length count should be: %d but has: %d", 0, l)
	}

	if l := len(Listeners("my_event")); l > 0 {
		t.Fatalf("Real event's listeners length count ( from Listeners) should be: %d but has: %d", 0, l)
	}

}

func TestRemoveListener(t *testing.T) {
	// on default
	e := New[[]any]()

	var count = 0
	listener := func(payload []any) {
		if count > 1 {
			t.Fatal("Event listener should be removed")
		}

		count++
	}

	e.AddListener("my_event", listener)
	e.AddListener("my_event", func(payload []any) {})
	e.AddListener("another_event", func(payload []any) {})

	e.Emit("my_event", []any{""})

	if e.RemoveListener("my_event", listener) != true {
		t.Fatal("Should return 'true' when removes found listener")
	}

	if e.RemoveListener("foo_bar", listener) != false {
		t.Fatal("Should return 'false' when removes nothing")
	}

	if e.Len() != 2 {
		t.Fatal("Length of all listeners must be 2")
	}

	if e.ListenerCount("my_event") != 1 {
		t.Fatal("Length of 'my_event' event listeners must be 1")
	}

	e.Emit("my_event", []any{""})
}

type W1 struct {
	Foo string
}
type W2 struct {
	Bar int
}
type W3 struct {
	XXX any
}

type WebEvent struct {
	Event1 *Event1[W1]
	Event2 *Event2[W1, W3]
	Event3 *Event3[W1, W2, W3]
}

var WebEvents = Events[WebEvent]{
	"event1": []Listener[WebEvent]{func(arg WebEvent) {
		fmt.Println(arg.Event1)
	}},

	"event2": []Listener[WebEvent]{func(arg WebEvent) {
		fmt.Println(arg.Event2)
	}},

	"event3": []Listener[WebEvent]{func(arg WebEvent) {
		fmt.Println(arg.Event3)
	}},
}

func TestGeneric(t *testing.T) {
	emiter := New[WebEvent]()
	WebEvents.CopyTo(emiter)
	emiter.Once("event1", func(arg WebEvent) {
		fmt.Println("once", arg.Event1)
	})

	emiter.Emit("event2", WebEvent{Event2: NewEvent2(W1{Foo: "ddd"}, W3{XXX: "ddfd"})})
	emiter.Emit("event1", WebEvent{Event1: NewEvent1(W1{Foo: "ddd"})})
	emiter.Emit("event1", WebEvent{Event1: NewEvent1(W1{Foo: "xxx"})})
	emiter.Emit("event3", WebEvent{Event3: NewEvent3(W1{Foo: "ddd"}, W2{Bar: 22}, W3{XXX: "ddfd"})})
}
