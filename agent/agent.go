package agent

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"

	"github.com/RichardKnop/paxos/acceptor"
	"github.com/RichardKnop/paxos/models"
	"github.com/RichardKnop/paxos/proposer"
	"github.com/RichardKnop/uuid"
)

// Agent ...
type Agent struct {
	ID    string
	Host  string
	Port  int
	Peers []string
}

// New returns new Agent instance
func New(theID string, host string, port int, peers []string) (*Agent, error) {
	if theID == "" {
		theID = uuid.New()
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return &Agent{
		ID:    theID,
		Host:  host,
		Port:  port,
		Peers: peers,
	}, nil
}

// Run ...
func (a *Agent) Run() error {
	// Start RPC server
	listener, err := net.Listen("tcp", a.GetAddress())
	if err != nil {
		return fmt.Errorf("TCP Listen: %v", err)
	}

	// Agent is acceptor
	acceptor, err := acceptor.New(a.ID, a.Host, a.Port)
	if err != nil {
		return fmt.Errorf("New Acceptor: %v", err)
	}
	acceptorRPC, err := acceptor.NewRPC()
	if err != nil {
		return fmt.Errorf("New Acceptor RPC: %v", err)
	}
	if err = rpc.Register(acceptorRPC); err != nil {
		return fmt.Errorf("Register Acceptor RPC: %v", err)
	}

	log.Printf("Starting agent ID: %s\n", a.ID)
	log.Printf("Listening on: %s\n", a.GetAddress())
	log.Printf("Peers: %s\n", a.Peers)

	rpc.HandleHTTP()

	// Agent will propose its own address to all peers,
	// we want to reach a consensus on who the cluster leader is
	proposer, err := proposer.New(a.ID, a.Host, a.Port, a.Peers)
	if err != nil {
		return fmt.Errorf("New Proposer: %v", err)
	}
	go proposer.Propose(models.NewProposal("leader", a.GetAddress()))

	return http.Serve(listener, nil)
}

// GetAddress concatenates host and port
func (a *Agent) GetAddress() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}
