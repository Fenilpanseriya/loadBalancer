package backend

import (
	"context"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend interface {
	SetAlive(alive bool) bool
	GetAlive() bool
	GetURL() *url.URL
	Serve(http.ResponseWriter, *http.Request)
	GetActiveConnections() int
}

type backend struct {
	url          *url.URL
	alive        bool
	mux          sync.RWMutex
	connections  int
	reverseProxy *httputil.ReverseProxy
}

func (b *backend) SetAlive(alive bool) bool {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.alive = alive
	return b.alive
}

func (b *backend) GetActiveConnections() int {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.connections
}

func (b *backend) GetAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.alive
}

func (b *backend) GetURL() *url.URL {
	return b.url
}

func (b *backend) Serve(w http.ResponseWriter, r *http.Request) {
	defer func() {
		b.mux.Lock()
		b.connections--
		b.mux.Unlock()
	}()
	b.mux.Lock()
	b.connections++
	b.mux.Unlock()
	b.reverseProxy.ServeHTTP(w, r)
}

func NewBackend(url *url.URL, rp *httputil.ReverseProxy) Backend {
	return &backend{
		url:          url,
		alive:        true,
		connections:  0,
		reverseProxy: rp,
	}
}
func IsServerAlive(ctx context.Context, aliveChan chan bool, u *url.URL) {
	var d net.Dialer
	conn, err := d.DialContext(ctx, "tcp", u.Host)
	if err != nil {
		aliveChan <- false
		return
	}
	defer conn.Close()
	aliveChan <- true
}
