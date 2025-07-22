package workflow

func NewWorkflowEngine() *WorkflowEngine {
	return &WorkflowEngine{
		commandMap: make(map[string]*Command),
	}
}

type WorkflowEngine struct {
	commandMap map[string]*Command
}

type CommandHandler func(s *ServiceCaller, triggerEventPayload any)

func (we *WorkflowEngine) RouteCommand(command string, handler CommandHandler) {
	we.commandMap[command] = &Command{
		commandHandler: handler,
	}
}
