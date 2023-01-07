package usecases

import (
	"context"
	"testing"

	"github.com/KhasanOrsaev/orse/internal/domain"
)

type TestController struct{}
type TestDevice struct{}

var rep ControllerUsecase

func init() {
	rep = ControllerUsecase{
		Storage: &TestController{},
		Device:  &TestDevice{},
	}
}

func (t *TestController) Store(ctx context.Context, c *domain.Controller) error { return nil }
func (t *TestController) FindByID(ctx context.Context, id int) (*domain.Controller, error) {
	return &domain.Controller{ID: id}, nil
}
func (t *TestController) FindByName(ctx context.Context, name string) (*domain.Controller, error) {
	return nil, nil
}
func (t *TestDevice) SendSignal(ctx context.Context, signal int, params ...interface{}) error {
	return nil
}
func (t *TestDevice) On(ctx context.Context, params ...interface{}) error  { return nil }
func (t *TestDevice) Off(ctx context.Context, params ...interface{}) error { return nil }

func TestControllerUsercase_Update(t *testing.T) {
	err := rep.Update(context.Background())
	if err != nil {
		t.Error(err)
	}
}

func TestControllerUsecase_TurnON(t *testing.T) {
	err := rep.TurnON(context.Background(), 1)
	if err != nil {
		t.Error(err)
	}
}
