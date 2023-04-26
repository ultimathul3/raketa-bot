package storage

import "github.com/vanyaio/raketa-bot/internal/types"

type stateWithData struct {
	state types.State
	url   string
	id    int64
}

type StateStorageWithData struct {
	storage map[int64]stateWithData
}

func NewStateStorageWithData() *StateStorageWithData {
	return &StateStorageWithData{
		storage: make(map[int64]stateWithData),
	}
}

func (s *StateStorageWithData) GetState(userID int64) types.State {
	return s.storage[userID].state
}

func (s *StateStorageWithData) GetURL(userID int64) string {
	return s.storage[userID].url
}

func (s *StateStorageWithData) GetID(userID int64) int64 {
	return s.storage[userID].id
}

func (s *StateStorageWithData) SetState(userID int64, state types.State) {
	s.storage[userID] = stateWithData{
		state: state,
	}
}

func (s *StateStorageWithData) SetStateWithID(userID int64, state types.State, id int64) {
	s.storage[userID] = stateWithData{
		state: state,
		id:    id,
	}
}

func (s *StateStorageWithData) SetStateWithURL(userID int64, state types.State, url string) {
	s.storage[userID] = stateWithData{
		state: state,
		url:   url,
	}
}
