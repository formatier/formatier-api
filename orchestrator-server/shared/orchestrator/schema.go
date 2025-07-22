package orchestrator

import (
	"formatier-api/pkg/service"
)

type CommandRequestSchema struct {
	CommandID     string                                   `json:"command_id"`
	Payload       string                                   `json:"payload"`
	TimeoutPolicy service.UniversalSagaTimeoutPolicySchema `json:"timeout_policy"`
}

var (
	TIMEOUT_MODE_PROMISE  = "promise"
	TIMEOUT_MODE_ROLLBACK = "rollback"
)
