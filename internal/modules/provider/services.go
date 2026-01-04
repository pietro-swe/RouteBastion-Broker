package provider

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
	"github.com/pietro-swe/RouteBastion-Broker/pkg/customerrors"
	uuid "github.com/satori/go.uuid"
)

func CreateProvider(
	ctx context.Context,
	store ProvidersStore,
	input shared.CreateProviderInput,
) (*Provider, error) {
	provider := NewProvider(input.Name)

	err := store.Create(ctx, provider)
	if err != nil {
		return nil, err
	}

	return provider, nil
}

func UpdateProvider(
	ctx context.Context,
	store ProvidersStore,
	input shared.UpdateProviderInput,
) (*Provider, error) {
	existingProvider, err := store.GetDetailsByID(ctx, input.ID)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.NewInfrastructureError(
			customerrors.ErrCodeDatabaseFailure,
			err.Error(),
			err,
		)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return nil, customerrors.NewApplicationError(
			customerrors.ErrCodeNotFound,
			"provider not found",
			err,
		)
	}

	now := time.Now()

	provider := &Provider{
		ID:         uuid.FromStringOrNil(input.ID),
		Name:       input.Name,
		CreatedAt:  existingProvider.CreatedAt,
		ModifiedAt: &now,
		DeletedAt:  existingProvider.DeletedAt,
	}

	updateErr := store.Update(ctx, provider)
	if updateErr != nil {
		return nil, customerrors.NewInfrastructureError(
			customerrors.ErrCodeDatabaseFailure,
			updateErr.Error(),
			updateErr,
		)
	}

	return provider, nil
}

func DeleteProvider(
	ctx context.Context,
	store ProvidersStore,
	id string,
) error {
	existingProvider, err := store.GetDetailsByID(ctx, id)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		return customerrors.NewInfrastructureError(
			customerrors.ErrCodeDatabaseFailure,
			err.Error(),
			err,
		)
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return customerrors.NewApplicationError(
			customerrors.ErrCodeNotFound,
			"provider not found",
			err,
		)
	}

	if existingProvider.DeletedAt != nil {
		return customerrors.NewApplicationError(
			customerrors.ErrCodeConflict,
			"provider already deleted",
			nil,
		)
	}

	deleteErr := store.Delete(ctx, id)
	if deleteErr != nil {
		return customerrors.NewInfrastructureError(
			customerrors.ErrCodeDatabaseFailure,
			deleteErr.Error(),
			deleteErr,
		)
	}

	return nil
}
