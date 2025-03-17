package mediasoupgo

type Only[T any, U any] struct {
	Values map[string]T
}

type Either[T any, U any] struct {
	OnlyT *Only[T, U]
	OnlyU *Only[U, T]
}

type AppData map[string]any
