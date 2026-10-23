package repo

import (
	"encoding/json"
	"errors"
	"fmt"

	"server-alarm/api/internal/app/entity"
	"server-alarm/api/internal/pkg/storage"
)

const (
	_keyLast    = "climate:last-record" // storage key for last record
	_keyHistory = "climate:history"     // storage key for records history
)

// ClimateRepoCache represents a cache repo for entity.Climate.
type ClimateRepoCache struct {
	storageInst storage.KeyValue
	HistoryLen  int
}

// NewClimateRepoCache returns a new instance of ClimateRepoCache.
func NewClimateRepoCache(storageInst storage.KeyValue, historyLen int) *ClimateRepoCache {
	return &ClimateRepoCache{
		storageInst: storageInst,
		HistoryLen:  historyLen,
	}
}

// Get gets last climate record from storage.
func (r *ClimateRepoCache) Get() (*entity.Climate, error) {
	data, err := r.storageInst.Get(_keyLast)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, nil
	}
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

	if err := r.storageInst.Set(_keyLast, data, 0); err != nil {
		return fmt.Errorf("set: %w", err)
	}
	return nil
}

// GetHistory gets history with climate records from storage.
func (r *ClimateRepoCache) GetHistory() ([]*entity.Climate, error) {
	data, err := r.storageInst.Get(_keyHistory)
	if errors.Is(err, storage.ErrNotFound) {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("get history: %w", err)
	}

	history := []*entity.Climate{}
	if err := json.Unmarshal(data, &history); err != nil {
		return nil, fmt.Errorf("unmarshal from json: %w", err)
	}
	return history, nil
}

// SetHistory sets given record to history into storage.
func (r *ClimateRepoCache) SetHistory(record *entity.Climate) error {
	history, err := r.GetHistory()
	if err != nil {
		return fmt.Errorf("old history: %w", err)
	}

	switch {
	case len(history) == 0:
		history = make([]*entity.Climate, 0, r.HistoryLen)
		history = append(history, record)
	case len(history) < r.HistoryLen:
		history = append([]*entity.Climate{record}, history...)
	default:
		history = append([]*entity.Climate{record}, history[:r.HistoryLen-1]...)
	}

	data, err := json.Marshal(history)
	if err != nil {
		return fmt.Errorf("marshal to json: %w", err)
	}

	if err := r.storageInst.Set(_keyHistory, data, 0); err != nil {
		return fmt.Errorf("set history: %w", err)
	}
	return nil
}
