package service

import (
	"context"
	"contract-server/shared/contract"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func NewRouter(ctx context.Context) *Router {
	return &Router{
		commandMap: make(map[string]*RouteCommand),
		ctx:        ctx,
	}
}

type Router struct {
	commandMap map[string]*RouteCommand
	ctx        context.Context
}

type RouteCommand struct {
	handlers         []Handler
	commandContracts *contract.CommandContractSchema
}

type Handler func(*Context) error

func (r *Router) Use(command string, commandContracts *contract.CommandContractSchema, handlers ...Handler) {
	if _, ok := r.commandMap[command]; ok {
		panic("command aready exits")
	}

	r.commandMap[command] = &RouteCommand{
		handlers:         handlers,
		commandContracts: commandContracts,
	}
}

func (r *Router) ExportCommandContracts(serviceName string, version string) map[string]contract.CommandContractSchema {
	contracts := make(map[string]contract.CommandContractSchema)
	for commandKey, command := range r.commandMap {
		contracts[commandKey] = *command.commandContracts
	}

	return contracts
}

func (r *Router) Listen(deliveryChan <-chan *amqp.Delivery) {

	for delivery := range deliveryChan {
		go r.route(delivery, r.ctx)
	}

}

func (r *Router) route(delivery *amqp.Delivery, ctx context.Context) {
	routeCtx /*ctxCancelFunc*/, _ := context.WithCancel(ctx)

	rawCommandPayload := delivery.Body

	commandPayload := &CommandSchema{}
	json.Unmarshal(rawCommandPayload, commandPayload)

	cancelChan := make(chan struct{})
	doneChan := make(chan struct{})

	fmt.Println(delivery.RoutingKey)

	handlerCtx := &Context{
		reqData:    commandPayload.Data,
		ctx:        routeCtx,
		cancelChan: cancelChan,
		doneChan:   doneChan,
	}

	handlerCtx.Done()
}

func (r *Router) runCommand(ctx *Context) {
}
