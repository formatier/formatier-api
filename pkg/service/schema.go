package service

type UniversalSagaMetadataSchema struct {
	SagaID    string `json:"saga_id"`
	Timestamp int64  `json:"timestamp"`
}

func NewSagaTimeoutPolicy(duration int, mode string) *UniversalSagaTimeoutPolicySchema {
	if mode == "" {
		mode = TIMEOUT_MODE_CANCEL
	}
	return &UniversalSagaTimeoutPolicySchema{
		Mode:     mode,
		Duration: duration,
	}
}

type UniversalSagaTimeoutPolicySchema struct {
	Mode     string `json:"mode"`
	Duration int    `json:"duration"`
}

const (
	TIMEOUT_MODE_CANCEL = "cancel"
	TIMEOUT_MODE_RETRY  = "retry"
)

type UniversalSagaTriggerEventSchema struct {
	TimeoutPolicy UniversalSagaTimeoutPolicySchema `json:"timeout_policy"`
	Metadata      UniversalSagaMetadataSchema      `json:"metadata"`
	Payload       any                              `json:"payload"`
}

type UniversalSagaReciverSchema struct {
	Metadata  UniversalSagaMetadataSchema `json:"metadata"`
	Response  any                         `json:"payload"`
	Status    string                      `json:"status"`
	Timestamp int64                       `json:"timestamp"`
}

type UniversalSagaError struct {
	Message  string                      `json:"message"`
	Metadata UniversalSagaMetadataSchema `json:"metadata"`
}
