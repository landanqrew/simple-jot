#!/bin/bash

# set api key
export GEMINI_API_KEY=$(cat ./IGNORE/gemini-api-key.txt )

go test ./...