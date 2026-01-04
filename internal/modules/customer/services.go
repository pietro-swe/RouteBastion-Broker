package customer

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/crypto"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/customerrors"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/dbutils"
	uuid "github.com/satori/go.uuid"
)

func CreateCustomer(
	ctx context.Context,
	tx dbutils.TxManager,
	store CustomersStore,
	keyGen crypto.HashGenerator,
	input shared.CreateCustomerInput,
) (*Customer, error) {
	customer, err := store.GetOneByBusinessIdentifier(ctx, input.BusinessIdentifier)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.NewInfrastructureError(
			customerrors.ErrCodeDatabaseFailure,
			err.Error(),
			err,
		)
	}

	if customer != nil {
		return nil, customerrors.NewApplicationError(
			customerrors.ErrCodeConflict,
			"customer already exists",
			err,
		)
	}

	return dbutils.WithinTransactionReturning(
		ctx,
		tx,
		func(txCtx context.Context) (*Customer, error) {
			rawKey := uuid.NewV4().String()
			hashed, err := keyGen.Generate(rawKey)
			if err != nil {
				return nil, customerrors.NewInfrastructureError(
					customerrors.ErrCodeEncryptionFailure,
					err.Error(),
					err,
				)
			}

			now := time.Now()

			customer := NewCustomer(
				uuid.NewV4(),
				input.Name,
				input.BusinessIdentifier,
				hashed,
				now,
				nil,
				nil,
			)

			customerInput := &shared.SaveCustomerInput{
				ID:                 uuid.NewV4(),
				Name:               input.Name,
				BusinessIdentifier: input.BusinessIdentifier,
				APIKey:             hashed,
				CreatedAt:          &now,
				ModifiedAt:         nil,
				DeletedAt:          nil,
			}

			err = store.Create(txCtx, customerInput)
			if err != nil {
				return nil, customerrors.NewInfrastructureError(
					customerrors.ErrCodeDatabaseFailure,
					err.Error(),
					err,
				)
			}

			err = store.SaveAPIKey(txCtx, &shared.SaveAPIKeyInput{
				ID:         uuid.NewV4(),
				APIKey:     hashed,
				CustomerID: customerInput.ID,
				CreatedAt:  &now,
				ModifiedAt: nil,
				DeletedAt:  nil,
			})
			if err != nil {
				return nil, customerrors.NewInfrastructureError(
					customerrors.ErrCodeDatabaseFailure,
					err.Error(),
					err,
				)
			}

			return customer, nil
		},
	)
}

func GetOneCustomerByAPIKey(
	ctx context.Context,
	store CustomersStore,
	apiKey string,
) (*Customer, error) {
	customer, err := store.GetOneByAPIKey(ctx, apiKey)

	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.NewInfrastructureError(
			customerrors.ErrCodeDatabaseFailure,
			err.Error(),
			err,
		)
	}

	if customer == nil {
		return nil, nil
	}

	return customer, nil
}
