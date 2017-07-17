package paxos

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
)

// Agent ...
type Agent struct {
	Node
	Peers []string
}

// NewAgent returns a new agent instance
func NewAgent(id string, host string, port int, peers []string) *Agent {
	return &Agent{
		Node:  NewNode(id, host, port),
		Peers: peers,
	}
}

// Run ...
func (a *Agent) Run() error {
	// Start RPC server
	listener, err := net.Listen("tcp", a.GetAddress())
	if err != nil {
		return fmt.Errorf("TCP Listen: %v", err)
	}

	// Agent is acceptor
	acceptor := NewAcceptor(a.ID, a.Host, a.Port)

	// Register acceptor to RPC
	if err = rpc.Register(acceptor); err != nil {
		return fmt.Errorf("Register acceptor to RPC: %v", err)
	}

	log.Printf("Starting agent ID: %s\n", a.ID)
	log.Printf("Listening on: %s\n", a.GetAddress())
	log.Printf("Peers: %s\n", a.Peers)

	rpc.HandleHTTP()

	// Agent will propose its own address to all peers,
	// we want to reach a consensus on who the cluster leader is
	proposer := NewProposer(a.ID, a.Host, a.Port, a.Peers)
	go func() {
		if err := proposer.Propose(NewProposal("leader", a.GetAddress())); err != nil {
			log.Print(err)
		}
	}()

	return http.Serve(listener, nil)
}
