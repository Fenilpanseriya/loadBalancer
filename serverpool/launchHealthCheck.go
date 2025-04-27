package serverpool

import (
	"context"
	"fmt"
	"time"
)

func LaunchHealthCheck(ctx context.Context, serverPool ServerPool) {
	t := time.NewTicker(time.Second * 20)
	defer ctx.Done()
	select {
	case <-ctx.Done():
		fmt.Println("health check timed out...............")
		return
	case <-t.C:
		fmt.Println("health check started")
		go HealthCheckServer(ctx, serverPool)
		fmt.Println("health check completed")
	}
}
