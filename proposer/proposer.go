package proposer

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/RichardKnop/paxos/models"
)

// Proposer ...
type Proposer struct {
	ID               string
	Host             string
	Port             int
	proposal         *models.Proposal
	acceptors        []string
	acceptorPromises map[string]*models.Proposal
}

// New returns new Proposer instance
func New(ID, host string, port int, proposedValue string, acceptors []string) (*Proposer, error) {
	return &Proposer{
		ID:   ID,
		Host: host,
		Port: port,
		proposal: &models.Proposal{
			Value: proposedValue,
		},
		acceptors:        acceptors, // addresses
		acceptorPromises: make(map[string]*models.Proposal, len(acceptors)),
	}, nil
}

// ToString returns a human readable representation
func (p *Proposer) ToString() string {
	return fmt.Sprintf("Proposer %s (%s:%d)", p.ID, p.Host, p.Port)
}

// Run ...
func (p *Proposer) Run() {
	// Stage 1: Prepare proposals until majority is reached
	for !p.majorityReached() {
		p.prepare()
	}
	log.Printf("%s reached majority %d", p.ToString(), p.majority())

	// Stage 2: Finalise proposal
	log.Printf("%s is proposing final proposal [%d: %s]", p.ToString(), p.proposal.Number, p.proposal.Value)
	p.propose()
}

// A proposer chooses a new proposal number n and sends a request to
// each member of some set of acceptors, asking it to respond with:
// (a) A promise never again to accept a proposal numbered less than n, and
// (b) The proposal with the highest number less than n that it has accepted, if any.
func (p *Proposer) prepare() {
	p.proposal.Number++

	for i := 0; i < p.majority(); i++ {
		acceptor := p.acceptors[i]

		promise, err := p.sendPrepareRequest(acceptor, p.proposal)
		if err != nil {
			log.Printf("Send Prepared Proposal: %v", err)
			continue
		}
		log.Printf("%s received promise %s from %s", p.ToString(), promise.ToString(), acceptor)

		// Get the previous promise
		previousPromise := p.acceptorPromises[acceptor]

		// Previous promise is euqual or greater than the new proposal, ignore
		if previousPromise != nil && promise.Number >= previousPromise.Number {
			continue
		}

		// Log the new promise
		p.acceptorPromises[acceptor] = promise

		// Update the proposal to the one with bigger number
		if promise.Number > p.proposal.Number {
			log.Printf("%s updated the proposal to %s", p.ToString(), promise.ToString())
			p.proposal = promise
		}
	}
}

// If the proposer receives the requested responses from a majority of
// the acceptors, then it can issue a proposal with number n and value
// v, where v is the value of the highest-numbered proposal among the
// responses, or is any value selected by the proposer if the responders
// reported no proposals.
func (p *Proposer) propose() {
	for _, acceptor := range p.acceptors {
		accepted, err := p.sendProposeRequest(acceptor, p.proposal)
		if err != nil {
			log.Printf("Send Prepared Proposal: %v", err)
			continue
		}
		log.Printf(
			"Accepted proposal [%d, %s]",
			accepted.Number,
			p.proposal.Value,
		)
	}
}

func (p *Proposer) majority() int {
	return len(p.acceptors)/2 + 1
}

func (p *Proposer) majorityReached() bool {
	m := 0
	for _, promised := range p.acceptorPromises {
		if promised.Number == p.proposal.Number {
			m++
		}
	}
	if m >= p.majority() {
		return true
	}
	return false
}

// sendPrepareRequest ...
func (p *Proposer) sendPrepareRequest(address string, proposal *models.Proposal) (*models.Proposal, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	var reply *models.Proposal
	err = client.Call("RPC.Prepare", proposal, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}

// sendProposeRequest ...
func (p *Proposer) sendProposeRequest(address string, proposal *models.Proposal) (*models.Proposal, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	var reply *models.Proposal
	err = client.Call("RPC.Propose", proposal, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
