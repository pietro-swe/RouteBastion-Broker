/*
Package dtos provides input and output data structures
*/
package dtos

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CreateCustomerInput struct {
	Name               string `json:"name"               binding:"required"`
	BusinessIdentifier string `json:"businessIdentifier" binding:"required"`
}

type SaveCustomerInput struct {
	ID uuid.UUID
	Name string
	BusinessIdentifier string
	APIKey string
	CreatedAt *time.Time
	ModifiedAt *time.Time
	DeletedAt *time.Time
}

type CustomerOutput struct {
	ID uuid.UUID
	Name               string
	BusinessIdentifier string
	APIKey string
	CreatedAt  *time.Time
	ModifiedAt *time.Time
	DeletedAt  *time.Time
}

type SaveAPIKeyInput struct {
	ID uuid.UUID
	APIKey     string
	CustomerID uuid.UUID
	CreatedAt  *time.Time
	ModifiedAt *time.Time
	DeletedAt  *time.Time
}

// type GetCustomerByApiKeyDTO struct {
// 	ApiKey string `json:"apiKey" binding:"required,uuid4"`
// }

// type CreateVehicleDTO struct {
// 	Plate     string             `json:"plate"     binding:"required"`
// 	CargoType entities.CargoKind `json:"cargoType" binding:"required,cargoKind"`
// 	Capacity  float64            `json:"capacity"  binding:"required,min=1"`
// }

type GetAllVehiclesByAPIKey struct {
	APIKey string
}

type WaypointInput struct {
	Lat  float64 `json:"lat"  binding:"required,number"`
	Long float64 `json:"long" binding:"required,number"`
}

type VehicleInput struct {
	ID        string  `json:"id"        binding:"required,uuid4"`
	StartLat  float64 `json:"startLat"  binding:"required,number"`
	StartLong float64 `json:"startLong" binding:"required,number"`
	EndLat    float64 `json:"endLat"    binding:"required,number"`
	EndLong   float64 `json:"endLong"   binding:"required,number"`
}

type OptimizationRequestInput struct {
	Pickups    []WaypointInput `json:"pickups"    binding:"required"`
	Deliveries []WaypointInput `json:"deliveries" binding:"required"`
	Vehicles   []VehicleInput  `json:"vehicles"   binding:"required"`
}

