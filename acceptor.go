package paxos

import (
	"fmt"
	"log"
)

const (
	// PrepareServiceMethod ...
	PrepareServiceMethod = "Acceptor.Prepare"
	// ProposeServiceMethod ...
	ProposeServiceMethod = "Acceptor.Propose"
)

// Acceptor ...
type Acceptor struct {
	Node
	promisedProposals map[string]*Proposal
	acceptedProposals map[string]*Proposal
}

// NewAcceptor creates a new acceptor instance
func NewAcceptor(id, host string, port int) *Acceptor {
	return &Acceptor{
		Node:              NewNode(id, host, port),
		promisedProposals: make(map[string]*Proposal),
		acceptedProposals: make(map[string]*Proposal),
	}
}

// String returns a human readable representation
func (a *Acceptor) String() string {
	return fmt.Sprintf("Acceptor %s (%s:%d)", a.ID, a.Host, a.Port)
}

// Prepare handles received preparation request from proposers
func (a *Acceptor) Prepare(proposal *Proposal, reply *Proposal) error {
	proposal, err := a.receivePrepare(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}

// Propose handles received proposal request from proposers
func (a *Acceptor) Propose(proposal *Proposal, reply *Proposal) error {
	proposal, err := a.receiveProposal(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}

// If an acceptor receives a prepare request with number n greater
// than that of any prepare request to which it has already responded,
// then it responds to the request with a promise not to accept any more
// proposals numbered less than n and with the highest-numbered proposal
// (if any) that it has accepted.
func (a *Acceptor) receivePrepare(proposal *Proposal) (*Proposal, error) {
	// Do we already have a promise for this proposal
	promised, ok := a.promisedProposals[proposal.Key]

	// Ignore lesser or equally numbered proposals
	if ok && promised.Number >= proposal.Number {
		return nil, fmt.Errorf(
			"%s already promised to accept %s which is >= than requested %s",
			a,
			promised,
			proposal,
		)
	}

	// Promise to accept the proposal
	a.promisedProposals[proposal.Key] = proposal
	log.Printf("%s promises to accept proposal %s", a, proposal)

	return proposal, nil
}

// If an acceptor receives an accept request for a proposal numbered
// n, it accepts the proposal unless it has already responded to a prepare
// request having a number greater than n.
func (a *Acceptor) receiveProposal(proposal *Proposal) (*Proposal, error) {
	// Do we already have a promise for this proposal
	promised, ok := a.promisedProposals[proposal.Key]

	// Ignore lesser or equally numbered proposals
	if ok && promised.Number >= proposal.Number {
		return nil, fmt.Errorf(
			"%s already promised to accept %s which is >= than requested %s",
			a,
			promised,
			proposal,
		)
	}

	// Unexpected proposal
	if ok && promised.Number < proposal.Number {
		return nil, fmt.Errorf("%s received unexpected proposal %s", a, proposal)
	}

	// Accept the proposal
	a.acceptedProposals[proposal.Key] = proposal
	log.Printf("%s accepted proposal %s", a, proposal)

	return proposal, nil
}
