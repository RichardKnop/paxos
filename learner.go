package paxos

// Learner ...
type Learner struct {
	ID string
}

// NewLearner creates a new learner instance
func NewLearner(id string) (*Learner, error) {
	return &Learner{
		ID: id,
	}, nil
}
