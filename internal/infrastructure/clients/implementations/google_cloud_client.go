package clients

import (
	"context"
	"fmt"

	routeoptimization "cloud.google.com/go/maps/routeoptimization/apiv1"
	"github.com/marechal-dev/RouteBastion-Broker/internal/application/dtos"
	"google.golang.org/genproto/googleapis/type/latlng"

	rpb "cloud.google.com/go/maps/routeoptimization/apiv1/routeoptimizationpb"
)

type GoogleCloudClient struct {}

func (gcc *GoogleCloudClient) Optimize(input *dtos.OptimizationRequestInput) error {
	ctx := context.Background()
	c, err := routeoptimization.NewClient(ctx)
	if err != nil {
		return fmt.Errorf("failed to initialize routeoptimization client: %w", err)
	}
	defer c.Close()

	// pickups := []*rpb.Shipment_VisitRequest{}
	// deliveries := []*rpb.Shipment_VisitRequest{}
	// vehicles := []*rpb.Vehicle{}

	req := &rpb.OptimizeToursRequest{
		Parent: "projects/" + "routebastion",
		Model: &rpb.ShipmentModel{
			Shipments: []*rpb.Shipment{
				{
					Pickups: []*rpb.Shipment_VisitRequest{},
					Deliveries: []*rpb.Shipment_VisitRequest{
						{ArrivalLocation: &latlng.LatLng{Latitude: 48.880942, Longitude: 2.323866}},
					},
				},
			},
			Vehicles: []*rpb.Vehicle{
				{
					StartLocation: &latlng.LatLng{Latitude: 48.863102, Longitude: 2.341204},
					EndLocation:   &latlng.LatLng{Latitude: 48.86311, Longitude: 2.341205},
				},
			},
		},
	}

	_, err = c.OptimizeTours(ctx, req)
	if err != nil {
		fmt.Printf("error when calling OptimizeTours: %v", err)
		return err
	}

	return nil
}
