package workflow

import (
	"context"
	"encoding/json"
	"errors"
	"formatier-api/pkg/service"
	"formatier-api/shared/future"
	"orchestrator-server/shared/orchestrator"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
)

func newCommand(
	timeoutPolicy *service.UniversalSagaTimeoutPolicySchema,
	nc *nats.Conn,
	commandHandler CommandHandler,
) *Command {
	return &Command{
		nc:             nc,
		timeoutPolicy:  timeoutPolicy,
		commandHandler: commandHandler,
	}
}

type Command struct {
	nc             *nats.Conn
	timeoutPolicy  *service.UniversalSagaTimeoutPolicySchema
	commandHandler CommandHandler
}

func (cmd *Command) Run(ctx context.Context, msg *nats.Msg) *future.ErrorFuture[any] {
	commandResponseFuture := future.NewErrorFuture[any]()

	commandRequest := &orchestrator.CommandRequestSchema{}
	err := json.Unmarshal(msg.Data, commandRequest)
	if err != nil {
		commandResponseFuture.SendError(err)
		return commandResponseFuture
	}

	sc := &ServiceCaller{
		timeoutPolicy: cmd.timeoutPolicy,
		nc:            cmd.nc,
		sagaMetadata: &service.UniversalSagaMetadataSchema{
			SagaID:    uuid.New().String(),
			Timestamp: time.Now().Unix(),
		},
		processedEvents: []sagaWorkflowSchema{},
	}

	handlerCtx, handlerCtxCancelFunc := context.WithTimeoutCause(
		ctx,
		time.Duration(commandRequest.TimeoutPolicy.Duration),
		errors.New("Timeout"),
	)
	defer handlerCtxCancelFunc()

	go func() {
		handlerResult, err := cmd.commandHandler(handlerCtx, sc)
		if err != nil {
			commandResponseFuture.SendError(err)
			return
		}
		commandResponseFuture.SendValue(handlerResult)
	}()

	return commandResponseFuture
}

func (cmd *Command) Rollback() {}
