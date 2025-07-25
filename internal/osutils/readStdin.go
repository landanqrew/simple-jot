package osutils

import (
	"bufio"
	"os"
)

// Function that will determine if stdin is from a pipe or redirect and return the content
func ReadStdin() (string, error) {
	stat, _ := os.Stdin.Stat()
	if (stat.Mode() & os.ModeCharDevice) == 0 { // Check if stdin is from a pipe or redirect
		scanner := bufio.NewScanner(os.Stdin)
		var stdinContent []byte
		for scanner.Scan() {
			stdinContent = append(stdinContent, scanner.Bytes()...)
			stdinContent = append(stdinContent, '\n') // Add newline after each scanned line
		}
		if err := scanner.Err(); err != nil {
			return "", err
		}
		return string(stdinContent), nil
	}

	return "", nil
}