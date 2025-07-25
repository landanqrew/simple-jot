#!/bin/bash

# set api key
export GEMINI_API_KEY=$(cat ./IGNORE/gemini-api-key.txt )

# Set nullglob option to make unmatched globs expand to nothing
shopt -s nullglob

# iterate through sub directories and run tests if *_test.go exists
for dir in $(find . -type d); do
    cd "$dir" || { echo "Failed to change directory to $dir"; exit 1; } # Added error handling for cd
    test_files=(*_test.go) # Create an array of matching files
    if [ ${#test_files[@]} -gt 0 ]; then # Check if the array is not empty
        echo "Running tests in $dir"
        go test ./... # Use ./... to run tests in current directory and sub-packages
    fi
    cd - > /dev/null # Go back to the previous directory quietly
done
shopt -u nullglob # Unset nullglob when done