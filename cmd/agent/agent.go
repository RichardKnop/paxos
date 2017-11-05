package agent

import (
	"fmt"
	"log"
	"net"
	"net/http"

	"github.com/RichardKnop/paxos/paxos"
	"github.com/RichardKnop/paxos/rpc"
	"github.com/satori/go.uuid"
)

// Agent ...
type Agent struct {
	ID    string
	Host  string
	Port  int
	Peers []string
}

// New returns a new agent instance
func New(id string, host string, port int, peers []string) *Agent {
	if id == "" {
		id = uuid.NewV4().String()
	}
	if host == "" {
		host = "127.0.0.1"
	}
	return &Agent{
		ID:    id,
		Host:  host,
		Port:  port,
		Peers: peers,
	}
}

// GetAddress concatenates host and port
func (a *Agent) GetAddress() string {
	return fmt.Sprintf("%s:%d", a.Host, a.Port)
}

// String returns a human readable representation
func (a *Agent) String() string {
	return a.GetAddress()
}

// Run ...
func (a *Agent) Run() error {
	// Start RPC server
	listener, err := net.Listen("tcp", a.GetAddress())
	if err != nil {
		return fmt.Errorf("TCP Listen: %v", err)
	}

	// Create a new acceptor
	acceptor := paxos.NewAcceptor()

	// Create a proposer with acceptor clients
	acceptorClients := make([]paxos.AcceptorClientInterface, len(a.Peers)+1)
	for i := 0; i < len(a.Peers); i++ {
		acceptorClients[i] = rpc.NewClient(a.Peers[i])
	}
	acceptorClients[len(a.Peers)] = rpc.NewClient(a.GetAddress())
	proposer := paxos.NewProposer(acceptorClients)

	// RPC server
	server := rpc.NewServer(acceptor)
	if err = rpc.RunServer(server); err != nil {
		return err
	}

	log.Printf("Starting agent ID: %s\n", a.ID)
	log.Printf("Listening on: %s\n", a.GetAddress())
	log.Printf("Peers: %s\n", a.Peers)

	// Agent will propose its own address to all peers,
	// we want to reach a consensus on who the cluster leader is
	go func() {
		proposal := &paxos.Proposal{
			Key:   "leader",
			Value: []byte(a.GetAddress()),
		}
		if err := proposer.Propose(proposal); err != nil {
			log.Print(err)
		}
	}()

	return http.Serve(listener, nil)
}
