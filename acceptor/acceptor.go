package acceptor

import (
	"fmt"
	"log"

	"github.com/RichardKnop/paxos/models"
)

// Acceptor ...
type Acceptor struct {
	ID                string
	Host              string
	Port              int
	promisedProposals map[string]*models.Proposal
	acceptedProposals map[string]*models.Proposal
}

// New returns new Acceptor instance
func New(ID, host string, port int) (*Acceptor, error) {
	return &Acceptor{
		ID:                ID,
		Host:              host,
		Port:              port,
		promisedProposals: make(map[string]*models.Proposal, 0),
		acceptedProposals: make(map[string]*models.Proposal, 0),
	}, nil
}

// ToString returns a human readable representation
func (a *Acceptor) ToString() string {
	return fmt.Sprintf("Acceptor %s (%s:%d)", a.ID, a.Host, a.Port)
}

// If an acceptor receives a prepare request with number n greater
// than that of any prepare request to which it has already responded,
// then it responds to the request with a promise not to accept any more
// proposals numbered less than n and with the highest-numbered proposal
// (if any) that it has accepted.
func (a *Acceptor) receivePrepare(proposal *models.Proposal) (*models.Proposal, error) {
	// Do we already have a promise for this proposal
	promised, ok := a.promisedProposals[proposal.Key]

	// Ignore lesser or equally numbered proposals
	if ok && promised.Number >= proposal.Number {
		return nil, fmt.Errorf(
			"%s promised to accept %s which is >= than requested %s",
			a.ToString(),
			promised.ToString(),
			proposal.ToString(),
		)
	}

	// Promise to accept the proposal
	a.promisedProposals[proposal.Key] = proposal
	log.Printf(
		"%s promises to accept proposal %s",
		a.ToString(),
		proposal.ToString(),
	)

	return proposal, nil
}

// If an acceptor receives an accept request for a proposal numbered
// n, it accepts the proposal unless it has already responded to a prepare
// request having a number greater than n.
func (a *Acceptor) receiveProposal(proposal *models.Proposal) (*models.Proposal, error) {
	// Do we already have a promise for this proposal
	promised, ok := a.promisedProposals[proposal.Key]

	// Ignore lesser or equally numbered proposals
	if ok && promised.Number >= proposal.Number {
		return nil, fmt.Errorf(
			"%s promised to accept %s which is >= than requested %s",
			a.ToString(),
			promised.ToString(),
			proposal.ToString(),
		)
	}

	// Unexpected proposal
	if ok && promised.Number < proposal.Number {
		return nil, fmt.Errorf(
			"%s received unexpected proposal %s",
			a.ToString(),
			proposal.ToString(),
		)
	}

	// Accept the proposal
	a.acceptedProposals[proposal.Key] = proposal
	log.Printf(
		"%s accepted proposal %s",
		a.ToString(),
		proposal.ToString(),
	)

	return proposal, nil
}
