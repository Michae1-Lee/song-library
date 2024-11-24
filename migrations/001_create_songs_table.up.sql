CREATE TABLE IF NOT EXISTS songs (
    id SERIAL PRIMARY KEY,               -- Уникальный идентификатор
    group_name TEXT NOT NULL,            -- Название группы
    song_name TEXT NOT NULL,             -- Название песни
    release_date TEXT NOT NULL,          -- Дата релиза
    text TEXT NOT NULL,                  -- Текст песни
    link TEXT NOT NULL,                  -- Ссылка на дополнительную информацию
    created_at TIMESTAMP DEFAULT now(),  -- Дата создания записи
    updated_at TIMESTAMP DEFAULT now()   -- Дата последнего обновления записи
);
