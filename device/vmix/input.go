package vmix

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/input"
	vmixtcp "github.com/FlowingSPDG/vmix-go/tcp"
)

func NewvMixTallyInput(tr <-chan *vmixtcp.TallyResponse, err <-chan error) input.Input[*vmixtcp.TallyResponse] {
	return &vMixTallyInput{
		tr:  tr,
		err: err,
	}
}

type vMixTallyInput struct {
	tr  <-chan *vmixtcp.TallyResponse
	err <-chan error
}

func (v *vMixTallyInput) Listen(ctx context.Context) (<-chan *vmixtcp.TallyResponse, <-chan error) {
	d := make(chan *vmixtcp.TallyResponse)
	e := make(chan error)

	go func() {
		for {
			select {

			case tr := <-v.tr:
				log.Println("received tally:", tr)
				d <- tr

			case err := <-v.err:
				if err != nil {
					e <- err
				}

			case <-ctx.Done():
				close(d)
				close(e)
			}
		}
	}()

	return d, e
}

func (v *vMixTallyInput) Name() string {
	return "vMixTallyInput"
}

func NewvMixActivatorInput(ar <-chan *vmixtcp.ActsResponse, err <-chan error) input.Input[*vmixtcp.ActsResponse] {
	return &vMixActivatorInput{
		ar:  ar,
		err: err,
	}
}

type vMixActivatorInput struct {
	ar  <-chan *vmixtcp.ActsResponse
	err <-chan error
}

func (v *vMixActivatorInput) Listen(ctx context.Context) (data <-chan *vmixtcp.ActsResponse, err <-chan error) {
	d := make(chan *vmixtcp.ActsResponse)
	e := make(chan error)

	go func() {
		for {
			select {

			case ar := <-v.ar:
				log.Println("received acts:", ar)
				d <- ar

			case err := <-v.err:
				if err != nil {
					e <- err
				}

			case <-ctx.Done():
				close(d)
				close(e)
			}
		}
	}()

	return d, e
}

func (v *vMixActivatorInput) Name() string {
	return "vMixActivatorInput"
}
