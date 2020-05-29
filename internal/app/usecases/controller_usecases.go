package usecases

import (
	"context"
	"github.com/KhasanOrsaev/orse/internal/app/domain"
	"github.com/go-errors/errors"
)

type ControllerUsecaseRepository interface {
	Update(ctx context.Context, controller *domain.Controller)(bool, error)
	TurnON() error
	TurnOff() error
	IncreasePower(rate int) error
	DecreasePower(rate int) error
	GetLevel() int
}

type DeviceRepository interface {
	SendSignal(ctx context.Context, controller *domain.Controller, signal string) error
}

type ControllerUsercase struct {
	Controller domain.ControllerRepository
	Device DeviceRepository
}

func (c *ControllerUsercase) Update(ctx context.Context, controller domain.Controller) error {
	ch := make(chan struct{})
	var err error
	go func() {
		err = c.Controller.Store(ctx, controller)
		ch <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), -1)
	case <-ch:
		return err
	}
}
func (c *ControllerUsercase) TurnON(ctx context.Context, id int) error {
	controller,err := c.Controller.FindByID(ctx, id)
	if err != nil {
		return err
	}
	controller.State = 1
	err = c.Device.SendSignal(ctx, controller, "on")
	if err != nil {
		return err
	}
	err = c.Update(ctx, controller)
	return nil
}
func (c *ControllerUsercase) TurnOff() error {
	return nil
}
func (c *ControllerUsercase) IncreasePower(rate int) error {
	return nil
}
func (c *ControllerUsercase) DecreasePower(rate int) error {
	return nil
}
func (c *ControllerUsercase) GetLevel() int {
	return 0
}

