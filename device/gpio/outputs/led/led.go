package led

import (
	"log"

	"github.com/FlowingSPDG/rdeck/output"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"golang.org/x/xerrors"
)

var _ output.Digital = (*ledOutput)(nil)

type ledOutput struct {
	driver *gpio.LedDriver
}

func NewLEDOutput(driver *gpio.LedDriver) output.Digital {
	return &ledOutput{
		driver: driver,
	}
}

func (l *ledOutput) Name() string {
	return l.driver.Name()
}

func (l *ledOutput) On() error {
	log.Println("ACTIVATING LED...")
	if err := l.driver.On(); err != nil {
		return xerrors.Errorf("failed to turn on LED interface: %w", err)
	}
	return nil
}

// Inactive implements tally.TallyLight.
func (l *ledOutput) Off() error {
	log.Println("DEACTIVATING LED...")
	if err := l.driver.Off(); err != nil {
		return xerrors.Errorf("failed to turn on LED interface: %w", err)
	}
	return nil
}
