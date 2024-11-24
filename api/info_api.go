package api

import (
	"encoding/json"
	"net/http"
	"song-library/service"
)

type InfoController struct {
	service *service.SongService
}

func NewInfoController(service *service.SongService) *InfoController {
	return &InfoController{service: service}
}

// InfoHandler обрабатывает запросы к API /info.
//
//	@Summary		Получить информацию о песне
//	@Description	Получение деталей о песне по имени группы и названию песни.
//	@Tags			Info
//	@Param			group	query	string	true	"Название группы"
//	@Param			song	query	string	true	"Название песни"
//	@Success		200		{object}	domain.SongDetail
//	@Failure		400		{string}	string	"Параметры обязательны"
//	@Failure		500		{string}	string	"Ошибка получения данных"
//	@Router			/info [get]
func (c *InfoController) InfoHandler(w http.ResponseWriter, r *http.Request) {
	group := r.PathValue("group")
	song := r.PathValue("song")

	// Проверка обязательных параметров
	if group == "" || song == "" {
		http.Error(w, "Параметры 'group' и 'song' обязательны", http.StatusBadRequest)
		return
	}

	response := map[string]string{
		"releaseDate": "16.07.2006",
		"text": `Ooh baby, don't you know I suffer?
Ooh baby, can you hear me moan?
You caught me under false pretenses
How long before you let me go?

Ooh
You set my soul alight
Ooh
You set my soul alight`,
		"link": "https://www.youtube.com/watch?v=Xsp3_a-PMTw",
	}

	// Возврат данных в формате JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}
