package vmix

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/input"
	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

func NewvMixTallyInput(v vmixtcp.Vmix) input.Input[*vmixtcp.TallyResponse] {
	return &vMixInput{
		v: v,
	}
}

type vMixInput struct {
	v vmixtcp.Vmix
}

func (v *vMixInput) Listen(ctx context.Context) (<-chan *vmixtcp.TallyResponse, <-chan error) {
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
