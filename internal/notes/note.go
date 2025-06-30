package notes

type Note struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Project   string `json:"project"`
	Tags      []string `json:"tags"`
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

