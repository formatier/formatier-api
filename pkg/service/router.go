package service

import (
	"context"
)

func NewRouter(ctx context.Context) *Router {
	return &Router{
		commandMap: make(map[string]*Command),
		ctx:        ctx,
	}
}

type Router struct {
	commandMap map[string]*Command
	ctx        context.Context
}

type Command struct {
	handlers []Handler
}

type Handler func(*Context) error

func (r *Router) Route(command string, handlers ...Handler) {
	if _, ok := r.commandMap[command]; ok {
		panic("command aready exits")
	}

	r.commandMap[command] = &Command{
		handlers: handlers,
	}
}
