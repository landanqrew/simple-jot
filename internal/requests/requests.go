package requests

import (
	"bytes"
	"io"
	"net/http"
)

func MakeRequest(url string, payload []byte) ([]byte, error) {
	res, err := http.Post(url, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return nil, err
	}	
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}