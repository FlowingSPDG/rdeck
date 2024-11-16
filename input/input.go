package input

import "context"

type Input[T any] interface {
	Name() string
	Listen(ctx context.Context) (data <-chan T, err <-chan error)
}
