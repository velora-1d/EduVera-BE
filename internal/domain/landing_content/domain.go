package landing_content

import (
	"context"
	outbound_port "prabogo/internal/port/outbound"
)

type Domain struct {
	dbPort outbound_port.LandingContentDatabasePort
}

func NewDomain(dbPort outbound_port.DatabasePort) *Domain {
	return &Domain{
		dbPort: dbPort.LandingContent(),
	}
}

func (d *Domain) Get(ctx context.Context, key string) (map[string]interface{}, error) {
	return d.dbPort.Get(ctx, key)
}

func (d *Domain) Set(ctx context.Context, key string, value map[string]interface{}) error {
	// Future: validation based on key?
	return d.dbPort.Set(ctx, key, value)
}
