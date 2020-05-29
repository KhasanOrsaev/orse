package usecases

import (
	"context"
	"fmt"
	"github.com/KhasanOrsaev/orse/internal/app/domain"
	"testing"
)

type TestController struct {}

func (t *TestController) Store(ctx context.Context, c domain.Controller) error { return nil }
func (t *TestController) FindByID(ctx context.Context, id int) (*domain.Controller, error) { return nil, nil }
func (t *TestController) FindByName(ctx context.Context, name string) (*domain.Controller,error) { return nil, nil }

func TestControllerUsercase_Update(t *testing.T) {
	rep := ControllerUsercase{
		Controller: &TestController{},
	}
	err := rep.Update(context.Background(), domain.Controller{})
	if err != nil {
		t.Error(err)
	}
}

func TestControllerUsercase_Update2(t *testing.T) {
	rep := ControllerUsercase{
		Controller: &TestController{},
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := rep.Update(ctx, domain.Controller{})
	if err == nil {
		t.Error("should be cancel")
	}
	fmt.Println(err)
}