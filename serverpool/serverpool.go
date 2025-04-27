package serverpool

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/fenilpanseriya/loadbalancer/backend"
)

type ServerPool interface {
	GetServers() []backend.Backend
	GetAliveServer() backend.Backend
	AddServer(backend.Backend)
	GetServerPoolSize() int
}

type roundRobbinServerPool struct {
	servers []backend.Backend
	current int
	mux     sync.RWMutex
}

func (r *roundRobbinServerPool) GetServerPoolSize() int {
	return len(r.servers)
}

func (r *roundRobbinServerPool) GetServers() []backend.Backend {
	return r.servers
}

func (r *roundRobbinServerPool) Rotate() backend.Backend {

	r.mux.Lock()
	defer r.mux.Unlock()
	r.current = (r.current + 1) % len(r.servers)
	r.current++
	return r.servers[r.current]
}

func (r *roundRobbinServerPool) GetAliveServer() backend.Backend {
	for i := 0; i < r.GetServerPoolSize(); i++ {
		nextPeer := r.Rotate()
		if nextPeer.GetAlive() {
			return nextPeer
		}
	}
	return nil
}
func (r *roundRobbinServerPool) AddServer(server backend.Backend) {
	r.servers = append(r.servers, server)
}

func NewServerPool(servers []backend.Backend) (ServerPool, error) {
	return &roundRobbinServerPool{
		servers: servers,
		current: 0,
	}, nil
}

func HealthCheckServer(ctx context.Context, serverPool ServerPool) {
	aliveChannel := make(chan bool, 1)
	for _, server := range serverPool.GetServers() {
		server := server
		requestCtx, stop := context.WithTimeout(ctx, 10*time.Second)
		defer stop()

		status := "up"
		go backend.IsServerAlive(requestCtx, aliveChannel, server.GetURL())

		select {
		case <-ctx.Done():
			fmt.Println("health check timed out")
		case alive := <-aliveChannel:
			if alive {
				fmt.Printf("Server %s is alive\n", server.GetURL())
				fmt.Println("server status is ", status)
				server.SetAlive(alive)
			} else {
				status = "down"
				fmt.Printf("Server %s is down\n", server.GetURL())
				server.SetAlive(alive)
				fmt.Println("server status is ", status)
			}
		}
	}
}
