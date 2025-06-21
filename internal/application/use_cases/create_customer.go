/*
Package usecases provides underlying functionality orchestration for the application
*/
package usecases

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/cryptography"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/repositories"
	"github.com/marechal-dev/RouteBastion-Broker/internal/shared"
	uuid "github.com/satori/go.uuid"
)

type CreateCustomerUseCase interface {
	Execute(ctx context.Context, dto *dtos.CreateCustomerInput) (*dtos.CustomerOutput, error)
}

type CreateCustomerUseCaseImpl struct {
	tx   persistence.TxManager
	repo repositories.CustomersRepository
	key []byte
	hashGen  cryptography.HashGenerator
}

func NewCreateCustomerUseCase(
	tx persistence.TxManager,
	repo repositories.CustomersRepository,
	key []byte,
	hashGen cryptography.HashGenerator,
) *CreateCustomerUseCaseImpl {
	return &CreateCustomerUseCaseImpl{
		tx:   tx,
		repo: repo,
		key: key,
		hashGen:  hashGen,
	}
}

func (uc *CreateCustomerUseCaseImpl) Execute(
	ctx context.Context,
	dto *dtos.CreateCustomerInput,
) (*dtos.CustomerOutput, error) {
	customer, err := uc.repo.GetOneByBusinessIdentifier(ctx, dto.BusinessIdentifier)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, shared.InfrastructureError{
			Code: shared.ErrCodeDatabaseFailure,
			Msg:  err.Error(),
		}
	}

	if customer != nil {
		return nil, shared.ApplicationError{
			Code: shared.ErrCodeConflict,
			Msg:  "customer already exists",
		}
	}

	return persistence.WithinTransactionReturning(
		uc.tx,
		ctx,
		func(txCtx context.Context) (*dtos.CustomerOutput, error) {
			rawKey := uuid.NewV4().String()
			hashed, err := uc.hashGen.Generate(uc.key, rawKey)
			if err != nil {
				return nil, shared.InfrastructureError{
					Code: shared.ErrCodeEncryptionFailure,
					Msg:  err.Error(),
				}
			}

			customerInput := &dtos.SaveCustomerInput{
				ID: uuid.NewV4(),
				Name: dto.Name,
				BusinessIdentifier: dto.BusinessIdentifier,
				APIKey: hashed,
			}

			err = uc.repo.Create(txCtx, customerInput)
			if err != nil {
				return nil, shared.InfrastructureError{
					Code: shared.ErrCodeDatabaseFailure,
					Msg:  err.Error(),
				}
			}

			err = uc.repo.SaveAPIKey(txCtx, &dtos.SaveAPIKeyInput{
				ID: uuid.NewV4(),
				APIKey: hashed,
				CustomerID: customerInput.ID,
				CreatedAt: &time.Time{},
				ModifiedAt: nil,
				DeletedAt: nil,
			})
			if err != nil {
				return nil, shared.InfrastructureError{
					Code: shared.ErrCodeDatabaseFailure,
					Msg:  err.Error(),
				}
			}

			return customer, nil
		},
	)
}
