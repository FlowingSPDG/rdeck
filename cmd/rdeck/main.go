package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/FlowingSPDG/rdeck"
	"github.com/FlowingSPDG/rdeck/connection"
	"github.com/FlowingSPDG/rdeck/connector"
	"github.com/FlowingSPDG/rdeck/determiner"
	"github.com/FlowingSPDG/rdeck/device/gpio/outputs/led"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	vMixConnection := connection.NewVMixConnection("192.168.1.10")
	vmixTallyInput := vMixConnection.ToTallyInput()

	raspiAdapter := raspi.NewAdaptor()
	ledDriver := gpio.NewLedDriver(raspiAdapter, "7")
	ledTallyOutput := led.NewLEDOutput(ledDriver)

	dt := determiner.NewvMixTallyDeterminer(1)
	vMixTallyConnector := connector.NewVMixTallyConnector(vmixTallyInput, ledTallyOutput, dt)
	rd := rdeck.NewRDeck(vMixTallyConnector)

	go func() {
		if err := vMixConnection.Start(ctx); err != nil {
			panic(err)
		}
	}()

	work := func() {
		gobot.After(500*time.Millisecond, func() {
			log.Println("STARTING!!")
			if err := rd.Start(ctx); err != nil {
				panic(err)
			}
			log.Println("END!!")
		})
	}

	robot := gobot.NewRobot("blinkBot",
		[]gobot.Connection{raspiAdapter},
		[]gobot.Device{ledDriver},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
