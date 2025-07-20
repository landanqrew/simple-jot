package ai

import (
	"encoding/json"
	"log"

	"github.com/landanqrew/simple-jot/internal/requests"
)


func SemanticSearch[T any](data []T, query string, apiKey string) ([]T, error) {
	// marshall json
	jsonData, err := json.Marshal(data)
	if err != nil {
		return []T{}, err
	}

	// call gemini api

	// 1. Prepare the request payload
	payload := map[string]any{
		"contents": []map[string]any{
			{
				"parts": []map[string]any{
					{
						"text": "Given the following data: " + string(jsonData) + ", find the items that are most relevant to the query: " + query,
					},
				},
			},
		},
	}
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return []T{}, err
	}

	// 2. Construct the API endpoint URL
	apiEndpoint := "https://generative-ai-api.googleapis.com/v1beta/models/gemini-1.0-pro:generateContent?key=" + apiKey

	// 3. Create and execute the HTTP request
	//log.Println("apiEndpoint: ", apiEndpoint)
	resp, err := requests.MakeRequest(apiEndpoint, payloadBytes)
	if err != nil {
		return []T{}, err
	}
	log.Println("resp: ", string(resp))


	// unmarshall json
	var items []T
	err = json.Unmarshal(resp, &items)
	if err != nil {
		return []T{}, err
	}
	// return filtered data
	return items, nil

	

}