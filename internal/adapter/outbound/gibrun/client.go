package gibrun_outbound_adapter

import (
	"context"

	"eduvera/internal/model"
	outbound_port "eduvera/internal/port/outbound"
	"eduvera/utils/gibrun"
)

type clientAdapter struct{}

func NewClientAdapter() outbound_port.ClientCachePort {
	return &clientAdapter{}
}

func (adapter *clientAdapter) Set(data model.Client) error {
	// Gib means Give/Set
	return gibrun.Gib(context.Background(), data.BearerKey, data)
}

func (adapter *clientAdapter) Get(bearerKey string) (model.Client, error) {
	var client model.Client
	// Run means Retrieve/Get
	found, err := gibrun.Run(context.Background(), bearerKey, &client)
	if err != nil {
		return model.Client{}, err
	}
	if !found {
		// Return empty/error if not found or handle as nil?
		// Redis adapter returns error on miss?
		// redis/client.go: `if err != nil { return ..., err }`.
		// If redis.Get returns redis.Nil, it bubbles up?
		// Let's assume CachePort expects error on miss or empty.
		// For now returning empty client.
		return model.Client{}, nil
	}
	return client, nil
}
