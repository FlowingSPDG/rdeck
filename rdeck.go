package rdeck

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/connector"
)

// RDeck is a vMix control panel made by Raspberry Pi 4 Model B and Go, gobot.
type RDeck interface {
	Start(ctx context.Context) error
	// TODO:
	// AddConnector(in input.Input, out output.Output) error // Adds connector with pre-configured input/output.
}

type rdeck struct {
	// input/output pool
	// TODO: I/O の接続をプールするフィールドの作成
	// 途中追加も可能にしたい

	// connectors
	// TODO: define with slice
	// スライスとして定義して、途中で追加したり停止・再起動を出来る様にしたい
	vmixTallyConnector connector.VMixTallyConnector
}

func NewRDeck(vmixTallyConnector connector.VMixTallyConnector) RDeck {
	return &rdeck{
		vmixTallyConnector: vmixTallyConnector,
	}
}

func (r *rdeck) Start(ctx context.Context) error {
	log.Println("STARTING RDECK.")
	return r.vmixTallyConnector.Start(ctx)
}
