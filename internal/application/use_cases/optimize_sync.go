package usecases

import (
	"context"
	"time"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
)

type OptimizeSyncUseCase interface {
	Execute(ctx context.Context) error
}

type OptimizeSyncUseCaseImpl struct {}

func NewOptimizeSyncUseCaseImpl() OptimizeSyncUseCaseImpl {
	return OptimizeSyncUseCaseImpl{}
}

func (uc *OptimizeSyncUseCaseImpl) Execute(ctx context.Context) ([]dtos.OptimizationRequestOutput, error) {
	time.Sleep(time.Millisecond * 10)

	return make([]dtos.OptimizationRequestOutput, 0), nil
}
