package usecases

import (
	"context"
)

type DeviceInterface interface {
	On(ctx context.Context, params ...interface{}) error
	Off(ctx context.Context, params ...interface{}) error
}
