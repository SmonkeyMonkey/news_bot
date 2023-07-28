package storage

import (
	"context"
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/samber/lo"
	"github.com/smonkeymonkey/news-bot/internal/model"
)

type ArticleStorage struct {
	db *sqlx.DB
}

type dbArticle struct {
	ID             int64          `db:"a_id"`
	SourcePriority int64          `db:"s_priority"`
	SourceID       int64          `db:"s_id"`
	Title          string         `db:"a_title"`
	Link           string         `db:"a_link"`
	Summary        sql.NullString `db:"a_summary"`
	PublishedAt    time.Time      `db:"a_published_at"`
	PostedAt       sql.NullTime   `db:"a_posted_at"`
	CreatedAt      time.Time      `db:"a_created_at"`
}

func NewArticleStorage(db *sqlx.DB) *ArticleStorage {
	return &ArticleStorage{db}
}

func (s *ArticleStorage) Store(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx,
		`INSERT INTO articles (source_id, title, link, published_at)
	    				VALUES ($1, $2, $3, $4, $5)
	    				ON CONFLICT DO NOTHING;`,
		article.SourceID,
		article.Title,
		article.Link,
		article.PostedAt); err != nil {
		return err
	}
	return nil
}

// articles that have not been published in the channel
func (s *ArticleStorage) AllNotPosted(ctx context.Context, since time.Time, limit uint64) ([]model.Article, error) {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	var articles []dbArticle
	if err := conn.SelectContext(ctx, &articles,
		`SELECT 
		a.id AS a_id, 
		s.priority AS s_priority,
		s.id AS s_id,
		a.title AS a_title,
		a.link AS a_link,
		a.published_at AS a_published_at,
		a.posted_at AS a_posted_at,
		a.created_at AS a_created_at
			FROM articles a JOIN sources s ON s.id = a.source_id
			WHERE a.posted_at IS NULL 
					AND a.published_at >= $1::timestamp
			ORDER BY a.created_at DESC, s_priority DESC LIMIT $2;`,
		since.UTC().Format(time.RFC3339),
		limit); err != nil {
		return nil, err
	}

	return lo.Map(articles, func(aritcle dbArticle, _ int) model.Article {
		return model.Article{
			SourceID:    aritcle.ID,
			Title:       aritcle.Title,
			Link:        aritcle.Link,
			PublishedAt: aritcle.PublishedAt,
			CreatedAt:   aritcle.CreatedAt,
		}
	}), nil
}

func (s *ArticleStorage) MarkAsPosted(ctx context.Context, article model.Article) error {
	conn, err := s.db.Connx(ctx)
	if err != nil {
		return err
	}
	defer conn.Close()

	if _, err := conn.ExecContext(ctx, `UPDATE articles SET posted_at = $1::timestamp WHERE id = $2;`,
		time.Now().UTC().Format(time.RFC3339),
		article.ID); err != nil {
		return err
	}
	return nil
}
