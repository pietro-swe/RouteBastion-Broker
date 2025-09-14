package httputils

import (
	"github.com/go-playground/validator/v10"
	"github.com/pietro-swe/RouteBastion-Broker/internal/application/enums"
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
