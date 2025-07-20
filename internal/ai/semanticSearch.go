package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"

	"google.golang.org/genai"
)

type SearchResponse struct {
	PrimaryKey string  `json:"primary_key"`
	Score      float64 `json:"score"`
}

// generateSchemaFromStruct generates a Gemini schema from a Go struct type
func generateSchemaFromStruct[T any]() *genai.Schema {
	var zero T
	t := reflect.TypeOf(zero)

	properties := make(map[string]*genai.Schema)
	var propertyOrdering []string

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag == "" {
			jsonTag = field.Name
		}

		propertyOrdering = append(propertyOrdering, jsonTag)

		switch field.Type.Kind() {
		case reflect.String:
			properties[jsonTag] = &genai.Schema{Type: genai.TypeString}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			properties[jsonTag] = &genai.Schema{Type: genai.TypeInteger}
		case reflect.Float32, reflect.Float64:
			properties[jsonTag] = &genai.Schema{Type: genai.TypeNumber}
		case reflect.Slice:
			if field.Type.Elem().Kind() == reflect.String {
				properties[jsonTag] = &genai.Schema{
					Type:  genai.TypeArray,
					Items: &genai.Schema{Type: genai.TypeString},
				}
			}
		}
	}

	return &genai.Schema{
		Type:             genai.TypeObject,
		Properties:       properties,
		PropertyOrdering: propertyOrdering,
	}
}

func SemanticSearch[T any](data []T, naturalLanguageQuery string, apiKey string) ([]SearchResponse, error) {
    // marshall json
    jsonDataBytes, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("failed to marshal data: %w", err)
    }
    jsonData := string(jsonDataBytes)

    
	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %w", err)
	}
	sysPrompt := fmt.Sprintf(`
You are a helpful assistant that can answer questions and evaluate data in accordance with the provided schema.
You are given a list of data and a natural language query.
You need to answer the question based on the data.
The data is in the following format:
%s
The natural language query is:
%s
The response should be in the following format:
[{"primary_key": "string", "score": float64}]
The score is a float64 between 0 and 1, where 1 is the best match.
The primary key is the primary key of the data. Return no more than 10 results. Return the results with the strongest match.
`, jsonData, naturalLanguageQuery)

	config := &genai.GenerateContentConfig{
		ResponseMIMEType: "application/json",
		ResponseSchema: &genai.Schema{
			Type:  genai.TypeArray,
			Items: generateSchemaFromStruct[SearchResponse](),
		},
	}

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(sysPrompt),
		config,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to generate content: %w", err)
	}

	var searchResponses []SearchResponse
	err = json.Unmarshal([]byte(result.Text()), &searchResponses)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return searchResponses, nil
}