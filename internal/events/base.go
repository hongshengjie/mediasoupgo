package events

type Event1[T any] struct {
	Arg1 T
}

func NewEvent1[T any](arg1 T) *Event1[T] {
	return &Event1[T]{Arg1: arg1}
}

type Event2[T, X any] struct {
	Arg1 T
	Arg2 X
}

func NewEvent2[T, X any](arg1 T, arg2 X) *Event2[T, X] {
	return &Event2[T, X]{
		Arg1: arg1,
		Arg2: arg2,
	}
}

type Event3[T, X, Y any] struct {
	Arg1 T
	Arg2 X
	Arg3 Y
}

func NewEvent3[T, X, Y any](arg1 T, arg2 X, arg3 Y) *Event3[T, X, Y] {
	return &Event3[T, X, Y]{
		Arg1: arg1,
		Arg2: arg2,
		Arg3: arg3,
	}
}

type Event4[T, X, Y, Z any] struct {
	Arg1 T
	Arg2 X
	Arg3 Y
	Arg4 Z
}

func NewEvent4[T, X, Y, Z any](arg1 T, arg2 X, arg3 Y, arg4 Z) *Event4[T, X, Y, Z] {
	return &Event4[T, X, Y, Z]{
		Arg1: arg1,
		Arg2: arg2,
		Arg3: arg3,
		Arg4: arg4,
	}
}
