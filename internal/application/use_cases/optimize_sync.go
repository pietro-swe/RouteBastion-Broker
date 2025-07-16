package usecases

import (
	"context"
	"time"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	// clientImpl "github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/clients/implementations"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
)

type OptimizeSyncUseCase interface {
	Execute(ctx context.Context, dto dtos.OptimizationRequestInput) error
}

type OptimizeSyncUseCaseDeps struct {
	DB persistence.DBProvider
}

type OptimizeSyncUseCaseImpl struct {
	deps OptimizeSyncUseCaseDeps
}

func (uc *OptimizeSyncUseCaseImpl) Execute(ctx context.Context, dto dtos.OptimizationRequestInput) ([]dtos.OptimizationRequestOutput, error) {
	// TODO: Algorithm to choose best fit
	// client := clientImpl.NewFakeRouteOptimizer()

	time.Sleep(time.Millisecond * 10)

	return make([]dtos.OptimizationRequestOutput, 0), nil
}
