/*
Package clients provides an interface and concrete implementations for HTTP Clients
*/
package clients

import (
	"context"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
)

type RouteOptimizationAPIClient interface {
	OptimizeSync(ctx context.Context, input dtos.OptimizationRequestInput) ([]dtos.OptimizationRequestOutput, error)
}
