package starter

type Module interface {
	Start(stop <-chan struct{})
}
