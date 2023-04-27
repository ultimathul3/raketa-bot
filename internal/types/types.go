package types

type Status string

const (
	Open     Status = "open"
	Closed   Status = "closed"
	Declined Status = "declined"
	Unknown  Status = "unknown"
)

type Task struct {
	Url    string
	Status Status
	UserID int64
}

const (
	UrlData = "url"
	IdData  = "id"
)
