package proposer

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/RichardKnop/paxos/acceptor"
	"github.com/RichardKnop/paxos/models"
)

// Proposer ...
type Proposer struct {
	ID               string
	Host             string
	Port             int
	proposal         *models.Proposal
	acceptorURLs     []string
	acceptorPromises map[string]*models.Proposal
}

// New returns new Proposer instance
func New(ID, host string, port int, proposedValue string, acceptorURLs []string) (*Proposer, error) {
	acceptorURLs = append(acceptorURLs, fmt.Sprintf("%s:%d", host, port))
	return &Proposer{
		ID:   ID,
		Host: host,
		Port: port,
		proposal: &models.Proposal{
			Value: proposedValue,
		},
		acceptorURLs:     acceptorURLs,
		acceptorPromises: make(map[string]*models.Proposal, len(acceptorURLs)),
	}, nil
}

// ToString returns a human readable representation
func (p *Proposer) ToString() string {
	return fmt.Sprintf("Proposer %s (%s:%d)", p.ID, p.Host, p.Port)
}

// Run fires off the proposal process
func (p *Proposer) Run() {
	// Stage 1: Prepare proposals until majority is reached
	for !p.majorityReached() {
		p.prepare()
	}
	log.Printf("%s reached majority %d", p.ToString(), p.majority())

	// Stage 2: Propose the value agreed on by majority of acceptors
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
		acceptorURL := p.acceptorURLs[i]

		promise, err := sendRPCRequest(acceptorURL, acceptor.PrepareServiceMethod, p.proposal)
		if err != nil {
			log.Printf("Prepare request failed: %v", err)
			continue
		}
		log.Printf("%s received promise %s from %s", p.ToString(), promise.ToString(), acceptorURL)

		// Get the previous promise
		previousPromise := p.acceptorPromises[acceptorURL]

		// Previous promise is equal or greater than the new proposal, ignore
		if previousPromise != nil && previousPromise.Number >= promise.Number {
			continue
		}

		// Log the new promise
		p.acceptorPromises[acceptorURL] = promise

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
	for _, acceptorURL := range p.acceptorURLs {
		accepted, err := sendRPCRequest(acceptorURL, acceptor.ProposeServiceMethod, p.proposal)
		if err != nil {
			log.Printf("Propose request failed: %v", err)
			continue
		}
		log.Printf("Accepted proposal [%d, %s]", accepted.Number, p.proposal.Value)
	}
}

// majority returns simple majority of acceptor nodes
func (p *Proposer) majority() int {
	return len(p.acceptorURLs)/2 + 1
}

// majorityReached returns true if number of matching promises from acceptors
// is equal or greater than simple majority of acceptor nodes
func (p *Proposer) majorityReached() bool {
	matches := 0
	for _, promised := range p.acceptorPromises {
		if promised.Number == p.proposal.Number {
			matches++
		}
	}
	return matches >= p.majority()
}

// sendRPCRequest is a generic function to make RPC call to acceptor's
// prepare or promise service methods
func sendRPCRequest(address, serviceMethod string, proposal *models.Proposal) (*models.Proposal, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	var reply *models.Proposal
	err = client.Call(serviceMethod, proposal, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
