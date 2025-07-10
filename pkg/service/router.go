package service

type Context struct{}

func (c *Context) Get() {

}

type Handler func(*Context) error

type Router struct {
	commandMap map[string][]*Handler
}

func (r *Router) Use(command string, handlers ...Handler) {
	if _, ok := r.commandMap[command]; ok {
		panic("command aready exits")
	}

	for _, handler := range handlers {
		r.commandMap[command] = append(r.commandMap[command], &handler)
	}
}

func (r *Router) Listen() {

}
