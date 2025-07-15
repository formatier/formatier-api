package service

type CommandSchema struct {
	Metadata CommandMetadata `json:"metadata"`
	Data     map[string]any  `json:"data"`
}

type CommandMetadata struct {
	SagaId    string
	Timestamp uint
}
