package handler

import (
	"context"
	"formatier-api/pkg/service"
	"orchestrator-server/internal/workflow"
)

func main() {
	we := workflow.NewWorkflowEngine()
	we.RouteCommand(
		"GetUserById",
		func(ctx context.Context, s *workflow.ServiceCaller) (any, error) {
			s.FutureCall(ctx, "", "", "", false)
			return nil, nil
		},
		&service.UniversalSagaTimeoutPolicySchema{},
	)
}
