package usecases

import (
	"context"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories"
)

type GetOneCustomerUseCase interface {
	Execute(apiKey string) *dtos.CustomerOutput
}

type GetOneCustomerUseCaseImpl struct {
	repo repositories.CustomersRepository
}

func NewGetOneCustomerUseCaseImpl(repo repositories.CustomersRepository) *GetOneCustomerUseCaseImpl {
	return &GetOneCustomerUseCaseImpl{
		repo: repo,
	}
}

func (uc *GetOneCustomerUseCaseImpl) Execute(apiKey string) *dtos.CustomerOutput {
	foundCustomer := uc.repo.GetOneByAPIKey(context.Background(), apiKey)

	if foundCustomer == nil {
		return nil
	}

	return foundCustomer
}
