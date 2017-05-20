package paxos

import (
	"fmt"
)

// Proposal represents a proposed value for specified key,
// the actual agreed on value will be decided by the consensus algorithm
type Proposal struct {
	Key    string // key to identify the value
	Number int    // proposal number used to decide which proposal to promise
	Value  string // the actual value / data we store once proposal is accepted
}

// NewProposal creates a new instance of Proposal
func NewProposal(key, value string) *Proposal {
	return &Proposal{Key: key, Value: value}
}

// String returns a human readable representation
func (p *Proposal) String() string {
	return fmt.Sprintf("[%d, %s]", p.Number, p.Value)
}
