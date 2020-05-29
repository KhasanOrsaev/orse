package domain

import (
	"context"
	"time"
)

type Controller struct {
	ID int
	Name string
	Address string
	State int
	Topic string
	IsActive bool
	CreatedAt time.Time
	UpdatedAt time.Time
}

type ControllerRepository interface {
	Store(ctx context.Context, c Controller) error
	FindByID(ctx context.Context, id int) (*Controller, error)
	FindByName(ctx context.Context, name string) (*Controller,error)
}

