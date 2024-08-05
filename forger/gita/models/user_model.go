package models

type Read struct {
	Chapter  int   `json:"chapter"`
	Verses   []int `json:"verses"`
	Progress int   `json:"progress"`
}

type User struct {
	Email       string  `json:"email"`
	DisplayName *string `json:"display_name,omitempty"`
	ProfileURL  *string `json:"profile_url,omitempty"`
	FCMToken    *string `json:"fcm_token,omitempty"`
	LastRead    string  `json:"last_read"`
	Reads       []Read  `json:"reads"`
	UpdatedAt   string  `json:"updated_at,omitempty"`
	CreatedAt   string  `json:"created_at,omitempty"`
}

type UpdateRead struct {
	ChapterNo int `json:"chapter_no"`
	VerseNo   int `json:"verse_no"`
}

var GitaChapters = []int{
	47,
	72,
	43,
	42,
	29,
	47,
	30,
	28,
	34,
	42,
	55,
	20,
	35,
	27,
	20,
	24,
	28,
	78,
}
