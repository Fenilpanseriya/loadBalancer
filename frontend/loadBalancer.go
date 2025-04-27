package frontend

import "net/http"

type LoadBalancer interface {
	Serve(http.ResponseWriter, *http.Request)
}

const (
	RETRY_COUNT = 0
)

type loadBalancer struct {
	serverPool []string
}

func (lb *loadBalancer) Serve(w http.ResponseWriter, r *http.Request) {

}

func NewLoadBalancer(serverPool []string) LoadBalancer {
	return &loadBalancer{
		serverPool: serverPool,
	}
}
