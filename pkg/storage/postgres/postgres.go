package postgres

import (
	"GoNews/pkg/storage"
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db *pgxpool.Pool
}

// Конструктор для БД
func New(connstr string) (*Storage, error) {
	db, err := pgxpool.New(context.Background(), connstr)
	if err != nil {
		return nil, err
	}

	p := Storage{
		db: db,
	}

	return &p, nil
}

// Функция для получения n новостей.
func (s *Storage) GetPosts(n int) ([]storage.Post, error) {
	if n == 0 {
		n = 10
	}
	rows, err := s.db.Query(context.Background(), `
		SELECT
		id,
		title,
		content,
		pub_time,
		link
		FROM posts
		ORDER BY pub_time DESC
		LIMIT $1;
	`,
		n,
	)

	if err != nil {
		return nil, err
	}

	var posts []storage.Post

	for rows.Next() {
		var p storage.Post
		err = rows.Scan(
			&p.ID,
			&p.Title,
			&p.Content,
			&p.PubTime,
			&p.Link,
		)
		if err != nil {
			return nil, err
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// Функция для добавления новостей в БД.
func (s *Storage) AddPosts(p []storage.Post) error {
	for _, post := range p {
		_, err := s.db.Exec(context.Background(), `
		INSERT INTO posts (title, content, pub_time, link) VALUES ($1, $2, $3, $4);
	`,
			post.Title,
			post.Content,
			post.PubTime,
			post.Link,
		)
		if err != nil {
			return err
		}
	}
	return nil
}
