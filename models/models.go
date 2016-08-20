package models

import (
	"fmt"
)

// Proposal ...
type Proposal struct {
	Number int
	Value  string
}

// ToString returns a human readable representation
func (p *Proposal) ToString() string {
	return fmt.Sprintf("[%d, %s]", p.Number, p.Value)
}
