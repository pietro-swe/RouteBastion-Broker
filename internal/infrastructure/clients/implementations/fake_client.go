package implementations

import (
	"context"
	"errors"
	"math"

	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/enums"
)

type stop struct {
    ShipmentID string
    Type       string
    Point      dtos.Point
}

type FakeRouteOptimizer struct{}

func NewFakeRouteOptimizer() *FakeRouteOptimizer {
    return &FakeRouteOptimizer{}
}

func (f *FakeRouteOptimizer) OptimizeSync(_ context.Context, input dtos.OptimizationRequestInput) ([]dtos.OptimizationRequestOutput, error) {
    if len(input.Vehicles) == 0 {
        return make([]dtos.OptimizationRequestOutput, 0), errors.New("no vehicles provided")
    }

    vehicle := input.Vehicles[0]
    start := vehicle.Start

    // Flatten pickups and deliveries into stops
    stops := make([]stop, len(input.Shipments)*2)
    for _, s := range input.Shipments {
        stops = append(stops, stop{
            ShipmentID: s.ID,
            Type:       "pickup",
            Point:      s.Pickup,
        })
        stops = append(stops, stop{
            ShipmentID: s.ID,
            Type:       "delivery",
            Point:      s.Delivery,
        })
    }

    visited := make(map[int]bool)
    pickupsDone := make(map[string]bool)
    steps := make([]dtos.RouteStep, len(stops))
    current := start
    totalDistance := 0.0

    for len(visited) < len(stops) {
        bestIdx := -1
        bestDist := math.MaxFloat64

        for i, s := range stops {
            if visited[i] {
                continue
            }
            // Can only visit a delivery if pickup is done
            if s.Type == "delivery" && !pickupsDone[s.ShipmentID] {
                continue
            }

            dist := haversineDistance(current.Latitude, current.Longitude, s.Point.Latitude, s.Point.Longitude)
            if dist < bestDist {
                bestDist = dist
                bestIdx = i
            }
        }

        if bestIdx == -1 {
            return make([]dtos.OptimizationRequestOutput, 0), errors.New("no valid next stop found — constraint deadlock")
        }

        stop := stops[bestIdx]
        visited[bestIdx] = true
        if stop.Type == "pickup" {
            pickupsDone[stop.ShipmentID] = true
        }

        totalDistance += haversineDistance(current.Latitude, current.Longitude, stop.Point.Latitude, stop.Point.Longitude)

		stepKind := enums.RouteStepKind(stop.Type)

		steps = append(steps, dtos.RouteStep{
            ShipmentID: stop.ShipmentID,
            Kind:       stepKind,
            Location:   stop.Point,
        })

        current = stop.Point
    }

    // Optionally return to end
    if vehicle.End != nil {
        totalDistance += haversineDistance(current.Latitude, current.Longitude, vehicle.End.Latitude, vehicle.End.Longitude)
        steps = append(steps, dtos.RouteStep{
            ShipmentID: "END",
            Kind:       "end",
            Location:   *vehicle.End,
        })
    }

	result := make([]dtos.OptimizationRequestOutput, 1)

	result = append(result, dtos.OptimizationRequestOutput{
        VehicleID:     vehicle.ID,
        Steps:         steps,
        TotalDistanceInKilometers: totalDistance,
    })

    return result, nil
}

func haversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
    const R = 6371 // Earth radius in km

    dLat := degreesToRadians(lat2 - lat1)
    dLon := degreesToRadians(lon2 - lon1)

    a := math.Sin(dLat / 2) * math.Sin(dLat / 2) +
		math.Cos(degreesToRadians(lat1)) * math.Cos(degreesToRadians(lat2)) *
		math.Sin(dLon / 2) * math.Sin(dLon / 2)

    c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
    return R * c
}

func degreesToRadians(deg float64) float64 {
    return deg * math.Pi / 180
}
