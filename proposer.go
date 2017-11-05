package paxos

import (
	"fmt"
	"log"
	"net/rpc"
	"time"
)

// Proposer proposes new values to acceptors
type Proposer struct {
	Node
	acceptorURLs     []string
	acceptorPromises map[string]map[string]*Proposal
	prepareFib       func() int
	proposeFib       func() int
}

// NewProposer creates a new proposer instance
func NewProposer(id, host string, port int, acceptorURLs []string) *Proposer {
	acceptorURLs = append(acceptorURLs, fmt.Sprintf("%s:%d", host, port))
	acceptorPromises := make(map[string]map[string]*Proposal, len(acceptorURLs))
	for _, acceptorURL := range acceptorURLs {
		acceptorPromises[acceptorURL] = make(map[string]*Proposal)
	}
	return &Proposer{
		Node:             NewNode(id, host, port),
		acceptorURLs:     acceptorURLs,
		acceptorPromises: acceptorPromises,
		prepareFib:       Fibonacci(),
		proposeFib:       Fibonacci(),
	}
}

// Propose sends a proposal request to the peers
func (p *Proposer) Propose(proposal *Proposal) error {
	// Stage 1: Prepare proposals until majority is reached
	for !p.majorityReached(proposal) {
		if err := p.prepare(proposal); err != nil {
			return err
		}
	}
	log.Printf("Reached majority %d", p.majority())

	// Stage 2: Propose the value agreed on by majority of acceptors
	log.Printf("Sending final proposal %s", proposal)

	return p.propose(proposal)
}

// A proposer chooses a new proposal number n and sends a request to
// each member of some set of acceptors, asking it to respond with:
// (a) A promise never again to accept a proposal numbered less than n, and
// (b) The proposal with the highest number less than n that it has accepted, if any.
func (p *Proposer) prepare(proposal *Proposal) error {
	// Increment the proposal number
	proposal.Number++

	for i := 0; i < p.majority(); i++ {
		acceptorURL := p.acceptorURLs[i]

		promise, err := sendRPCRequest(acceptorURL, PrepareServiceMethod, proposal)
		if err != nil {
			// Use fibonacci sequence to space out retry attempts
			waitMs := p.prepareFib() * 1
			<-time.After(time.Duration(waitMs) * time.Millisecond)

			continue
		}

		// Reset the fibonacci sequence
		p.prepareFib = Fibonacci()

		log.Printf("%s promises to accept proposal %s", acceptorURL, promise)

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
			log.Printf("Updating the proposal to %s", promise)
			proposal = promise
		}
	}

	return nil
}

// If the proposer receives the requested responses from a majority of
// the acceptors, then it can issue a proposal with number n and value
// v, where v is the value of the highest-numbered proposal among the
// responses, or is any value selected by the proposer if the responders
// reported no proposals.
func (p *Proposer) propose(proposal *Proposal) error {
	for _, acceptorURL := range p.acceptorURLs {
		accepted, err := sendRPCRequest(acceptorURL, ProposeServiceMethod, proposal)
		if err != nil {
			// Use fibonacci sequence to space out retry attempts
			waitMs := p.proposeFib() * 1
			<-time.After(time.Duration(waitMs) * time.Millisecond)

			continue
		}

		// Reset the fibonacci sequence
		p.proposeFib = Fibonacci()

		log.Printf("%s has accepted the proposal %s", acceptorURL, accepted)
	}

	return nil
}

// majority returns simple majority of acceptor nodes
func (p *Proposer) majority() int {
	return len(p.acceptorURLs)/2 + 1
}

// majorityReached returns true if number of matching promises from acceptors
// is equal or greater than simple majority of acceptor nodes
func (p *Proposer) majorityReached(proposal *Proposal) bool {
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
func sendRPCRequest(address, serviceMethod string, proposal *Proposal) (*Proposal, error) {
	client, err := rpc.DialHTTP("tcp", address)
	if err != nil {
		return nil, err
	}

	var reply *Proposal
	err = client.Call(serviceMethod, proposal, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
