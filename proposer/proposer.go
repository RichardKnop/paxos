package proposer

import (
	"fmt"
	"log"
	"net/rpc"

	"github.com/RichardKnop/paxos/acceptor"
	"github.com/RichardKnop/paxos/models"
)

// Proposer proposes new values to acceptors
type Proposer struct {
	ID               string
	Host             string
	Port             int
	acceptorURLs     []string
	acceptorPromises map[string]map[string]*models.Proposal
}

// New returns new Proposer instance
func New(theID, host string, port int, acceptorURLs []string) (*Proposer, error) {
	acceptorURLs = append(acceptorURLs, fmt.Sprintf("%s:%d", host, port))
	acceptorPromises := make(map[string]map[string]*models.Proposal, len(acceptorURLs))
	for _, acceptorURL := range acceptorURLs {
		acceptorPromises[acceptorURL] = make(map[string]*models.Proposal, 0)
	}
	return &Proposer{
		ID:               theID,
		Host:             host,
		Port:             port,
		acceptorURLs:     acceptorURLs,
		acceptorPromises: acceptorPromises,
	}, nil
}

// ToString returns a human readable representation
func (p *Proposer) ToString() string {
	return fmt.Sprintf("Proposer %s (%s:%d)", p.ID, p.Host, p.Port)
}

// Propose sends a proposal request to the peers
func (p *Proposer) Propose(proposal *models.Proposal) {
	// Stage 1: Prepare proposals until majority is reached
	for !p.majorityReached(proposal) {
		p.prepare(proposal)
	}
	log.Printf("%s reached majority %d", p.ToString(), p.majority())

	// Stage 2: Propose the value agreed on by majority of acceptors
	log.Printf(
		"%s is proposing final proposal [%s=%s (proposal number: %d)]",
		p.ToString(),
		proposal.Key,
		proposal.Value,
		proposal.Number,
	)
	p.propose(proposal)
}

// A proposer chooses a new proposal number n and sends a request to
// each member of some set of acceptors, asking it to respond with:
// (a) A promise never again to accept a proposal numbered less than n, and
// (b) The proposal with the highest number less than n that it has accepted, if any.
func (p *Proposer) prepare(proposal *models.Proposal) {
	// Increment the proposal number
	proposal.Number++

	for i := 0; i < p.majority(); i++ {
		acceptorURL := p.acceptorURLs[i]

		promise, err := sendRPCRequest(acceptorURL, acceptor.PrepareServiceMethod, proposal)
		if err != nil {
			log.Printf("Prepare request failed: %v", err)
			continue
		}
		log.Printf("%s received promise %s from %s", p.ToString(), promise.ToString(), acceptorURL)

		// Get the previous promise
		previousPromise, ok := p.acceptorPromises[acceptorURL][proposal.Key]

		// Previous promise is equal or greater than the new proposal, continue
		if ok && previousPromise.Number >= promise.Number {
			continue
		}

		// Save the new promise
		p.acceptorPromises[acceptorURL][proposal.Key] = promise

		// Update the proposal to the one with bigger number
		if promise.Number > proposal.Number {
			log.Printf("%s updated the proposal to %s", p.ToString(), promise.ToString())
			proposal = promise
		}
	}
}

// If the proposer receives the requested responses from a majority of
// the acceptors, then it can issue a proposal with number n and value
// v, where v is the value of the highest-numbered proposal among the
// responses, or is any value selected by the proposer if the responders
// reported no proposals.
func (p *Proposer) propose(proposal *models.Proposal) {
	for _, acceptorURL := range p.acceptorURLs {
		accepted, err := sendRPCRequest(acceptorURL, acceptor.ProposeServiceMethod, proposal)
		if err != nil {
			log.Printf("Propose request failed: %v", err)
			continue
		}
		log.Printf("Accepted proposal [%d, %s]", accepted.Number, proposal.Value)
	}
}

// majority returns simple majority of acceptor nodes
func (p *Proposer) majority() int {
	return len(p.acceptorURLs)/2 + 1
}

// majorityReached returns true if number of matching promises from acceptors
// is equal or greater than simple majority of acceptor nodes
func (p *Proposer) majorityReached(proposal *models.Proposal) bool {
	var matches = 0

	// Iterate over promised values for each acceptor
	for _, promiseMap := range p.acceptorPromises {
		// Skip if the acceptor has not yet promised a proposal for this key
		promised, ok := promiseMap[proposal.Key]
		if !ok {
			continue
		}

		// If the promised and proposal number is the same, increment matches count
		if promised.Number == proposal.Number {
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
