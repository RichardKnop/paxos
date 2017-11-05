package paxos

import (
	"fmt"
)

// Acceptor ...
type Acceptor struct {
	promisedProposals map[string]*Proposal
	acceptedProposals map[string]*Proposal
}

// AcceptorClientInterface ...
type AcceptorClientInterface interface {
	GetName() string
	SendPrepare(proposal *Proposal) (*Proposal, error)
	SendPropose(proposal *Proposal) (*Proposal, error)
}

// NewAcceptor creates a new acceptor instance
func NewAcceptor() *Acceptor {
	return &Acceptor{
		promisedProposals: make(map[string]*Proposal),
		acceptedProposals: make(map[string]*Proposal),
	}
}

// If an acceptor receives a prepare request with number n greater
// than that of any prepare request to which it has already responded,
// then it responds to the request with a promise not to accept any more
// proposals numbered less than n and with the highest-numbered proposal
// (if any) that it has accepted.
func (a *Acceptor) ReceivePrepare(proposal *Proposal) (*Proposal, error) {
	// Do we already have a promise for this proposal
	promised, ok := a.promisedProposals[proposal.Key]

	// Ignore lesser or equally numbered proposals
	if ok && promised.Number > proposal.Number {
		return nil, fmt.Errorf(
			"Already promised to accept %s which is > than requested %s",
			promised,
			proposal,
		)
	}

	// Promise to accept the proposal
	a.promisedProposals[proposal.Key] = proposal

	return proposal, nil
}

// If an acceptor receives a propose request for a proposal numbered
// n, it accepts the proposal unless it has already responded to a prepare
// request having a number greater than n.
func (a *Acceptor) ReceivePropose(proposal *Proposal) (*Proposal, error) {
	// Do we already have a promise for this proposal
	promised, ok := a.promisedProposals[proposal.Key]

	// Ignore lesser numbered proposals
	if ok && promised.Number > proposal.Number {
		return nil, fmt.Errorf(
			"Already promised to accept %s which is > than requested %s",
			promised,
			proposal,
		)
	}

	// Unexpected proposal
	if ok && promised.Number < proposal.Number {
		return nil, fmt.Errorf("Received unexpected proposal %s", proposal)
	}

	// Accept the proposal
	a.acceptedProposals[proposal.Key] = proposal

	return proposal, nil
}
