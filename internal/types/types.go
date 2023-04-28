package types

type Status int
type Role int

const (
	TaskOpen Status = iota
	TaskClosed
	TaskDeclined
	TaskUnknown
)

const (
	AdminRole Role = iota
	RegularRole
	UnknownRole
)

type Task struct {
	Url    string
	Status Status
	UserID int64
	Price  uint64
}

const (
	UrlDataKey      = "url"
	UsernameDataKey = "username"
)
