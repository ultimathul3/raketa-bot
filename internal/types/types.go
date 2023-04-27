package types

type Status string

const (
	TaskOpen     Status = "open"
	TaskClosed   Status = "closed"
	TaskDeclined Status = "declined"
	TaskUnknown  Status = "unknown"
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
