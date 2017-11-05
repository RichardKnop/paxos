package rpc

import (
	"fmt"
	"net/rpc"

	"github.com/RichardKnop/paxos/paxos"
)

// Service methods
const (
	// AcceptorReceivePrepare ...
	AcceptorReceivePrepare = "Server.AcceptorReceivePrepare"
	// AcceptorReceivePropose ...
	AcceptorReceivePropose = "Server.AcceptorReceivePropose"
)

// RunServer registers service methods and HTTP handler for RPC messages
func RunServer(server *Server) error {
	if server.acceptor != nil {
		// Register acceptor to RPC
		if err := rpc.Register(server); err != nil {
			return fmt.Errorf("Register acceptor to RPC: %v", err)
		}
	}

	rpc.HandleHTTP()

	return nil
}

// Server ...
type Server struct {
	acceptor *paxos.Acceptor
}

// NewServer returns new instance of Server
func NewServer(acceptor *paxos.Acceptor) *Server {
	return &Server{
		acceptor: acceptor,
	}
}

// AcceptorReceivePrepare handles prepare requests from proposers
func (s *Server) AcceptorReceivePrepare(proposal, reply *paxos.Proposal) error {
	proposal, err := s.acceptor.ReceivePrepare(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}

// AcceptorReceivePropose handles propose requests from proposers
func (s *Server) AcceptorReceivePropose(proposal, reply *paxos.Proposal) error {
	proposal, err := s.acceptor.ReceivePropose(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}
