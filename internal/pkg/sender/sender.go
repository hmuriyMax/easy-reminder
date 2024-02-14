package sender

type (
	Sender interface {
		Send(message *Message) (*Result, error)
	}

	Message struct {
		ID   string
		Text string
	}

	Result struct {
		Message string
		Status  SendStatus
	}

	SendStatus int
)

const (
	SendStatusUnknown SendStatus = iota
	SendStatusInternalError
	SendStatusBadRequest
	SendStatusClientError
	SendStatusSuccess
)
