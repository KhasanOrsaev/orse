package devices

import (
	"context"
	"github.com/KhasanOrsaev/logger-client"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

type ESP8266 struct {
	
}
// On turn on controller
// params consist of {address, pin}
func (e *ESP8266) On(ctx context.Context, params ...interface{}) error {
	firmataAdaptor := firmata.NewTCPAdaptor(params[0])
	led := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {
		err := led.On()
		if err != nil {
			logger.Error("error on turn on", "", &map[string]interface{}{
				"address": params[0],
				"pin": params[1],
			},err)
		}
	}

	robot := gobot.NewRobot(
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	return robot.Start()
}

// On turn on controller
// params consist of {address, pin}
func (e *ESP8266) Off(ctx context.Context, params ...interface{}) error {
	firmataAdaptor := firmata.NewTCPAdaptor(params[0])
	led := gpio.NewLedDriver(firmataAdaptor, "2")

	work := func() {
		err := led.Off()
		if err != nil {
			logger.Error("error on turn off", "", &map[string]interface{}{
				"address": params[0],
				"pin": params[1],
			},err)
		}
	}

	robot := gobot.NewRobot(
		[]gobot.Connection{firmataAdaptor},
		[]gobot.Device{led},
		work,
	)

	return robot.Start()
}