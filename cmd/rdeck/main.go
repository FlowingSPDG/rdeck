package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/FlowingSPDG/rdeck"
	"github.com/FlowingSPDG/rdeck/connection"
	"github.com/FlowingSPDG/rdeck/connector/vmix"
	"github.com/FlowingSPDG/rdeck/determiner"
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

	// vMix related
	vMixConnectionPool := connection.NewvMixConnectionPool()
	vMixConnection := vMixConnectionPool.AddNew("192.168.1.10")
	vMixTallyInput := vMixConnection.ToTallyInput()
	vMixActivatorInput := vMixConnection.ToActivatorInput()
	vMixOutput := vMixConnection.ToOutput()

	// raspi related
	raspiAdapter := raspi.NewAdaptor()
	ledDriver1 := gpio.NewLedDriver(raspiAdapter, "7")
	ledOutput1 := led.NewLEDOutput(ledDriver1)
	ledDriver2 := gpio.NewLedDriver(raspiAdapter, "8")
	ledOutput2 := led.NewLEDOutput(ledDriver2)

	buttonDriver := gpio.NewButtonDriver(raspiAdapter, "40")
	buttonInput := button.NewButtonInput(buttonDriver)

	// determiner/logic

	// 1: vMix Tally -> LED
	tallyDeterminer := determiner.NewvMixTallyDeterminer(1)
	vMixTallyConnector := vmix.NewVMixTallyConnector(vMixTallyInput, ledOutput1, tallyDeterminer)
	rd.Add(ctx, vMixTallyConnector)

	// 2: vMix Activator -> LED
	activatorDeterminer := determiner.NewVMixActivatorDeterminer("InputPlaying", 1, 1)
	vMixActivatorConnector := vmix.NewVMixActivatorConnector(vMixActivatorInput, ledOutput2, activatorDeterminer)
	rd.Add(ctx, vMixActivatorConnector)

	// 3: Button -> vMix Function
	vMixSendFunctionConnector := vmix.NewSendFunction(buttonInput, vMixOutput, "Cut", "Input=1")
	rd.Add(ctx, vMixSendFunctionConnector)

	go func() {
		if err := vMixConnection.Start(ctx); err != nil {
			panic(err)
		}
	}()

	work := func() {
		gobot.After(500*time.Millisecond, func() {
			if err := rd.Start(ctx); err != nil {
				log.Println("Failed to start rdeck.] ", err)
			}
		})
	}

	robot := gobot.NewRobot("rdeck",
		[]gobot.Connection{raspiAdapter},
		[]gobot.Device{ledDriver1, ledDriver2, buttonDriver},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
