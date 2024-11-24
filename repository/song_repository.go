package repository

import (
	"database/sql"
	"errors"
	"log"
	"song-library/domain"
)

type SongRepository struct {
	db  *sql.DB
	log *log.Logger
}

func NewSongRepository(db *sql.DB, logger *log.Logger) *SongRepository {
	return &SongRepository{db: db, log: logger}
}

func (repo *SongRepository) GetSongs(offset, limit int) ([]domain.Song, error) {
	rows, err := repo.db.Query("SELECT id, group_name, song_name, release_date, text, link FROM songs LIMIT $1 OFFSET $2", limit, offset)
	if err != nil {
		repo.log.Printf("ошибка выполнения GetSongs: %v", err)
		return nil, err
	}
	defer func() {
		if closeErr := rows.Close(); closeErr != nil {
			repo.log.Printf("ошибка закрытия rows в GetSongs: %v", closeErr)
		}
	}()

	var songs []domain.Song
	for rows.Next() {
		var song domain.Song
		if err := rows.Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link); err != nil {
			repo.log.Printf("ошибка сканирования строки в GetSongs: %v", err)
			return nil, err
		}
		songs = append(songs, song)
	}

	if err := rows.Err(); err != nil {
		repo.log.Printf("ошибка итерации строк в GetSongs: %v", err)
		return nil, err
	}

	repo.log.Printf("успешно выполнен GetSongs (offset=%d, limit=%d)", offset, limit)
	return songs, nil
}

func (repo *SongRepository) AddSong(song domain.Song) error {
	_, err := repo.db.Exec(
		"INSERT INTO songs (group_name, song_name, release_date, text, link) VALUES ($1, $2, $3, $4, $5)",
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link,
	)
	if err != nil {
		repo.log.Printf("ошибка добавления песни: group=%s, song=%s, error=%v", song.Group, song.Song, err)
		return err
	}

	repo.log.Printf("песня успешно добавлена: group=%s, song=%s", song.Group, song.Song)
	return nil
}

func (repo *SongRepository) UpdateSong(song domain.Song) error {
	res, err := repo.db.Exec(
		"UPDATE songs SET group_name = $1, song_name = $2, release_date = $3, text = $4, link = $5 WHERE id = $6",
		song.Group, song.Song, song.ReleaseDate, song.Text, song.Link, song.ID,
	)
	if err != nil {
		repo.log.Printf("ошибка обновления песни: id=%d, error=%v", song.ID, err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		repo.log.Printf("ошибка получения количества затронутых строк в UpdateSong: %v", err)
		return err
	}

	if rowsAffected == 0 {
		repo.log.Printf("песня для обновления не найдена: id=%d", song.ID)
		return errors.New("песня не найдена")
	}

	repo.log.Printf("песня успешно обновлена: id=%d", song.ID)
	return nil
}

func (repo *SongRepository) DeleteSong(id int) error {
	res, err := repo.db.Exec("DELETE FROM songs WHERE id = $1", id)
	if err != nil {
		repo.log.Printf("ошибка удаления песни: id=%d, error=%v", id, err)
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		repo.log.Printf("ошибка получения количества затронутых строк в DeleteSong: %v", err)
		return err
	}

	if rowsAffected == 0 {
		repo.log.Printf("песня для удаления не найдена: id=%d", id)
		return errors.New("песня не найдена")
	}

	repo.log.Printf("песня успешно удалена: id=%d", id)
	return nil
}

func (repo *SongRepository) GetSongByID(id int) (*domain.Song, error) {
	var song domain.Song
	err := repo.db.QueryRow("SELECT id, group_name, song_name, release_date, text, link FROM songs WHERE id = $1", id).
		Scan(&song.ID, &song.Group, &song.Song, &song.ReleaseDate, &song.Text, &song.Link)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			repo.log.Printf("песня не найдена: id=%d", id)
			return nil, errors.New("песня не найдена")
		}
		repo.log.Printf("ошибка получения песни по ID: id=%d, error=%v", id, err)
		return nil, err
	}

	repo.log.Printf("песня успешно получена: id=%d", id)
	return &song, nil
}
