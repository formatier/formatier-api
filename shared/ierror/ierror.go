package ierror

const (
	ErrorSchema = "ierrschema"
)

type IActionSchema struct {
	WhereShouldErrorShow string `json:"where"` // just_do, popup, sonner
	Todo                 string `json:"todo"`  // go, reload, recall, show_error, rollback, wait, open

	Path             string   `json:"path"`   // if todo = go
	ThingsHaveToOpen []string `json:"things"` // if todo = open
}

type IErrorSchema struct {
	UseSchema    string          `json:"use_schema"`
	ErrorMessage string          `json:"error_message"`
	Actions      []IActionSchema `json:"transactions,omitempty"`
}
