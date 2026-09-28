package repo

import (
	"encoding/json"
	"fmt"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/pkg/storage"
)

const _key = "climate:last-record" // storage key

// ClimateRepoCache represents a cache repo for entity.Climate.
type ClimateRepoCache struct {
	storageInst storage.KeyValue
}

// NewClimateRepoCache returns a new instance of ClimateRepoCache.
func NewClimateRepoCache(storageInst storage.KeyValue) *ClimateRepoCache {
	return &ClimateRepoCache{
		storageInst: storageInst,
	}
}

// Get gets last climate record from storage.
func (r *ClimateRepoCache) Get() (*entity.Climate, error) {
	data, err := r.storageInst.Get(_key)
	if err != nil {
		return nil, fmt.Errorf("get: %w", err)
	}

	record := &entity.Climate{}
	if err := json.Unmarshal(data, &record); err != nil {
		return nil, fmt.Errorf("unmarshal from json: %w", err)
	}
	return record, nil
}

// Set sets given climate record as last climate record into storage.
func (r *ClimateRepoCache) Set(record *entity.Climate) error {
	data, err := json.Marshal(record)
	if err != nil {
		return fmt.Errorf("marshal to json: %w", err)
	}

	if err := r.storageInst.Set(_key, data, 0); err != nil {
		return fmt.Errorf("set: %w", err)
	}
	return nil
}
