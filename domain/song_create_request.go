package domain

// SongCreateRequest представляет данные, которые клиент отправляет для добавления новой песни.
type SongCreateRequest struct {
	Group string `json:"group"` // Название группы (обязательно)
	Song  string `json:"song"`  // Название песни (обязательно)
}
