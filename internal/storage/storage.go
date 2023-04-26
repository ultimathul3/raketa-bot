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
	stateWithData := s.storage[userID]
	stateWithData.state = state
	s.storage[userID] = stateWithData
}

func (s *StateStorageWithData) SetStateWithID(userID int64, state types.State, id int64) {
	stateWithData := s.storage[userID]
	stateWithData.state = state
	stateWithData.id = id
	s.storage[userID] = stateWithData
}

func (s *StateStorageWithData) SetStateWithURL(userID int64, state types.State, url string) {
	stateWithData := s.storage[userID]
	stateWithData.state = state
	stateWithData.url = url
	s.storage[userID] = stateWithData
}
