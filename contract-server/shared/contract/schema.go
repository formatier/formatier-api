package contract

type ContractSchema struct {
	ServiceName string
	Version     string
	Commands    map[string]CommandContractSchema
}

type CommandContractSchema struct {
	/*
		The data that the command wants.

		map-key: key-name
		map-value: key-type
	*/
	InputFields map[string]ContractValueSchema `json:"input_fields"`

	/*
		The data that the command will send to orch. server.

		map-key: key-name
		map-value: key-type
	*/
	OutputFields map[string]ContractValueSchema `json:"output_fields"`

	SaveMode string `json:"save_mode"`
}

type ContractValueSchema struct {
	/*
		The data that the command will send to orch. server.

		map-key: key-name
		map-value: key-type
	*/
	Type        string `json:"type"`
	Description string `json:"description"`

	Map   map[string]string  `json:"map"`
	Array ContractArrayValue `json:"arr"`
}

type ContractArrayValue struct {
	ItemType string            `json:"item_type"`
	Map      map[string]string `json:"map"`
}

const (
	DATA_TYPE_STR   = "str"
	DATA_TYPE_INT   = "int"
	DATA_TYPE_FLOAT = "float"
	DATA_TYPE_BOOL  = "bool"
	DATA_TYPE_MAP   = "map"
	DATA_TYPE_ARR   = "arr"
)

const (
	SAVE_MODE_REPLACE = "rep"
	SAVE_MODE_SAFE    = "safe"
)
