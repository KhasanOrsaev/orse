package storage

import (
	"context"
	"github.com/KhasanOrsaev/orse/internal/domain"
	"github.com/KhasanOrsaev/orse/internal/repository/config"
	"github.com/go-errors/errors"
)

type Controller struct {
	domain.Controller
}

func (controller *Controller) Store(ctx context.Context) error {
	ch := make(chan struct{})
	go func() {
		c := config.Config()
		if c.DBClient.NewRecord(controller) {
			c.DBClient.Create(&controller)
		} else {
			c.DBClient.Save(&controller)
		}
		ch <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return errors.Wrap(ctx.Err(), -1)
	case <-ch:
		return nil
	}
}
func (controller *Controller) FindByID(ctx context.Context, id int) (*Controller, error) {
	ch := make(chan struct{})
	go func() {
		c := config.Config()
		c.DBClient.First(&controller, id)
		ch <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), -1)
	case <-ch:
		return controller, nil
	}
}
func (controller *Controller) FindByName(ctx context.Context, name string) (*Controller,error) {
	ch := make(chan struct{})
	go func() {
		c := config.Config()
		c.DBClient.Where("name=?",name).First(&controller)
		ch <- struct{}{}
	}()
	select {
	case <-ctx.Done():
		return nil, errors.Wrap(ctx.Err(), -1)
	case <-ch:
		return controller, nil
	}
}