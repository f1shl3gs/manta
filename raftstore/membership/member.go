package membership

type Member struct {
	ID uint64

	// Addresses is the list of peers in the raft cluster
	Address string `json:"address"`

	// Learner indicates if the member is raft learner
	Learner bool `json:"learner,omitempty"`
}
