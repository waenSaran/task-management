package utils

type Status string

const (
	Todo       Status = "TODO"
	InProgress Status = "IN_PROGRESS"
	Done       Status = "DONE"
	Archive    Status = "ARCHIVE"
)
