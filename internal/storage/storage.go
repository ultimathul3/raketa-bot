package storage

import "github.com/vanyaio/raketa-bot/internal/types"

type stateInfo struct {
	state    types.State
	callback types.Callback
}

type StateStorage struct {
	storage map[int64]stateInfo
}

func NewStateStorage() *StateStorage {
	return &StateStorage{
		storage: make(map[int64]stateInfo),
	}
}

func (s *StateStorage) GetState(ID int64) types.State {
	return s.storage[ID].state
}

func (s *StateStorage) GetCallback(ID int64) types.Callback {
	return s.storage[ID].callback
}

func (s *StateStorage) SetState(ID int64, state types.State, callback types.Callback) {
	s.storage[ID] = stateInfo{
		state:    state,
		callback: callback,
	}
}
