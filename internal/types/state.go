package types

type State int64

const (
	Menu State = iota
	CreateTaskUrlInput
	DeleteTaskUrlInput
	AssignWorkerUrlInput
	AssignWorkerIdInput
	CloseTaskUrlInput
	StateUnknown
)
