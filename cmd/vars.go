package cmd

// Shared variables used across multiple commands
var (
	// noteContent is used by both create and edit commands
	noteContent string
	// setNote is used by the create command
	setNote bool
	// appendContent is used by the edit command
	appendContent string
)
