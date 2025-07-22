package handler

import "orchestrator-server/internal/workflow"

func main() {
	we := workflow.NewWorkflowEngine()
	we.RouteCommand("GetUserById", func(sc *workflow.ServiceCaller, a any) {})
}
