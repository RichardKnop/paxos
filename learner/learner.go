package learner

// Learner ...
type Learner struct {
	id string
}

// New returns new Learner instance
func New(id string) (*Learner, error) {
	return &Learner{
		id: id,
	}, nil
}
