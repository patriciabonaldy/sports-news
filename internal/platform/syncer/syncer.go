package syncer

import "context"

// Syncer represents a synchronization process.
type Syncer interface {
	// Sync synchronizes data among the providers.
	Sync(context.Context) error
}
