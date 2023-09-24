package utils

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Workers interface {
	Do(f func() error)
	Wait() error
}

type workers struct {
	eg *errgroup.Group
}

func NewWorkers(ctx context.Context, limit int) Workers {
	eg, _ := errgroup.WithContext(ctx)
	eg.SetLimit(limit)

	return &workers{
		eg: eg,
	}
}

func (w *workers) Do(f func() error) {
	w.eg.Go(f)
}

func (w *workers) Wait() error {
	return w.eg.Wait()
}
