package domain

type Song struct {
	ID          int    `json:"id"`           // Уникальный идентификатор песни
	Group       string `json:"group"`        // Название группы
	Song        string `json:"song"`         // Название песни
	ReleaseDate string `json:"release_date"` // Дата релиза песни
	Text        string `json:"text"`         // Текст песни
	Link        string `json:"link"`         // Ссылка на дополнительную информацию
}
