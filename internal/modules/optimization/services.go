package optimization

import (
	"context"
	"time"

	"github.com/pietro-swe/RouteBastion-Broker/internal/shared"
)

func OptimizeSync(ctx context.Context) ([]shared.OptimizationRequestOutput, error) {
	time.Sleep(time.Millisecond * 10)

	return make([]shared.OptimizationRequestOutput, 0), nil
}
