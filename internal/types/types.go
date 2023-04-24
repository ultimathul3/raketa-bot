package types

type State int64

const (
	Menu State = iota
	UrlInput
)

type Callback func(string)
