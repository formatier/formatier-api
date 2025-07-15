package validator

import (
	"contract-server/shared/contract"
	"errors"
	"fmt"
)

func Validate(data map[string]any, schema map[string]contract.ContractValueSchema) error {
	for fieldKey, valueSchema := range schema {
		data, ok := data[fieldKey]
		if !ok {
			return fmt.Errorf("key %s does not exist", fieldKey)
		}

		switch valueSchema.Type {
		case contract.DATA_TYPE_MAP:
			err := validateMap(data, valueSchema.Map)
			if err != nil {
				return fmt.Errorf("key %s: %w", fieldKey, err)
			}
		case contract.DATA_TYPE_ARR:
			parsedData, ok := data.([]any)
			if !ok {
				return fmt.Errorf("key %s: data is not array", fieldKey)
			}

			for idx, data := range parsedData {
				if valueSchema.Array.ItemType == contract.DATA_TYPE_MAP {
					err := validateMap(data, valueSchema.Array.Map)
					if err != nil {
						return fmt.Errorf("data at %d in key %s: %w", idx, fieldKey, err)
					}
				} else {
					err := validateStandardType(valueSchema.Array.ItemType, data)
					if err != nil {
						return fmt.Errorf("data at %d in key %s: %w", idx, fieldKey, err)
					}
				}

			}
		default:
			err := validateStandardType(valueSchema.Type, data)
			if err != nil {
				return fmt.Errorf("key %s: %w", fieldKey, err)
			}
		}
	}

	return nil
}

func validateMap(data any, schema map[string]string) error {
	parsedData, ok := data.(map[string]any)
	if !ok {
		return errors.New("data is not map")
	}

	for schemaKey, schemaType := range schema {
		data, ok := parsedData[schemaKey]
		if !ok {
			return fmt.Errorf("key %s: data does not exist", schemaKey)
		}

		err := validateStandardType(schemaType, data)
		if err != nil {
			return fmt.Errorf("key %s: %w", schemaKey, err)
		}
	}

	return nil
}

func validateStandardType(typeName string, data any) error {
	switch typeName {
	case contract.DATA_TYPE_STR:
		_, ok := data.(string)
		if !ok {
			return errors.New("data is not string")
		}
	case contract.DATA_TYPE_INT:
		_, ok := data.(int)
		if !ok {
			return errors.New("data is not int")
		}
	case contract.DATA_TYPE_FLOAT:
		_, ok := data.(float32)
		if !ok {
			return errors.New("data is not float")
		}
	case contract.DATA_TYPE_BOOL:
		_, ok := data.(bool)
		if !ok {
			return errors.New("data is not boolean")
		}
	default:
		return fmt.Errorf("type %s is not a valid type", typeName)
	}

	return nil
}
