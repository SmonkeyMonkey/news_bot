package source

import (
	"context"

	"github.com/smonkeymonkey/news-bot/internal/model"

	"github.com/SlyMarbo/rss"
	"github.com/samber/lo"
)

type RSSSource struct {
	URL        string
	SourceID   int64
	SourceName string
}

func NewRssSourceFromModel(m model.Source) RSSSource {
	return RSSSource{
		URL:        m.FeedURL,
		SourceID:   m.ID,
		SourceName: m.Name,
	}
}
func (rs RSSSource) Fetch(ctx context.Context) ([]model.Item, error) {
	feed, err := rs.loadFeed(ctx, rs.URL)
	if err != nil {
		return nil, err
	}
	return lo.Map(feed.Items, func(item *rss.Item, _ int) model.Item {
		return model.Item{
			Title:      item.Title,
			Categories: item.Categories,
			Link:       item.Link,
			Date:       item.Date,
			SourceName: rs.SourceName,
		}
	}), nil
}

func (rs RSSSource) loadFeed(ctx context.Context, url string) (*rss.Feed, error) {
	var feedCh = make(chan *rss.Feed)
	var errCh = make(chan error)

	go func() {
		feed, err := rss.Fetch(url)
		if err != nil {
			errCh <- err
			return
		}
		feedCh <- feed
	}()

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case err := <-errCh:
		return nil, err
	case feed := <-feedCh:
		return feed, nil
	}
}

func (rs RSSSource) ID() int64 {
	return rs.SourceID
}

func (rs RSSSource) Name() string {
	return rs.SourceName
}
