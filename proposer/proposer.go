package proposer

import (
	"log"
	"net/rpc"

	"github.com/RichardKnop/paxos/models"
)

// Proposer ...
type Proposer struct {
	id               string
	host             string
	port             int
	proposal         *models.Proposal
	acceptors        []string
	acceptorPromises map[string]*models.Proposal
}

// New returns new Proposer instance
func New(id, host string, port int, proposedValue string, acceptors []string) (*Proposer, error) {
	return &Proposer{
		id:   id,
		host: host,
		port: port,
		proposal: &models.Proposal{
			Value: proposedValue,
		},
		acceptors:        acceptors, // addresses
		acceptorPromises: make(map[string]*models.Proposal, len(acceptors)),
	}, nil
}

// Run ...
func (p *Proposer) Run() {
	// Stage 1: Prepare proposals until majority is reached
	for !p.majorityReached() {
		p.prepare()
	}
	log.Printf("Reached majority %d", p.majority())

	// Stage 2: Finalise proposal
	log.Printf(
		"Starting to propose [%d: %s]",
		p.proposal.Number,
		p.proposal.Value,
	)
	p.propose()
}

// A proposer chooses a new proposal number n and sends a request to
// each member of some set of acceptors, asking it to respond with:
// (a) A promise never again to accept a proposal numbered less than n, and
// (b) The proposal with the highest number less than n that it has accepted, if any.
func (p *Proposer) prepare() {
	p.proposal.Number++

	for i := 0; i < p.majority(); i++ {
		promised, err := p.sendPrepareRequest(p.acceptors[i], p.proposal)
		if err != nil {
			log.Printf("Send Prepared Proposal: %v", err)
			continue
		}
		log.Printf(
			"Acceptor %s returned promise [%d, %s]",
			p.acceptors[i],
			promised.Number,
			promised.Value,
		)

		previusPromise := p.acceptorPromises[p.acceptors[i]]
		if previusPromise.Number < promised.Number {
			log.Printf(
				"Received a new promise [%d, %s]",
				promised.Number,
				promised.Value,
			)
			p.acceptorPromises[p.acceptors[i]] = promised

			// Update the proposal to the one with bigger number
			if promised.Number > p.proposal.Number {
				log.Printf(
					"Updating the proposal to [%d, %s]",
					promised.Number,
					promised.Value,
				)
				p.proposal = promised
			}
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
	err = client.Call("Acceptor.Prepare", proposal, &reply)
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
	err = client.Call("Acceptor.Propose", proposal, &reply)
	if err != nil {
		return nil, err
	}

	return reply, nil
}
