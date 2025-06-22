/*
Package implementations provides concrete implementations of data-access repositories
*/
package implementations

import (
	"context"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/mappers"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/generated"
	"github.com/marechal-dev/RouteBastion-Broker/internal/shared"
)

type PgCustomersRepository struct {
	queries *generated.Queries
}

func NewPgCustomersRepository(db persistence.DBProvider) *PgCustomersRepository {
	return &PgCustomersRepository{
		queries: generated.New(db.GetConn()),
	}
}

func (r *PgCustomersRepository) Create(ctx context.Context, input *dtos.SaveCustomerInput) error {
	tx, err := persistence.ExtractTx(ctx)
	if err != nil {
		return err
	}

	q := r.queries.WithTx(tx)

	_, err = q.CreateCustomer(ctx, generated.CreateCustomerParams{
		ID:                 input.ID,
		Name:               input.Name,
		BusinessIdentifier: input.BusinessIdentifier,
	})

	return err
}

func (r *PgCustomersRepository) GetOneByAPIKey(ctx context.Context, apiKey string) *dtos.CustomerOutput {
	row, err := r.queries.GetCustomerByApiKey(ctx, apiKey)

	if err != nil {
		return nil
	}

	// TODO: Fetch Vehicles here
	return mappers.CustomerToDomain(row.ModelCustomer, row.ModelApiKey, make([]generated.ModelVehicle, 0))
}

func (r *PgCustomersRepository) GetOneByBusinessIdentifier(
	ctx context.Context,
	businessIdentifier string,
) (*dtos.CustomerOutput, error) {
	row, err := r.queries.GetOneCustomerByBusinessIdentifier(
		ctx,
		businessIdentifier,
	)

	if err != nil {
		return nil, err
	}

	return mappers.CustomerToDomain(
		row.ModelCustomer,
		row.ModelApiKey,
		make([]generated.ModelVehicle, 0),
	), nil
}

func (r *PgCustomersRepository) SaveAPIKey(ctx context.Context, input *dtos.SaveAPIKeyInput) error {
	tx, err := persistence.ExtractTx(ctx)
	if err != nil {
		return err
	}

	q := r.queries.WithTx(tx)

	if input.CreatedAt != nil {
		createdAt := shared.ConvertTimeToPgtypeTimestamp(*input.CreatedAt)

		_, err = q.CreateApiKey(ctx, generated.CreateApiKeyParams{
			ID:         input.ID,
			Key:        input.APIKey,
			CustomerID: input.CustomerID,
			CreatedAt:  createdAt,
		})

		return err
	}

	if input.ModifiedAt != nil {
		modifiedAt := shared.ConvertTimeToPgtypeTimestamp(*input.ModifiedAt)

		err = q.UpdateApiKey(ctx, generated.UpdateApiKeyParams{
			ID: input.ID,
			Key: input.APIKey,
			ModifiedAt: modifiedAt,
		})

		return err
	}

	modifiedAt := shared.ConvertTimeToPgtypeTimestamp(*input.ModifiedAt)
	deletedAt := shared.ConvertTimeToPgtypeTimestamp(*input.DeletedAt)

	err = q.DeleteApiKey(ctx, generated.DeleteApiKeyParams{
		ID: input.ID,
		ModifiedAt: modifiedAt,
		DeletedAt: deletedAt,
	})

	return err
}
