/*
Package presenters provides a presentation layer for internal data structures
*/
package presenters

import (
	"time"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
)

type CustomerPresenter struct {
	Name               string     `json:"name"`
	BusinessIdentifier string     `json:"businessIdentifier"`
	APIKey             string     `json:"apiKey"`
	CreatedAt          *time.Time `json:"createdAt"`
}

func CustomerFromDomain(customer *dtos.CustomerOutput) *CustomerPresenter {
	return &CustomerPresenter{
		Name:               customer.Name,
		BusinessIdentifier: customer.BusinessIdentifier,
		APIKey:             customer.APIKey,
		CreatedAt:          customer.CreatedAt,
	}
}
