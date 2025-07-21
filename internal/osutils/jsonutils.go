package osutils

import "encoding/json"

func ToJsonString[T any](data T) (string, error) {
	jsonData, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}

func ReadJson[T any](data []byte) ([]T, error) {
	var output []T
	err := json.Unmarshal(data, &output)
	if err != nil {
		return nil, err
	}
	return output, nil
}
