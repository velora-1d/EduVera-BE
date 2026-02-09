package outbound_port

import (
	"context"
)

// LandingContentDatabasePort defines database operations for dynamic landing page content
type LandingContentDatabasePort interface {
	// Get retrieves content by key
	Get(ctx context.Context, key string) (map[string]interface{}, error)
	// Set upserts content by key
	Set(ctx context.Context, key string, value map[string]interface{}) error
}
