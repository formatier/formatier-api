package service

type UniversalSagaSagaMetadataSchema struct {
	SagaID    string `json:"saga_id"`
	Timestamp int64  `json:"timestamp"`
}

type UniversalSagaTimeoutPolicySchema struct {
	Mode     string `json:"mode"`
	Duration int    `json:"duration"`
}

type UniversalSagaTriggerEventSchema struct {
	TimeoutPolicy UniversalSagaTimeoutPolicySchema `json:"timeout_policy"`
	Metadata      UniversalSagaSagaMetadataSchema  `json:"metadata"`
	Payload       any                              `json:"payload"`
}

type UniversalSagaReciverSchema struct {
	Metadata  UniversalSagaSagaMetadataSchema `json:"metadata"`
	Response  string                          `json:"payload"`
	Status    string                          `json:"status"`
	Timestamp int64                           `json:"timestamp"`
}
