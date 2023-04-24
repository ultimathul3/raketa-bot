package storage

import "github.com/vanyaio/raketa-bot/internal/types"

type pair struct {
	state    types.State
	callback types.CallbackFunc
}

type StateStorage struct {
	storage map[int64]pair
}

func NewStateStorage() *StateStorage {
	return &StateStorage{
		storage: make(map[int64]pair),
	}
}

func (s *StateStorage) GetState(ID int64) types.State {
	return s.storage[ID].state
}

func (s *StateStorage) GetCallback(ID int64, state types.State) types.CallbackFunc {
	return s.storage[ID].callback
}

func (s *StateStorage) SetState(ID int64, state types.State, callback types.CallbackFunc) {
	s.storage[ID] = pair{
		state:    state,
		callback: callback,
	}
}
