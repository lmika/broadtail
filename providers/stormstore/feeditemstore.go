package stormstore

import (
	"context"
	"strings"

	"github.com/asdine/storm/v3"
	"github.com/asdine/storm/v3/q"
	"github.com/google/uuid"
	"github.com/lmika/broadtail/models"
	"github.com/pkg/errors"
)

type FeedItemStore struct {
	db *storm.DB
}

func (f *FeedItemStore) PutIfAbsent(ctx context.Context, item *models.FeedItem) error {
	if err := f.db.Select(q.Eq("EntryID", item.EntryID)).First(&models.FeedItem{}); err != nil {
		if !errors.Is(err, storm.ErrNotFound) {
			return err
		}
	} else {
		// Item exists.  Do nothing
		return nil
	}

	item.ID = uuid.New()

	return f.db.Save(item)
}

func NewFeedItemStore(filename string) (*FeedItemStore, error) {
	db, err := storm.Open(filename)
	if err != nil {
		return nil, err
	}

	return &FeedItemStore{db: db}, nil
}

func (f *FeedItemStore) ListRecentsFromAllFeeds(ctx context.Context, limit int) (feedItems []models.FeedItem, err error) {
	err = f.db.Select().OrderBy("Published").Reverse().Limit(limit).Find(&feedItems)
	if err == storm.ErrNotFound {
		return []models.FeedItem{}, nil
	}
	return feedItems, err
}

func (f *FeedItemStore) ListRecent(ctx context.Context, feedID uuid.UUID, filterExpression models.FeedItemFilter) (feedItems []models.FeedItem, err error) {
	query := q.Eq("FeedID", feedID)
	if len(filterExpression.ContainKeyword) > 0 {
		query = q.And(query, q.NewFieldMatcher("Title", fieldContainsAnyCase(filterExpression.ContainKeyword)))
	}

	err = f.db.Select(query).OrderBy("Published").Reverse().Limit(50).Find(&feedItems)
	if err == storm.ErrNotFound {
		return []models.FeedItem{}, nil
	}
	return feedItems, err
}

func (f *FeedItemStore) Close() {
	f.db.Close()
}

func fieldContainsAnyCase(tokens []string) q.FieldMatcher {
	lowerCaseTokens := make([]string, len(tokens))
	for i, t := range tokens {
		lowerCaseTokens[i] = strings.ToLower(t)
	}
	return fieldContainsAnyCaseMatcher(lowerCaseTokens)
}

type fieldContainsAnyCaseMatcher []string

func (f fieldContainsAnyCaseMatcher) MatchField(v interface{}) (bool, error) {
	s, isS := v.(string)
	if !isS {
		return false, nil
	}

	lc := strings.ToLower(s)

	for _, t := range f {
		if !strings.Contains(lc, t) {
			return false, nil
		}
	}
	return true, nil
}
