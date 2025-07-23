package orchestrator

import (
	"formatier-api/pkg/service"
)

type CommandRequestSchema struct {
	CommandID     string                                   `json:"command_id"`
	Payload       any                                      `json:"payload"`
	TimeoutPolicy service.UniversalSagaTimeoutPolicySchema `json:"timeout_policy"`
}

type CommandResponseSchema struct {
	CommandID string `json:"command_id"`
	Response  any    `json:"payload"`
	Status    string `json:"status"`
}

const (
	TIMEOUT_MODE_FUTURE   = "future"
	TIMEOUT_MODE_ROLLBACK = "rollback"
)
