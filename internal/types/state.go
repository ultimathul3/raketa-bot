package types

type State int64

const (
	Menu State = iota
	UrlInput
	IdInput
)

type Callback func(string)
