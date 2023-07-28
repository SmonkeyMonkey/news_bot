package fetcher

import (
	"context"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/smonkeymonkey/news-bot/internal/model"
	src "github.com/smonkeymonkey/news-bot/internal/source"
	"golang.org/x/exp/slices"
)

type Fetcher struct {
	articles ArticleStorage
	sources  SourceProvider

	fetchInvertal  time.Duration
	filterKeywords []string
}
type ArticleStorage interface {
	Store(ctx context.Context, article model.Article) error
}
type SourceProvider interface {
	Sources(ctx context.Context) ([]model.Source, error)
}
type Source interface {
	ID() int64
	Name() string
	Fetch(ctx context.Context) ([]model.Item, error)
}

func New(articlestorage ArticleStorage,
	sourceprovider SourceProvider,
	fetchInterval time.Duration,
	filterKeywords []string) *Fetcher {
	return &Fetcher{
		articles:       articlestorage,
		sources:        sourceprovider,
		fetchInvertal:  fetchInterval,
		filterKeywords: filterKeywords,
	}
}

func (f *Fetcher) Fetch(ctx context.Context) error {
	sources, err := f.sources.Sources(ctx)
	if err != nil {
		fmt.Printf("error when get sources: %v", err)
		return err
	}

	var wg sync.WaitGroup
	for _, source := range sources {
		wg.Add(1)

		rssSource := src.NewRssSourceFromModel(source)

		go func(source Source) {
			defer wg.Done()

			items, err := source.Fetch(ctx)
			if err != nil {
				log.Printf("error fetching items from source.Error mgs: %v", err)
				return
			}
			if err := f.processItems(ctx, source, items); err != nil {
				log.Println("error processing items from source")
				return
			}
		}(rssSource)
	}
	wg.Wait()
	return nil
}

func (f *Fetcher) Start(ctx context.Context) error {
	ticker := time.NewTicker(f.fetchInvertal)
	defer ticker.Stop()
	if err := f.Fetch(ctx); err != nil {
		log.Printf("error fetching")
		return err
	}

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()

		case <-ticker.C:
			if err := f.Fetch(ctx); err != nil {
				return err
			}
		}
	}
}
func (f *Fetcher) processItems(ctx context.Context, source Source, items []model.Item) error {
	for _, item := range items {
		item.Date = item.Date.UTC()

		if f.skipItem(item) {
			log.Printf("item %s should be skipped", item.Title)
			continue
		}

		if err := f.articles.Store(ctx, model.Article{
			SourceID:    source.ID(),
			Title:       item.Title,
			Link:        item.Link,
			PublishedAt: item.Date,
		}); err != nil {
			return err
		}

	}
	return nil
}

// if the item exists in filter keywords - it should be skipped
func (f *Fetcher) skipItem(item model.Item) bool {
	for _, keyword := range f.filterKeywords {
		if slices.Contains(item.Categories, keyword) || strings.Contains(strings.ToLower(item.Title), keyword) {
			return true
		}
	}
	return false
}
