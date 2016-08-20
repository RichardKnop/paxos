package agent

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/RichardKnop/paxos/acceptor"
	"github.com/RichardKnop/paxos/proposer"
	"github.com/RichardKnop/uuid"
)

// Agent ...
type Agent struct {
	id    string
	host  string
	port  int
	peers []string
}

// New returns new Agent instance
func New(id string, host string, port int, peers []string) (*Agent, error) {
	if id == "" {
		id = uuid.New()
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return &Agent{
		id:    id,
		host:  host,
		port:  port,
		peers: peers,
	}, nil
}

// Run ...
func (a *Agent) Run() error {
	// Start RPC server
	listener, err := net.Listen("tcp", a.getAddress())
	if err != nil {
		return fmt.Errorf("TCP Listen: %v", err)
	}

	// Agent is acceptor
	acceptor, err := acceptor.New(a.id, a.host, a.port)
	if err != nil {
		return fmt.Errorf("New Acceptor: %v", err)
	}
	acceptorRPC, err := acceptor.NewRPC()
	if err != nil {
		return fmt.Errorf("New Acceptor RPC: %v", err)
	}
	if err := rpc.Register(acceptorRPC); err != nil {
		return fmt.Errorf("Register Acceptor RPC: %v", err)
	}

	// Agent is proposer
	proposedValue := a.getAddress()
	proposer, err := proposer.New(a.id, a.host, a.port, proposedValue, a.peers)
	if err != nil {
		return fmt.Errorf("New Proposer: %v", err)
	}
	go proposer.Run()

	log.Printf("Starting agent ID: %s\n", a.id)
	log.Printf("Listening on: %s\n", a.getAddress())
	log.Printf("Peers: %s\n", a.peers)

	rpc.HandleHTTP()
	return http.Serve(listener, nil)
}

func (a *Agent) getAddress() string {
	return fmt.Sprintf("%s:%d", a.host, a.port)
}
