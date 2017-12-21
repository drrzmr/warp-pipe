package state

const (
	// Idle state
	Idle = "state.Idle" // -> start
	// Start stare
	Start = "state.Start" // -> started
	// Started state
	Started = "state.Started" // -> store, publish
	// Store state
	Store = "state.Store" // -> stored
	// Stored state
	Stored = "state.Stored" // -> publish, store
	// Publish state
	Publish = "state.Publish" // -> idle
)
