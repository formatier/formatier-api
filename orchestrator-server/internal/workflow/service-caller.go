package workflow

import (
	"context"
	"encoding/json"
	"fmt"
	"formatier-api/pkg/service"
	"formatier-api/shared/future"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type ServiceCaller struct {
	mx sync.Mutex

	timeoutPolicy *service.UniversalSagaTimeoutPolicySchema
	nc            *nats.Conn
	ctxCancelFunc context.CancelCauseFunc

	processedEvents []sagaWorkflowSchema
	sagaMetadata    *service.UniversalSagaMetadataSchema
}

func (sc *ServiceCaller) Call(
	ctx context.Context,
	serviceName string,
	eventName string,
	eventTriggerPayload any,
	cancelable bool,
) (*service.UniversalSagaReciverSchema, error) {
	triggerEvent := service.UniversalSagaTriggerEventSchema{
		TimeoutPolicy: *sc.timeoutPolicy,
		Metadata:      *sc.sagaMetadata,
		Payload:       eventTriggerPayload,
	}
	parsedTriggerEvent, err := json.Marshal(triggerEvent)
	if err != nil {
		return nil, err
	}

	requestCtx, requestCtxCancelFunc := context.WithTimeout(ctx, time.Duration(sc.timeoutPolicy.Duration))
	defer requestCtxCancelFunc()
	reciverEventMsg, err := sc.nc.RequestWithContext(
		requestCtx,
		fmt.Sprintf("%s.%s", serviceName, eventName),
		parsedTriggerEvent,
	)
	if err != nil {
		return nil, err
	}

	reciverEvent := &service.UniversalSagaReciverSchema{}
	err = json.Unmarshal(reciverEventMsg.Data, reciverEvent)
	if err != nil {
		return nil, err
	}

	sc.commit(serviceName, eventName, cancelable)

	return reciverEvent, nil
}

func (sc *ServiceCaller) FutureCall(
	ctx context.Context,
	serviceName string,
	eventName string,
	eventTriggerPayload any,
	cancelable bool,
) *future.ErrorFuture[*service.UniversalSagaReciverSchema] {
	triggerEventFuture := future.NewErrorFuture[*service.UniversalSagaReciverSchema]()

	go func() {
		reciverEvent, err := sc.Call(ctx, serviceName, eventName, eventTriggerPayload, cancelable)
		if err != nil {
			triggerEventFuture.SendError(err)
			return
		}
		triggerEventFuture.SendValue(reciverEvent)
	}()

	return triggerEventFuture
}

func (sc *ServiceCaller) commit(serviceName string, eventName string, cancelable bool) {
	sc.mx.Lock()
	defer sc.mx.Unlock()
	sc.processedEvents = append(sc.processedEvents, sagaWorkflowSchema{
		ServiceName: serviceName,
		EventName:   eventName,
		Status:      WORKFLOW_STATUS_SUCCESS,
	})
}

func (sc *ServiceCaller) rollback() {
	sc.mx.Lock()
	defer sc.mx.Unlock()
}
