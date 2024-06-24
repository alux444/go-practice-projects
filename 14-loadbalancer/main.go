package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
)

const PORT = "8000"

type Server interface {
	Address() string
	IsAlive() bool
	Serve(writer http.ResponseWriter, request *http.Request)
}

type simpleServer struct {
	address string
	proxy   *httputil.ReverseProxy
}

func createSimpleServer(address string) *simpleServer {
	serverUrl, err := url.Parse(address)
	handleError(err)
	return &simpleServer{address: address, proxy: httputil.NewSingleHostReverseProxy(serverUrl)}
}

type LoadBalancer struct {
	port            string
	roundRobinCount int
	servers         []Server
}

func createLoadBalancer(port string, servers []Server) *LoadBalancer {
	return &LoadBalancer{
		port:            port,
		roundRobinCount: 0,
		servers:         servers,
	}
}

func (s *simpleServer) Address() string { return s.address }

// just have if the server exists, say it is alive.
func (s *simpleServer) IsAlive() bool { return true }

func (s *simpleServer) Serve(writer http.ResponseWriter, request *http.Request) {
	s.proxy.ServeHTTP(writer, request)
}

func (lb *LoadBalancer) getNextAvailableServer() Server {
	server := lb.servers[lb.roundRobinCount%len(lb.servers)]
	// while server isnt alive
	for !server.IsAlive() {
		lb.roundRobinCount++
		server = lb.servers[lb.roundRobinCount%len(lb.servers)]
	}
	lb.roundRobinCount++
	//return the available servers
	return server
}

func (lb *LoadBalancer) serveProxy(writer http.ResponseWriter, request *http.Request) {
	targetServer := lb.getNextAvailableServer()
	fmt.Printf("Forwarding address to: %q\n", targetServer.Address())
	targetServer.Serve(writer, request)
}

func handleError(err error) {
	if err != nil {
		fmt.Printf("error: %v", err)
		os.Exit(1)
	}
}

func main() {
	servers := []Server{
		createSimpleServer("https://github.com/"),
		createSimpleServer("https://www.duckduckgo.com/"),
		createSimpleServer("https://alux444.github.io/"),
	}
	lb := createLoadBalancer(PORT, servers)
	handleRedirect := func(writer http.ResponseWriter, request *http.Request) {
		lb.serveProxy(writer, request)
	}
	http.HandleFunc("/", handleRedirect)
	fmt.Printf("Server started at 'localhost:%v'\n", PORT)
	http.ListenAndServe(":"+lb.port, nil)
}
