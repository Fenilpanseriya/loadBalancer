package serverpool

import (
	"sync"

	"github.com/fenilpanseriya/loadbalancer/backend"
)

type leastConnectionServerPool struct {
	servers []backend.Backend
	mux     sync.RWMutex
}

func (l *leastConnectionServerPool) GetServers() []backend.Backend {
	return l.servers
}

func (l *leastConnectionServerPool) GetServerPoolSize() int {
	return len(l.servers)
}

func (l *leastConnectionServerPool) GetAliveServer() backend.Backend {
	var minConnServer backend.Backend
	for _, server := range l.servers {
		if server.GetAlive() {
			if minConnServer == nil || server.GetActiveConnections() < minConnServer.GetActiveConnections() {
				minConnServer = server
			}
		}
	}
	return minConnServer
}
