package saga

type SagaWorkflowSchema struct {
	ServiceName string `json:"service"`
	Command     string `json:"command"`
	Status      string `json:"status"`

	Cancelable     bool `json:"cancelable"`
	StartTimestamp int64
	DoneTimestamp  int64 `json:"timestamp"`
}

const (
	CMD_STATUS_WAITING = "waiting"
	CMD_STATUS_PENDING = "pending"
	CMD_STATUS_ASYNC   = "async"

	CMD_STATUS_SUCCESS = "success"
	CMD_STATUS_FAILED  = "failed"
)

const (
	CMD_STATE_NEXT     = "next"
	CMD_STATE_ROLLBACK = "rollback"
)

const (
	TIMEOUT_MODE_CANCEL = "cancel"
	TIMEOUT_MODE_RETRY  = "retry"
)
