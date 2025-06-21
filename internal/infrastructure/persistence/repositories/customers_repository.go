/*
Package repositories provides interfaces for data-access
*/
package repositories

import (
	"context"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
)

type CustomersRepository interface {
	Create(ctx context.Context, input *dtos.SaveCustomerInput) error
	SaveAPIKey(ctx context.Context, input *dtos.SaveAPIKeyInput) error
	GetOneByAPIKey(ctx context.Context, apiKey string) *dtos.CustomerOutput
	GetOneByBusinessIdentifier(ctx context.Context, businessIdentifier string) (*dtos.CustomerOutput, error)
}
