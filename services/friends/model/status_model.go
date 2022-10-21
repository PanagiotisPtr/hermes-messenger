package model

type Status string

const (
	StatusAccepted Status = "accepted"
	StatusPending         = "pending"
	StatusRejected        = "rejected"
)
