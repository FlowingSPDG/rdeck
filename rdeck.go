package rdeck

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/connector"
	"golang.org/x/xerrors"
)

// RDeck is a vMix control panel made by Raspberry Pi 4 Model B and Go, gobot.
type RDeck interface {
	Add(ctx context.Context, connector connector.VMixTallyConnector) error
	Start(ctx context.Context) error
}

type rdeck struct {
	started bool

	// connectors
	vmixTallyConnectors []connector.VMixTallyConnector
}

func NewRDeck() RDeck {
	return &rdeck{
		vmixTallyConnectors: []connector.VMixTallyConnector{},
	}
}

func (r *rdeck) Add(ctx context.Context, connector connector.VMixTallyConnector) error {
	r.vmixTallyConnectors = append(r.vmixTallyConnectors, connector)
	if r.started {
		if err := connector.Start(ctx); err != nil {
			return xerrors.Errorf("failed to start connector: %w", err)
		}
	}
	return nil
}

func (r *rdeck) Start(ctx context.Context) error {
	log.Println("STARTING RDECK.")
	r.started = true
	for _, connector := range r.vmixTallyConnectors {
		if err := connector.Start(ctx); err != nil {
			return xerrors.Errorf("failed to start connector: %w", err)
		}
	}
	<-ctx.Done()
	return ctx.Err()
}
