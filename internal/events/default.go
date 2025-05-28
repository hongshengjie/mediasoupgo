package events

var defaultEmmiter = New[any]()

func AddListener(evt EventName, listener ...Listener[any]) {
	defaultEmmiter.AddListener(evt, listener...)
}

func Emit(evt EventName, data any) {
	defaultEmmiter.Emit(evt, data)
}

func EventNames() []EventName {
	return defaultEmmiter.EventNames()
}
func GetMaxListeners() int {
	return defaultEmmiter.GetMaxListeners()
}

func Listeners(evt EventName) []Listener[any] {
	return defaultEmmiter.Listeners(evt)
}
func ListenerCount(evt EventName) int {
	return defaultEmmiter.ListenerCount(evt)
}

func On(evt EventName, listener ...Listener[any]) {
	defaultEmmiter.On(evt, listener...)
}

func Once(evt EventName, listener ...Listener[any]) {
	defaultEmmiter.Once(evt, listener...)
}

func Clear() {
	defaultEmmiter.Clear()
}

func SetMaxListeners(n int) {
	defaultEmmiter.SetMaxListeners(n)
}

func RemoveAllListeners(evt EventName) bool {
	return defaultEmmiter.RemoveAllListeners(evt)
}

func Len() int {
	return defaultEmmiter.Len()
}
