package acceptor

import (
	"fmt"
	"log"

	"github.com/RichardKnop/paxos/models"
)

// Acceptor ...
type Acceptor struct {
	ID               string
	Host             string
	Port             int
	promisedProposal *models.Proposal
	acceptedProposal *models.Proposal
}

// New returns new Acceptor instance
func New(ID, host string, port int) (*Acceptor, error) {
	return &Acceptor{
		ID:   ID,
		Host: host,
		Port: port,
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
	// Ignore lesser proposals
	if a.promisedProposal.Number >= proposal.Number {
		return nil, fmt.Errorf(
			"%s promised to accept %s which is >= than requested %s",
			a.ToString(),
			a.promisedProposal.ToString(),
			proposal.ToString(),
		)
	}

	// Promise to accept the proposal
	a.promisedProposal = proposal
	log.Printf(
		"%s promises to accept proposal %s",
		a.ToString(),
		a.promisedProposal.ToString(),
	)

	return a.promisedProposal, nil
}

// If an acceptor receives an accept request for a proposal numbered
// n, it accepts the proposal unless it has already responded to a prepare
// request having a number greater than n.
func (a *Acceptor) receiveProposal(proposal *models.Proposal) (*models.Proposal, error) {
	// Ignore lesser proposals
	if a.promisedProposal.Number > proposal.Number {
		return nil, fmt.Errorf(
			"%s promised to accept %s which is >= than requested %s",
			a.ToString(),
			a.promisedProposal.ToString(),
			proposal.ToString(),
		)
	}

	// Unexpected proposal
	if a.promisedProposal.Number < proposal.Number {
		return nil, fmt.Errorf(
			"%s received unexpected proposal %d",
			a.ToString(),
			proposal.ToString(),
		)
	}

	// Accept the proposal
	a.acceptedProposal = proposal
	log.Printf(
		"%s accepted proposal %s",
		a.ToString(),
		a.acceptedProposal.ToString(),
	)

	return a.acceptedProposal, nil
}
