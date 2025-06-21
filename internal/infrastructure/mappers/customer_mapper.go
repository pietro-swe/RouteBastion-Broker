package mappers

import (
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	"github.com/marechal-dev/RouteBastion-Broker/internal/infrastructure/persistence/generated"
)

func CustomerToDomain(
	rawCustomer generated.ModelCustomer,
	rawAPIKey generated.ModelApiKey,
	rawVehicles []generated.ModelVehicle,
) *dtos.CustomerOutput {
	return &dtos.CustomerOutput{
		ID: rawCustomer.ID,
		Name: rawCustomer.Name,
		BusinessIdentifier: rawCustomer.BusinessIdentifier,
		APIKey: rawAPIKey.Key,
		CreatedAt: &rawCustomer.CreatedAt.Time,
		ModifiedAt: &rawCustomer.ModifiedAt.Time,
		DeletedAt: &rawCustomer.DeletedAt.Time,
	}
}
