package types

type State int64

const (
	Menu State = iota
	CreateTaskUrlInput
	CreateTaskPriceInput
	DeleteTaskUrlInput
	AssignWorkerUrlInput
	AssignWorkerUsernameInput
	CloseTaskUrlInput
	StateUnknown
)
