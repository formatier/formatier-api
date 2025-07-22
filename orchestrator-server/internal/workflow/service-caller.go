package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"formatier-api/pkg/service"
	"formatier-api/shared/saga"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type ServiceCaller struct {
	mx            sync.Mutex
	TimeoutPolicy *service.UniversalSagaTimeoutPolicySchema
	nc            *nats.Conn
	errorChan     chan error

	SagaWorkflow []saga.SagaWorkflowSchema `json:"workflow"`

	isRollbackRequire bool
}

func (sc *ServiceCaller) Call(
	ctx context.Context,
	serviceName string,
	eventName string,
	eventTriggerPayload any,
) *service.UniversalSagaReciverSchema {
	sc.mx.Lock()
	if sc.isRollbackRequire {
		sc.mx.Unlock()
		return nil
	}
	sc.mx.Unlock()

	triggerEvent := service.UniversalSagaTriggerEventSchema{
		TimeoutPolicy: *sc.TimeoutPolicy,
		Metadata:      service.UniversalSagaSagaMetadataSchema{},
		Payload:       eventTriggerPayload,
	}
	parsedTriggerEvent, err := json.Marshal(triggerEvent)
	if err != nil {
		sc.errorChan <- err
		return nil
	}

	requestCtx, requestCtxCancelFunc := context.WithTimeout(ctx, time.Duration(sc.TimeoutPolicy.Duration))
	defer requestCtxCancelFunc()
	reciverEventMsg, err := sc.nc.RequestWithContext(
		requestCtx,
		fmt.Sprintf("%s.%s", serviceName, eventName),
		parsedTriggerEvent,
	)
	if err != nil {
		sc.errorChan <- err
		return nil
	}

	reciverEvent := &service.UniversalSagaReciverSchema{}
	err = json.Unmarshal(reciverEventMsg.Data, reciverEvent)
	if err != nil {
		sc.errorChan <- err
		return nil
	}

	return reciverEvent
}

func (sc *ServiceCaller) ChanCall(
	ctx context.Context,
	serviceName string,
	eventName string,
	eventTriggerPayload any,
	eventReciverPayloadChan chan<- any,
) {
	go func() {
		reciverEvent := sc.Call(ctx, serviceName, eventName, eventTriggerPayload)
		eventReciverPayloadChan <- reciverEvent
	}()
}
