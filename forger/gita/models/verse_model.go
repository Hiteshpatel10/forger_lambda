package models

// Sloka represents the structure of the JSON data
type Verse struct {
	ID              string         `json:"_id"`
	Chapter         int            `json:"chapter"`
	Verse           int            `json:"verse"`
	Slok            string         `json:"slok"`
	Transliteration string         `json:"transliteration"`
	Comments        []VerseComment `json:"comments"`
}

type VerseLanguage struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}

// Comment represents a comment with its languages and author
type VerseComment struct {
	Languages []VerseLanguage `json:"languages"`
	Author    string          `json:"author"`
}
