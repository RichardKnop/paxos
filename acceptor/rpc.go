package acceptor

import (
	"github.com/RichardKnop/paxos/models"
)

const (
	// PrepareServiceMethod ...
	PrepareServiceMethod = "RPC.Prepare"
	// ProposeServiceMethod ...
	ProposeServiceMethod = "RPC.Propose"
)

// RPC ...
type RPC struct {
	acceptor *Acceptor
}

// NewRPC returns new RPC instance
func (a *Acceptor) NewRPC() (*RPC, error) {
	return &RPC{
		acceptor: a,
	}, nil
}

// Prepare ...
func (r *RPC) Prepare(proposal *models.Proposal, reply *models.Proposal) error {
	proposal, err := r.acceptor.receivePrepare(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}

// Propose ...
func (r *RPC) Propose(proposal *models.Proposal, reply *models.Proposal) error {
	proposal, err := r.acceptor.receiveProposal(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}
