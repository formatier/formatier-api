package workflow

type sagaWorkflowSchema struct {
	ServiceName string `json:"service_name"`
	EventName   string `json:"event_name"`
	Status      string `json:"status"`
	Cancelable  bool   `json:"cancelable"`
}

const (
	WORKFLOW_STATUS_WAITING = "waiting"
	WORKFLOW_STATUS_SUCCESS = "success"
	WORKFLOW_STATUS_FAILED  = "failed"
)
