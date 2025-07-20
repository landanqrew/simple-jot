package notes

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
)

func GenerateTestNotes(numNotes int, filePath string) error {
	var notes []Note
	for i := 0; i < numNotes; i++ {
		id := uuid.New().String()
		title := fmt.Sprintf("Note Title %d", i)
		content := fmt.Sprintf("This is the content for note %d. It talks about various topics.", i)
		tags := []string{"tag1", "tag2"}
		createdAt := time.Now().Add(time.Duration(-i) * time.Hour).Format(time.DateTime)
		updatedAt := time.Now().Format(time.DateTime)

		if i == 5 {
			title = "Go Programming Best Practices"
			content = "Deep dive into Go concurrency patterns and error handling."
			tags = []string{"go", "concurrency", "error_handling"}
		}
		if i == 10 {
			title = "Python Machine Learning Libraries"
			content = "Exploring scikit-learn and TensorFlow for machine learning in Python."
			tags = []string{"python", "ml", "tensorflow"}
		}

		notes = append(notes, Note{
			ID:        id,
			Title:     title,
			Tags:      tags,
			Content:   content,
			CreatedAt: createdAt,
			UpdatedAt: updatedAt,
		})
	}

	jsonData, err := json.MarshalIndent(notes, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal notes: %w", err)
	}

	err = os.WriteFile(filePath, jsonData, 0644)
	if err != nil {
		return fmt.Errorf("failed to write notes to file: %w", err)
	}

	return nil
}


