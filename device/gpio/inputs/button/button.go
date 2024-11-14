package button

import (
	"context"
	"log"

	"github.com/FlowingSPDG/rdeck/input"
	"gobot.io/x/gobot/v2/drivers/gpio"
)

var _ input.Input[bool] = (*buttonInput)(nil)

type buttonInput struct {
	driver *gpio.ButtonDriver
}

func NewButtonInput(driver *gpio.ButtonDriver) input.Input[bool] {
	return &buttonInput{
		driver: driver,
	}
}

func (b *buttonInput) Listen(ctx context.Context) (data <-chan bool, err <-chan error) {
	log.Println("Listening for button input...")
	d := make(chan bool)
	e := make(chan error)
	b.driver.On("push", func(data interface{}) {
		log.Println("Button pushed")
		d <- true
	})
	b.driver.On(gpio.ButtonRelease, func(data interface{}) {
		log.Println("Button released")
		d <- false
	})
	go func() {
		<-ctx.Done()
		close(d)
		close(e)
	}()

	return d, e
}
