package customer

import (
	"context"
	"time"

	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db"
	"github.com/pietro-swe/RouteBastion-Broker/internal/infra/db/generated"
	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/dbutils"
	uuid "github.com/satori/go.uuid"
)

type CustomersStore interface {
	Create(ctx context.Context, input *shared.SaveCustomerInput) error
	Delete(ctx context.Context, customerID uuid.UUID) error
	GetByBusinessIdentifier(ctx context.Context, businessIdentifier string) (*Customer, error)
	GetByID(ctx context.Context, customerID uuid.UUID) (*Customer, error)
	SaveAPIKey(ctx context.Context, input *shared.SaveAPIKeyInput) error
	GetByAPIKey(ctx context.Context, apiKey string) (*Customer, error)
	DeleteAllAPIKeysByCustomerID(ctx context.Context, customerID uuid.UUID) error
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

func (s *CustomersStoreImpl) Delete(ctx context.Context, customerID uuid.UUID) error {
	tx, err := dbutils.ExtractTx(ctx)
	if err != nil {
		return err
	}

	q := s.queries.WithTx(tx)

	return q.DisableCustomer(ctx, generated.DisableCustomerParams{
		ID:        customerID,
		DeletedAt: dbutils.ConvertTimeToPgtypeTimestamp(time.Now()),
	})
}

func (s *CustomersStoreImpl) GetByBusinessIdentifier(
	ctx context.Context,
	businessIdentifier string,
) (*Customer, error) {
	row, err := s.queries.GetCustomerByBusinessIdentifier(
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

func (s *CustomersStoreImpl) GetByID(ctx context.Context, customerID uuid.UUID) (*Customer, error) {
	row, err := s.queries.GetCustomerByID(ctx, customerID)
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

func (s *CustomersStoreImpl) DeleteAllAPIKeysByCustomerID(ctx context.Context, customerID uuid.UUID) error {
	tx, err := dbutils.ExtractTx(ctx)
	if err != nil {
		return err
	}

	q := s.queries.WithTx(tx)

	modifiedAt := dbutils.ConvertTimeToPgtypeTimestamp(time.Now())
	deletedAt := dbutils.ConvertTimeToPgtypeTimestamp(time.Now())

	return q.DeleteAllApiKeysByCustomerID(ctx, generated.DeleteAllApiKeysByCustomerIDParams{
		CustomerID: customerID,
		ModifiedAt: modifiedAt,
		DeletedAt:  deletedAt,
	})
}

func (s *CustomersStoreImpl) GetByAPIKey(ctx context.Context, apiKey string) (*Customer, error) {
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
