/*
Package validators provides custom HTTP validators for the Gin Server Engine
*/
package validators

import (
	"github.com/go-playground/validator/v10"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/enums"
)

func IsCargoKindValid(fl validator.FieldLevel) bool {
	cargoKind := fl.Field().String()
	switch enums.CargoKind(cargoKind) {
	case enums.CargoKindBulkCargo,
		enums.CargoKindContainerizedCargo,
		enums.CargoKindRefrigeratedCargo,
		enums.CargoKindDryCargo,
		enums.CargoKindAliveCargo,
		enums.CargoKindDangerousCargo,
		enums.CargoKindFragileCargo,
		enums.CargoKindIndivisibleAndExceptionalCargo,
		enums.CargoKindVehicleCargo:
		return true
	}

	return false
}
