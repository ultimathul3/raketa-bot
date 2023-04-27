package storage

import "github.com/vanyaio/raketa-bot/internal/types"

type stateWithData struct {
	state types.State
	data  map[string]any
}

type StateStorageWithData struct {
	storage map[int64]stateWithData
}

func NewStateStorageWithData() *StateStorageWithData {
	return &StateStorageWithData{
		storage: make(map[int64]stateWithData),
	}
}

func (s *StateStorageWithData) GetState(userID int64) (types.State, bool) {
	storage, ok := s.storage[userID]
	if !ok {
		return types.StateUnknown, false
	}

	return storage.state, true
}

func (s *StateStorageWithData) GetData(userID int64, key string) (any, bool) {
	storage, ok := s.storage[userID]
	if !ok {
		return nil, false
	}

	data, ok := storage.data[key]
	return data, ok
}

func (s *StateStorageWithData) SetState(userID int64, state types.State) {
	stateWithData := s.storage[userID]
	stateWithData.state = state
	s.storage[userID] = stateWithData
}

func (s *StateStorageWithData) SetStateWithData(userID int64, state types.State, key string, value any) {
	stateWithData := s.storage[userID]
	stateWithData.state = state

	if stateWithData.data == nil {
		stateWithData.data = make(map[string]any)
	}
	stateWithData.data[key] = value

	s.storage[userID] = stateWithData
}
