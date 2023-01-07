package usecases

import (
	"context"
	"github.com/KhasanOrsaev/logger-client"
	"github.com/KhasanOrsaev/orse/internal/domain"
	"github.com/go-errors/errors"
)

type ControllerUsecaseInterface interface {
	Update(ctx context.Context)(bool, error)
	TurnON(ctx context.Context, id int) error
	TurnOff(ctx context.Context, id int) error
}

type ControllerUsecase struct {
	Controller *domain.Controller
	Device     DeviceInterface
	Storage    domain.ControllerInterface
}

// Update update controller
func (c *ControllerUsecase) Update(ctx context.Context) error {
	ch := make(chan struct{})
	var err error
	go func() {
		err = c.Storage.Store(ctx, c.Controller)
		ch <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), -1)
	case <-ch:
		return err
	}
}
// TurnON turn on controller
func (c *ControllerUsecase) TurnON(ctx context.Context, id int) (err error) {
	c.Controller, err = c.Storage.FindByID(ctx, id)
	if err != nil {
		return
	}
	if c.Controller == nil {
		logger.Warning("controller not found", "", &map[string]interface{}{
			"id": id,
		}, nil)
		return errors.New("controller not found")
	}
	c.Controller.State = domain.ON_SIGNAL
	err = c.Device.On(ctx, c.Controller.State, c.Controller.Address, c.Controller.Topic)
	if err != nil {
		return
	}
	err = c.Update(ctx)
	return nil
}
// TurnOff turn off controller
func (c *ControllerUsecase) TurnOff(ctx context.Context, id int) (err error) {
	c.Controller, err = c.Storage.FindByID(ctx, id)
	if err != nil {
		return
	}
	if c.Controller == nil {
		logger.Warning("controller not found", "", &map[string]interface{}{
			"id": id,
		}, nil)
		return errors.New("controller not found")
	}
	c.Controller.State = domain.OFF_SIGNAL
	err = c.Device.Off(ctx, c.Controller.Address, c.Controller.Topic)
	if err != nil {
		return
	}
	err = c.Update(ctx)
	return nil
}

func (c *ControllerUsecase) IncreasePower(rate int) error {
	return nil
}
func (c *ControllerUsecase) DecreasePower(rate int) error {
	return nil
}
func (c *ControllerUsecase) GetLevel() int {
	return 0
}

