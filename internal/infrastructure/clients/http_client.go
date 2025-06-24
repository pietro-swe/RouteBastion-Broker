/*
Package clients provides an interface and concrete implementations for HTTP Clients
*/
package clients

import "github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"

type RouteOptimizationAPIClient interface {
	Optimize(input *dtos.OptimizationRequestInput)
}
