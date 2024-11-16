package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/FlowingSPDG/rdeck"
	"github.com/FlowingSPDG/rdeck/connector"
	"github.com/FlowingSPDG/rdeck/device/gpio/inputs/button"
	"github.com/FlowingSPDG/rdeck/device/gpio/outputs/led"

	"gobot.io/x/gobot/v2"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

func main() {
	// context
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	// rdeck
	rd := rdeck.NewRDeck()

	// raspi related
	raspiAdapter := raspi.NewAdaptor()
	ledDriver := gpio.NewLedDriver(raspiAdapter, "7")
	ledOutput := led.NewLEDOutput(ledDriver)
	buttonDriver := gpio.NewButtonDriver(raspiAdapter, "40")
	buttonInput := button.NewButtonInput(buttonDriver)

	// determiner/logic

	ltikaConnector := connector.NewLTikaConnector(buttonInput, ledOutput)
	rd.Add(ctx, ltikaConnector)

	work := func() {
		gobot.After(500*time.Millisecond, func() {
			if err := rd.Start(ctx); err != nil {
				log.Println("Failed to start rdeck. ", err)
			}
		})
	}

	robot := gobot.NewRobot("ltika",
		[]gobot.Connection{raspiAdapter},
		[]gobot.Device{ledDriver, buttonDriver},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
