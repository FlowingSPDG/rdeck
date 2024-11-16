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
	vMixOutput := vMixConnection.ToOutput()

	// raspi related
	raspiAdapter := raspi.NewAdaptor()
	ledDriver1 := gpio.NewLedDriver(raspiAdapter, "11") // GPIO:17
	ledOutput1 := led.NewLEDOutput(ledDriver1)
	ledDriver2 := gpio.NewLedDriver(raspiAdapter, "13") // GPIO:27
	ledOutput2 := led.NewLEDOutput(ledDriver2)
	ledDriver3 := gpio.NewLedDriver(raspiAdapter, "15") // GPIO:22
	ledOutput3 := led.NewLEDOutput(ledDriver3)

	buttonDriver := gpio.NewButtonDriver(raspiAdapter, "37")
	buttonInput := button.NewButtonInput(buttonDriver)

	// determiner/logic

	// 1: vMix Tally[PGM] -> Pin 11 LED
	tallyDeterminer := determiner.NewvMixTallyDeterminer(1)
	vMixPGMTallyConnector := vmix.NewVMixTallyConnector(vMixConnection, ledOutput1, tallyDeterminer, vmix.VMixTallyConnectorSettings{
		Target: vmix.Program,
	})
	rd.Add(ctx, vMixPGMTallyConnector)

	// 2: vMix Activator[InputPreview 1 1] -> Pin 13 LED
	previewTallyActivatorDeterminer := determiner.NewVMixActivatorDeterminer("InputPreview", 1, 1)
	previewTallyActivatorConnector := vmix.NewVMixActivatorConnector(vMixConnection, ledOutput2, previewTallyActivatorDeterminer)
	rd.Add(ctx, previewTallyActivatorConnector)

	// 2: vMix Activator[InputPlaying 1 1] -> Pin 15 LED
	inputPlayingActivatorDeterminer := determiner.NewVMixActivatorDeterminer("InputPlaying", 1, 1)
	inputPlayingActivatorConnector := vmix.NewVMixActivatorConnector(vMixConnection, ledOutput3, inputPlayingActivatorDeterminer)
	rd.Add(ctx, inputPlayingActivatorConnector)

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
				log.Println("Failed to start RDeck. ", err)
			}
		})
	}

	robot := gobot.NewRobot("RDeck",
		[]gobot.Connection{raspiAdapter},
		[]gobot.Device{ledDriver1, ledDriver2, ledDriver3, buttonDriver},
		work,
	)

	if err := robot.Start(); err != nil {
		panic(err)
	}
}
