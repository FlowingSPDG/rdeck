package led

import (
	"log"

	"github.com/FlowingSPDG/rdeck/output"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"golang.org/x/xerrors"
)

var _ output.Tally = (*ledOutput)(nil)

type TallyData struct{}

type ledOutput struct {
	driver *gpio.LedDriver
}

func NewLEDOutput(driver *gpio.LedDriver) output.Tally {
	return &ledOutput{
		driver: driver,
	}
}

// Active implements tally.TallyLight.
func (l *ledOutput) Active() error {
	log.Println("ACTIVATING LED...")
	if err := l.driver.On(); err != nil {
		return xerrors.Errorf("failed to turn on LED interface: %w", err)
	}
	return nil
}

// Inactive implements tally.TallyLight.
func (l *ledOutput) Inactive() error {
	log.Println("DEACTIVATING LED...")
	if err := l.driver.Off(); err != nil {
		return xerrors.Errorf("failed to turn on LED interface: %w", err)
	}
	return nil
}

// Preview implements tally.TallyLight.
func (l *ledOutput) Preview() error {
	log.Println("DEACTIVATING LED...")
	// Since LED only supports On/Off, lets leave preview tally.
	if err := l.driver.Off(); err != nil {
		return xerrors.Errorf("failed to turn on LED interface: %w", err)
	}
	return nil
}
