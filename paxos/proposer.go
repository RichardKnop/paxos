package paxos

import (
	"log"
)

// Proposer proposes new values to acceptors
type Proposer struct {
	acceptorClients  []AcceptorClientInterface
	acceptorPromises map[string]map[string]*Proposal
}

// NewProposer creates a new proposer instance
func NewProposer(acceptorClients []AcceptorClientInterface) *Proposer {
	acceptorPromises := make(map[string]map[string]*Proposal, len(acceptorClients))
	for _, acceptorClient := range acceptorClients {
		acceptorPromises[acceptorClient.GetName()] = make(map[string]*Proposal)
	}
	return &Proposer{
		acceptorClients:  acceptorClients,
		acceptorPromises: acceptorPromises,
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
		acceptorClient := p.acceptorClients[i]

		promise, err := acceptorClient.SendPrepare(proposal)
		if err != nil {
			continue
		}

		// Get the previous promise
		previousPromise, ok := p.acceptorPromises[acceptorClient.GetName()][proposal.Key]

		// Previous promise is equal or greater than the new proposal, continue
		if ok && previousPromise.Number >= promise.Number {
			continue
		}

		// Save the new promise
		p.acceptorPromises[acceptorClient.GetName()][proposal.Key] = promise

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
	for _, acceptorClient := range p.acceptorClients {
		accepted, err := acceptorClient.SendPropose(proposal)
		if err != nil {
			continue
		}

		log.Printf("%s has accepted the proposal %s", acceptorClient.GetName(), accepted)
	}

	return nil
}

// majority returns simple majority of acceptor nodes
func (p *Proposer) majority() int {
	return len(p.acceptorClients)/2 + 1
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
