package domain

// DocumentChunk is a fraction of the text in the DB.
// It contains the text and the original source (ex: the title).
type DocumentChunk struct {
	Title   string
	Content string
	Score   float32
}
