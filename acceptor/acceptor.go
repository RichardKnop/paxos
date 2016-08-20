package acceptor

import (
	"fmt"
	"log"

	"github.com/RichardKnop/paxos/models"
)

// Acceptor ...
type Acceptor struct {
	id       string
	host     string
	port     int
	promised *models.Proposal
	accepted *models.Proposal
}

// New returns new Acceptor instance
func New(id, host string, port int) (*Acceptor, error) {
	return &Acceptor{
		id:   id,
		host: host,
		port: port,
	}, nil
}

// Prepare ...
func (a *Acceptor) Prepare(proposal *models.Proposal, reply *models.Proposal) error {
	proposal, err := a.receivePrepare(proposal)
	if err != nil {
		return err
	}
	*reply = *proposal
	return nil
}

// Propose ...
func (a *Acceptor) Propose(proposal *models.Proposal, reply *models.Proposal) error {
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
func (a *Acceptor) receivePrepare(proposal *models.Proposal) (*models.Proposal, error) {
	// Ignore lesser proposals
	if a.promised.Number >= proposal.Number {
		return nil, fmt.Errorf(
			"Acceptor already promised to accept %d which is >= than requested %d",
			a.promised.Number,
			proposal.Number,
		)
	}

	// Promise to accept the proposal
	log.Printf("Acceptor promises to accept proposal %d\n", proposal.Number)
	a.promised = proposal
	return a.promised, nil
}

// If an acceptor receives an accept request for a proposal numbered
// n, it accepts the proposal unless it has already responded to a prepare
// request having a number greater than n.
func (a *Acceptor) receiveProposal(proposal *models.Proposal) (*models.Proposal, error) {
	// Ignore lesser proposals
	if a.promised.Number > proposal.Number {
		return nil, fmt.Errorf(
			"Acceptor already promised to accept %d which is >= than requested %d",
			a.promised.Number,
			proposal.Number,
		)
	}

	// Unexpected proposal
	if a.promised.Number < proposal.Number {
		return nil, fmt.Errorf("Received unexpected proposal %d", proposal.Number)
	}

	// Accept the proposal
	log.Printf("Accepted proposal %d", proposal.Number)
	a.accepted = proposal
	return a.accepted, nil
}
