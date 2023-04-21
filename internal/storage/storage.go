package storage

type State int64

const (
	Menu State = iota
	UrlInput
)

type callbackFunc func(ctx ...any)

type Storage interface {
	GetState(ID int64) State
	GetCallback(ID int64, state State) callbackFunc
	SetState(ID int64, state State, callback callbackFunc)
}

type pair struct {
	state    State
	callback callbackFunc
}

type StateStorage struct {
	storage map[int64]pair
}

func NewStateStorage() *StateStorage {
	return &StateStorage{
		storage: make(map[int64]pair),
	}
}

func (s *StateStorage) GetState(ID int64) State {
	return s.storage[ID].state
}

func (s *StateStorage) GetCallback(ID int64, state State) callbackFunc {
	return s.storage[ID].callback
}

func (s *StateStorage) SetState(ID int64, state State, callback callbackFunc) {
	s.storage[ID] = pair{
		state:    state,
		callback: callback,
	}
}
