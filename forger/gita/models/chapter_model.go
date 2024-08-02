package models

type Chapter struct {
	ChapterNumber   int           `json:"chapter_number"`
	VersesCount     int           `json:"verses_count"`
	Name            string        `json:"name"`
	Translation     string        `json:"translation"`
	Transliteration string        `json:"transliteration"`
	Meaning         []ChapterText `json:"meaning"`
	Summary         []ChapterText `json:"summary"`
}

type ChapterText struct {
	Language string `json:"language"`
	Text     string `json:"text"`
}
