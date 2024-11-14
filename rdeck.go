package rdeck

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/connector"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"
)

// RDeck is a vMix control panel made by Raspberry Pi 4 Model B and Go, gobot.
type RDeck interface {
	Add(ctx context.Context, connector connector.Connector) error
	Start(ctx context.Context) error
}

type rdeck struct {
	started bool

	connectors []connector.Connector
}

func NewRDeck() RDeck {
	return &rdeck{
		connectors: []connector.Connector{},
	}
}

func (r *rdeck) Add(ctx context.Context, connector connector.Connector) error {
	r.connectors = append(r.connectors, connector)
	if r.started {
		go func() {}()
		if err := connector.Start(ctx); err != nil {
			//TODO: handle
			log.Println("Failed to start connector:", err)
		}
	}
	return nil
}

func (r *rdeck) Start(ctx context.Context) error {
	log.Println("STARTING RDECK.")
	r.started = true

	eg := errgroup.Group{}
	for _, connector := range r.connectors {
		eg.Go(func() error {
			if err := connector.Start(ctx); err != nil {
				return xerrors.Errorf("failed to start connector: %w", err)
			}
			return nil
		})
	}
	if err := eg.Wait(); err != nil {
		return xerrors.Errorf("failed to start connectors: %w", err)
	}
	<-ctx.Done()
	return ctx.Err()
}
