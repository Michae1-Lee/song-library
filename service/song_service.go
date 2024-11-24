package service

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"song-library/domain"
	"song-library/repository"
)

type SongService struct {
	repo       *repository.SongRepository
	log        *log.Logger
	apiBaseURL string
}

func NewSongService(repo *repository.SongRepository, logger *log.Logger, apiBaseURL string) *SongService {
	return &SongService{repo: repo, log: logger, apiBaseURL: apiBaseURL}
}

// GetLibrary получает список песен с учетом пагинации.
func (service *SongService) GetLibrary(page, limit int) ([]domain.Song, error) {
	if page <= 0 || limit <= 0 {
		err := fmt.Errorf("некорректные параметры пагинации: page=%d, limit=%d", page, limit)
		service.log.Printf("ошибка в GetLibrary: %v", err)
		return nil, err
	}

	offset := calculateOffset(page, limit)
	songs, err := service.repo.GetSongs(offset, limit)
	if err != nil {
		service.log.Printf("ошибка получения песен в GetLibrary: offset=%d, limit=%d, error=%v", offset, limit, err)
		return nil, fmt.Errorf("ошибка получения библиотеки песен: %w", err)
	}

	service.log.Printf("успешно выполнен GetLibrary: page=%d, limit=%d", page, limit)
	return songs, nil
}

// AddSong добавляет новую песню с запросом к внешнему API для получения деталей.
func (service *SongService) AddSong(song domain.Song) error {
	if song.Group == "" || song.Song == "" {
		err := fmt.Errorf("группа и название песни не могут быть пустыми")
		service.log.Printf("ошибка в AddSong: %v", err)
		return err
	}

	// Получение данных из внешнего API
	details, err := service.fetchSongDetails(song.Group, song.Song)
	if err != nil {
		service.log.Printf("ошибка получения данных из внешнего API: group=%s, song=%s, error=%v", song.Group, song.Song, err)
		return fmt.Errorf("ошибка получения данных из внешнего API: %w", err)
	}

	// Обновляем структуру песни деталями из API
	song.ReleaseDate = details.ReleaseDate
	song.Text = details.Text
	song.Link = details.Link

	// Добавляем песню в базу данных
	if err := service.repo.AddSong(song); err != nil {
		service.log.Printf("ошибка добавления песни в базу данных: group=%s, song=%s, error=%v", song.Group, song.Song, err)
		return fmt.Errorf("ошибка добавления песни в базу данных: %w", err)
	}

	service.log.Printf("песня успешно добавлена: group=%s, song=%s", song.Group, song.Song)
	return nil
}

// UpdateSong обновляет существующую песню.
func (service *SongService) UpdateSong(song domain.Song) error {
	if song.ID <= 0 {
		err := fmt.Errorf("некорректный ID песни: %d", song.ID)
		service.log.Printf("ошибка в UpdateSong: %v", err)
		return err
	}

	if err := service.repo.UpdateSong(song); err != nil {
		service.log.Printf("ошибка обновления песни: id=%d, error=%v", song.ID, err)
		return fmt.Errorf("ошибка обновления песни: %w", err)
	}

	service.log.Printf("песня успешно обновлена: id=%d", song.ID)
	return nil
}

// DeleteSong удаляет песню по ID.
func (service *SongService) DeleteSong(id int) error {
	if id <= 0 {
		err := fmt.Errorf("некорректный ID песни: %d", id)
		service.log.Printf("ошибка в DeleteSong: %v", err)
		return err
	}

	if err := service.repo.DeleteSong(id); err != nil {
		service.log.Printf("ошибка удаления песни: id=%d, error=%v", id, err)
		return fmt.Errorf("ошибка удаления песни: %w", err)
	}

	service.log.Printf("песня успешно удалена: id=%d", id)
	return nil
}

// GetSongByID получает песню по ID.
func (service *SongService) GetSongByID(id int) (*domain.Song, error) {
	if id <= 0 {
		err := fmt.Errorf("некорректный ID песни: %d", id)
		service.log.Printf("ошибка в GetSongByID: %v", err)
		return nil, err
	}

	song, err := service.repo.GetSongByID(id)
	if err != nil {
		service.log.Printf("ошибка получения песни: id=%d, error=%v", id, err)
		return nil, fmt.Errorf("ошибка получения песни: %w", err)
	}

	service.log.Printf("песня успешно получена: id=%d", id)
	return song, nil
}

// fetchSongDetails делает запрос к внешнему API и возвращает детали песни.
func (service *SongService) fetchSongDetails(group, song string) (*domain.SongDetail, error) {
	url := fmt.Sprintf("%s/info?group=%s&song=%s", service.apiBaseURL, group, song)

	resp, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса к API: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ошибка API: получен статус %d", resp.StatusCode)
	}

	var details domain.SongDetail
	if err := json.NewDecoder(resp.Body).Decode(&details); err != nil {
		return nil, fmt.Errorf("ошибка декодирования ответа API: %w", err)
	}

	return &details, nil
}
func (service *SongService) GetSongDetails(group, song string) (*domain.SongDetail, error) {
	if group == "" || song == "" {
		return nil, fmt.Errorf("параметры 'group' и 'song' обязательны")
	}

	// Используем fetchSongDetails для запроса к внешнему API
	details, err := service.fetchSongDetails(group, song)

	if err != nil {
		service.log.Printf("ошибка получения данных о песне из API: group=%s, song=%s, error=%v", group, song, err)
		return nil, err
	}

	service.log.Printf("успешно получены данные о песне: group=%s, song=%s", group, song)
	return details, nil
}

// calculateOffset вычисляет смещение для пагинации.
func calculateOffset(page, limit int) int {
	return (page - 1) * limit
}
