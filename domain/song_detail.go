package domain

type SongDetail struct {
	ReleaseDate string `json:"releaseDate"` // Дата релиза песни
	Text        string `json:"text"`        // Текст песни
	Link        string `json:"link"`        // Ссылка на дополнительную информацию о песне
}
