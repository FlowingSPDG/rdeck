package connection

import (
	"context"
	"time"

	"github.com/FlowingSPDG/rdeck/device/vmix"
	"github.com/FlowingSPDG/rdeck/input"
	"github.com/FlowingSPDG/rdeck/output"
	"golang.org/x/xerrors"

	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

type VMixConnection interface {
	Start(ctx context.Context) error
	// ToInput returns a new input.Input[*vmixtcp.TallyResponse]. Do not reuse result input.Input for multiple connections.
	ToTallyInput() input.Input[*vmixtcp.TallyResponse]
	// ToActivatorInput returns a new input.Input[*vmixtcp.ActsResponse]. Do not reuse result input.Input for multiple connections.
	ToActivatorInput() input.Input[*vmixtcp.ActsResponse]
	ToOutput() output.VMixOutput
}

func NewVMixConnection(addr string) VMixConnection {
	v := vmixtcp.New(addr)
	return &vMixConnection{
		v:   v,
		tr:  make([]chan *vmixtcp.TallyResponse, 0),
		ar:  make([]chan *vmixtcp.ActsResponse, 0),
		err: make([]chan error, 0),
	}
}

type vMixConnection struct {
	v vmixtcp.Vmix

	// channels
	tr  []chan *vmixtcp.TallyResponse
	ar  []chan *vmixtcp.ActsResponse
	err []chan error
}

// Start implements VMixConnection.
func (v *vMixConnection) Start(ctx context.Context) error {
	if err := v.v.Connect(ctx, 3*time.Second); err != nil {
		return xerrors.Errorf("failed to connect: %w", err)
	}
	if err := v.v.Subscribe(vmixtcp.EventTally, ""); err != nil {
		return xerrors.Errorf("failed to subscribe to tally: %w", err)
	}
	if err := v.v.Subscribe(vmixtcp.EventActs, ""); err != nil {
		return xerrors.Errorf("failed to subscribe to acts: %w", err)
	}
	v.v.OnActs(func(acts *vmixtcp.ActsResponse, err error) {
		if err != nil {
			for _, c := range v.err {
				c <- err
			}
			return
		}
		for _, c := range v.ar {
			c <- acts
		}
	})
	v.v.OnTally(func(tally *vmixtcp.TallyResponse, err error) {
		if err != nil {
			for _, c := range v.err {
				c <- err
			}
			return
		}
		for _, c := range v.tr {
			c <- tally
		}
	})
	if err := v.v.Run(ctx); err != nil {
		return xerrors.Errorf("unknown error on running vmix: %w", err)
	}
	return nil
}

// ToInput implements VMixConnection.
func (v *vMixConnection) ToTallyInput() input.Input[*vmixtcp.TallyResponse] {
	c := make(chan *vmixtcp.TallyResponse)
	e := make(chan error)
	v.tr = append(v.tr, c)
	v.err = append(v.err, e)
	return vmix.NewvMixTallyInput(c, e)
}

func (v *vMixConnection) ToActivatorInput() input.Input[*vmixtcp.ActsResponse] {
	c := make(chan *vmixtcp.ActsResponse)
	e := make(chan error)
	v.ar = append(v.ar, c)
	v.err = append(v.err, e)

	return vmix.NewvMixActivatorInput(c, e)
}

// ToOutput implements VMixConnection.
func (v *vMixConnection) ToOutput() output.VMixOutput {
	return vmix.NewVMixOutput(v.v)
}
