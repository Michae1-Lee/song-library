package controller

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"song-library/domain"
	"song-library/service"
)

// SongController представляет контроллер для работы с песнями.
type SongController struct {
	service *service.SongService
}

// NewSongController создает новый SongController.
func NewSongController(service *service.SongService) *SongController {
	return &SongController{service: service}
}

// GetLibraryHandler получает список песен с фильтрацией и пагинацией.
//
//	@Summary		Получить библиотеку песен
//	@Description	Получение списка песен с фильтрацией по группе, названию и дате релиза.
//	@Tags			Songs
//	@Param			page			query		int		false	"Номер страницы"					default(1)
//	@Param			limit			query		int		false	"Количество элементов на странице"	default(10)
//	@Param			group			query		string	false	"Фильтр по группе"
//	@Param			song			query		string	false	"Фильтр по названию песни"
//	@Param			release_date	query		string	false	"Фильтр по дате релиза"
//	@Success		200				{array}		domain.Song
//	@Failure		500				{string}	string	"Ошибка получения библиотеки"
//	@Router			/library [get]
func (c *SongController) GetLibraryHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()

	page, err := strconv.Atoi(r.PathValue("page"))
	if err != nil || page < 1 {
		page = 1
	}
	limit, err := strconv.Atoi(r.PathValue("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	group := query.Get("group")
	songName := query.Get("song")
	releaseDate := query.Get("release_date")

	songs, err := c.service.GetLibrary(page, limit)
	if err != nil {
		http.Error(w, "Ошибка получения библиотеки: "+err.Error(), http.StatusInternalServerError)
		return
	}

	filtered := []domain.Song{}
	for _, song := range songs {
		if group != "" && !strings.Contains(song.Group, group) {
			continue
		}
		if songName != "" && !strings.Contains(song.Song, songName) {
			continue
		}
		if releaseDate != "" && song.ReleaseDate != releaseDate {
			continue
		}
		filtered = append(filtered, song)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(filtered)
}

// GetSongTextHandler получает текст песни по ID.
//
//	@Summary		Получить текст песни
//	@Description	Получение текста песни по ID с возможностью пагинации по куплетам.
//	@Tags			Songs
//	@Param			id		path		int	true	"ID песни"
//	@Param			page	query		int	false	"Номер страницы (по куплетам)"	default(1)
//	@Success		200		{object}	map[string]interface{}
//	@Failure		400		{string}	string	"Неверный ID песни"
//	@Failure		404		{string}	string	"Песня или куплеты не найдены"
//	@Router			/song/{id}/text [get]
func (c *SongController) GetSongTextHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id") // Динамический сегмент {id}
	songID, err := strconv.Atoi(id)
	if err != nil || songID < 1 {
		http.Error(w, "Неверный ID песни", http.StatusBadRequest)
		return
	}

	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit := 1
	song, err := c.service.GetSongByID(songID)
	if err != nil {
		http.Error(w, "Ошибка получения песни: "+err.Error(), http.StatusNotFound)
		return
	}

	verses := strings.Split(song.Text, "\n\n")
	start := (page - 1) * limit
	if start >= len(verses) {
		http.Error(w, "Куплеты не найдены", http.StatusNotFound)
		return
	}

	end := start + limit
	if end > len(verses) {
		end = len(verses)
	}

	response := map[string]interface{}{
		"id":     song.ID,
		"group":  song.Group,
		"song":   song.Song,
		"verses": verses[start:end],
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// DeleteSongHandler удаляет песню по ID.
//
//	@Summary		Удалить песню
//	@Description	Удаление песни из библиотеки по ID.
//	@Tags			Songs
//	@Param			id	path	int	true	"ID песни"
//	@Success		204	"Песня удалена"
//	@Failure		400	{string}	string	"Неверный ID песни"
//	@Failure		500	{string}	string	"Ошибка удаления песни"
//	@Router			/song/{id} [delete]
func (c *SongController) DeleteSongHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	songID, err := strconv.Atoi(id)
	if err != nil || songID < 1 {
		http.Error(w, "Неверный ID песни", http.StatusBadRequest)
		return
	}

	if err := c.service.DeleteSong(songID); err != nil {
		http.Error(w, "Ошибка удаления песни: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateSongHandler обновляет данные песни по ID.
//
//	@Summary		Обновить данные песни
//	@Description	Обновление информации о песне по ID.
//	@Tags			Songs
//	@Param			id		path	int			true	"ID песни"
//	@Param			song	body	domain.Song	true	"Данные песни"
//	@Success		200		"Песня обновлена"
//	@Failure		400		{string}	string	"Ошибка декодирования данных или неверный ID"
//	@Failure		500		{string}	string	"Ошибка обновления песни"
//	@Router			/song/{id} [put]
func (c *SongController) UpdateSongHandler(w http.ResponseWriter, r *http.Request) {
	id := r.PathValue("id")
	songID, err := strconv.Atoi(id)
	if err != nil || songID < 1 {
		http.Error(w, "Неверный ID песни", http.StatusBadRequest)
		return
	}

	var song domain.Song
	if err := json.NewDecoder(r.Body).Decode(&song); err != nil {
		http.Error(w, "Ошибка декодирования данных песни: "+err.Error(), http.StatusBadRequest)
		return
	}
	song.ID = songID

	if err := c.service.UpdateSong(song); err != nil {
		http.Error(w, "Ошибка обновления песни: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// AddSongHandler добавляет новую песню в библиотеку.
//
//	@Summary		Добавить песню
//	@Description	Добавление новой песни в библиотеку.
//	@Tags			Songs
//	@Param			song	body	domain.SongCreateRequest	true	"Данные для создания песни"
//	@Success		201		"Песня добавлена"
//	@Failure		400		{string}	string	"Ошибка декодирования данных песни"
//	@Failure		500		{string}	string	"Ошибка добавления песни"
//	@Router			/song [post]
func (c *SongController) AddSongHandler(w http.ResponseWriter, r *http.Request) {
	var request domain.SongCreateRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Ошибка декодирования данных песни: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Проверка обязательных полей
	if request.Group == "" || request.Song == "" {
		http.Error(w, "Поля 'group' и 'song' обязательны", http.StatusBadRequest)
		return
	}

	// Преобразование данных в доменную модель
	newSong := domain.Song{
		Group: request.Group,
		Song:  request.Song,
	}

	// Добавление песни через сервис
	if err := c.service.AddSong(newSong); err != nil {
		http.Error(w, "Ошибка добавления песни: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}
