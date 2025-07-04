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
	// TODO: Implement Gemini API call here
	// This is a placeholder
	// log.Println("jsonData: ", string(jsonData))
	// log.Println("query: ", query)
	// log.Println("apiKey: ", apiKey)

	// Implement Gemini API call here
	// 1. Prepare the request payload
	payload := map[string]interface{}{
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]interface{}{
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
	// Assuming the Gemini API endpoint is something like this:
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