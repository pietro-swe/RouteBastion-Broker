package customer

import (
	"context"

	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db"
	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db/generated"
	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/dbutils"
)

type CustomersStore interface {
	Create(ctx context.Context, input *shared.SaveCustomerInput) error
	SaveAPIKey(ctx context.Context, input *shared.SaveAPIKeyInput) error
	GetOneByAPIKey(ctx context.Context, apiKey string) (*Customer, error)
	GetOneByBusinessIdentifier(ctx context.Context, businessIdentifier string) (*Customer, error)
}

type CustomersStoreImpl struct {
	queries *generated.Queries
}

func NewCustomersStore(
	provider db.DBProvider,
) CustomersStore {
	return &CustomersStoreImpl{
		queries: generated.New(provider.GetConn()),
	}
}

func (s *CustomersStoreImpl) Create(ctx context.Context, input *shared.SaveCustomerInput) error {
	tx, err := dbutils.ExtractTx(ctx)
	if err != nil {
		return err
	}

	q := s.queries.WithTx(tx)

	_, err = q.CreateCustomer(ctx, generated.CreateCustomerParams{
		ID:                 input.ID,
		Name:               input.Name,
		BusinessIdentifier: input.BusinessIdentifier,
	})

	return err
}

func (s *CustomersStoreImpl) GetOneByAPIKey(ctx context.Context, apiKey string) (*Customer, error) {
	row, err := s.queries.GetCustomerByApiKey(ctx, apiKey)
	if err != nil {
		return nil, err
	}

	return NewCustomer(
		row.ModelCustomer.ID,
		row.ModelCustomer.BusinessIdentifier,
		row.ModelCustomer.Name,
		row.ModelApiKey.Key,
		row.ModelCustomer.CreatedAt.Time,
		nil,
		nil,
	), nil
}

func (s *CustomersStoreImpl) GetOneByBusinessIdentifier(
	ctx context.Context,
	businessIdentifier string,
) (*Customer, error) {
	row, err := s.queries.GetOneCustomerByBusinessIdentifier(
		ctx,
		businessIdentifier,
	)
	if err != nil {
		return nil, err
	}

	return NewCustomer(
		row.ModelCustomer.ID,
		row.ModelCustomer.BusinessIdentifier,
		row.ModelCustomer.Name,
		row.ModelApiKey.Key,
		row.ModelCustomer.CreatedAt.Time,
		nil,
		nil,
	), nil
}

func (s *CustomersStoreImpl) SaveAPIKey(ctx context.Context, input *shared.SaveAPIKeyInput) error {
	tx, err := dbutils.ExtractTx(ctx)
	if err != nil {
		return err
	}

	q := s.queries.WithTx(tx)

	if input.CreatedAt != nil {
		createdAt := dbutils.ConvertTimeToPgtypeTimestamp(*input.CreatedAt)

		_, err = q.CreateApiKey(ctx, generated.CreateApiKeyParams{
			ID:         input.ID,
			Key:        input.APIKey,
			CustomerID: input.CustomerID,
			CreatedAt:  createdAt,
		})

		return err
	}

	if input.ModifiedAt != nil {
		modifiedAt := dbutils.ConvertTimeToPgtypeTimestamp(*input.ModifiedAt)

		err = q.UpdateApiKey(ctx, generated.UpdateApiKeyParams{
			ID:         input.ID,
			Key:        input.APIKey,
			ModifiedAt: modifiedAt,
		})

		return err
	}

	modifiedAt := dbutils.ConvertTimeToPgtypeTimestamp(*input.ModifiedAt)
	deletedAt := dbutils.ConvertTimeToPgtypeTimestamp(*input.DeletedAt)

	err = q.DeleteApiKey(ctx, generated.DeleteApiKeyParams{
		ID:         input.ID,
		ModifiedAt: modifiedAt,
		DeletedAt:  deletedAt,
	})

	return err
}
