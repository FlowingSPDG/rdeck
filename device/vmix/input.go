package vmix

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/input"
	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

func NewvMixTallyInput(v vmixtcp.Vmix) input.Input[*vmixtcp.TallyResponse] {
	return &vMixTallyInput{
		v: v,
	}
}

type vMixTallyInput struct {
	v vmixtcp.Vmix
}

func (v *vMixTallyInput) Listen(ctx context.Context) (<-chan *vmixtcp.TallyResponse, <-chan error) {
	d := make(chan *vmixtcp.TallyResponse)
	e := make(chan error)
	if err := v.v.Subscribe(vmixtcp.EventTally, ""); err != nil {
		log.Println("failed to subscribe... ", err)
	}
	v.v.OnTally(func(tr *vmixtcp.TallyResponse, err error) {
		log.Println("received tally:", tr)
		d <- tr
		if err != nil {
			e <- err
		}
	})
	go func() {
		<-ctx.Done()
		close(d)
		close(e)
	}()

	return d, e
}

func NewvMixActivatorInput(v vmixtcp.Vmix) input.Input[*vmixtcp.ActsResponse] {
	return &vMixActivatorInput{
		v: v,
	}
}

type vMixActivatorInput struct {
	v vmixtcp.Vmix
}

func (v *vMixActivatorInput) Listen(ctx context.Context) (data <-chan *vmixtcp.ActsResponse, err <-chan error) {
	d := make(chan *vmixtcp.ActsResponse)
	e := make(chan error)
	if err := v.v.Subscribe(vmixtcp.EventActs, ""); err != nil {
		log.Println("failed to subscribe... ", err)
	}
	v.v.OnActs(func(acts *vmixtcp.ActsResponse, err error) {
		log.Println("received activator:", acts)
		d <- acts
		if err != nil {
			e <- err
		}
	})
	go func() {
		<-ctx.Done()
		close(d)
		close(e)
	}()

	return d, e
}
