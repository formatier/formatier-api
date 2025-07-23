package workflow

import (
	"context"
	"formatier-api/pkg/service"

	"github.com/nats-io/nats.go"
)

func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		commandMap: make(map[string]*Command),
	}
}

type WorkflowEngine struct {
	nc *nats.Conn

	commandMap map[string]*Command
}

type CommandHandler func(ctx context.Context, s *ServiceCaller) (any, error)

func (we *WorkflowEngine) RouteCommand(commandName string, handler CommandHandler, timeoutPolicy *service.UniversalSagaTimeoutPolicySchema) {
	we.commandMap[commandName] = newCommand(
		timeoutPolicy,
		we.nc,
		handler,
	)
}
