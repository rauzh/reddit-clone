package items

import (
	"errors"
	"fmt"
	"redditclone/pkg/session"
	"sort"
	"sync"
	"time"
)

var (
	ErrNotPostAuthor error = errors.New("not post author")
)

type ItemsRepo struct {
	mu     *sync.RWMutex
	lastID uint32
	data   []*Item
}

func NewRepo() *ItemsRepo {
	return &ItemsRepo{
		mu:     &sync.RWMutex{},
		lastID: 0,
		data:   make([]*Item, 0, 10),
	}
}

func (repo *ItemsRepo) GetAll() ([]*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()

	sort.Slice(repo.data, func(i, j int) bool {
		if repo.data[i].Score != repo.data[j].Score {
			return repo.data[i].Score > repo.data[j].Score
		}
		return repo.data[i].ID > repo.data[j].ID // если равный рейтинг, то по айди
	})

	return repo.data, nil
}

func (repo *ItemsRepo) GetByCategory(categoryName string) ([]*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	reqItems := make([]*Item, 0)
	for _, item := range repo.data {
		if item.Category == categoryName {
			reqItems = append(reqItems, item)
		}
	}

	sort.Slice(repo.data, func(i, j int) bool {
		if repo.data[i].Score != repo.data[j].Score {
			return repo.data[i].Score > repo.data[j].Score
		}
		return repo.data[i].ID > repo.data[j].ID // если равный рейтинг, то по айди
	})

	return reqItems, nil
}

func (repo *ItemsRepo) GetByItemID(id string) (*Item, error) {
	repo.mu.RLock()
	defer repo.mu.RUnlock()
	for _, item := range repo.data {
		if item.ID == id {
			return item, nil
		}
	}
	return nil, nil
}

func (repo *ItemsRepo) GetByUserID(username string) ([]*Item, error) {
	repo.mu.Lock()
	defer repo.mu.Unlock()
	reqItems := make([]*Item, 0)
	for _, item := range repo.data {
		if item.Author.Username == username {
			reqItems = append(reqItems, item)
		}
	}

	sort.Slice(repo.data, func(i, j int) bool {
		if repo.data[i].Score != repo.data[j].Score {
			return repo.data[i].Score > repo.data[j].Score
		}
		return repo.data[i].ID > repo.data[j].ID // если равный рейтинг, то по айди
	})

	return reqItems, nil
}

func (repo *ItemsRepo) Add(sess *session.Session, item *Item) (uint32, error) {
	item.Author.Username = sess.Username
	item.Author.ID = sess.UserID
	item.Created = time.Now().Format("2006-01-02 15:04:05")

	defaultVote := 0
	item.NewVote(defaultVote, sess.Username)

	repo.mu.Lock()
	repo.lastID++
	item.ID = fmt.Sprintf("%d", repo.lastID)
	repo.data = append(repo.data, item)
	lastID := repo.lastID
	repo.mu.Unlock()

	return lastID, nil
}

func (repo *ItemsRepo) Delete(sess *session.Session, id string) (bool, error) {
	i := -1
	repo.mu.Lock()
	defer repo.mu.Unlock()
	for idx, item := range repo.data {
		if item.ID == id && sess.Username != item.Author.Username {
			return false, ErrNotPostAuthor
		}
		if item.ID == id {
			i = idx
			break
		}
	}
	if i < 0 {
		return false, nil
	}

	if i < len(repo.data)-1 {
		copy(repo.data[i:], repo.data[i+1:])
	}
	repo.data[len(repo.data)-1] = nil
	repo.data = repo.data[:len(repo.data)-1]

	return true, nil
}

func (repo *ItemsRepo) Unvote(sess *session.Session, item *Item) error {
	item.muVt.Lock()
	_, err := item.DeleteVote(sess.UserID)
	item.muVt.Unlock()
	if err != nil {
		return fmt.Errorf("DB err")
	}
	return nil
}

func (repo *ItemsRepo) Upvote(sess *session.Session, item *Item) error {
	item.muVt.Lock()
	_, err := item.DeleteVote(sess.UserID)
	item.muVt.Unlock()
	if err != nil {
		return fmt.Errorf("DB err")
	}

	item.muVt.Lock()
	item.NewVote(1, sess.UserID)
	item.muVt.Unlock()
	return nil
}

func (repo *ItemsRepo) Downvote(sess *session.Session, item *Item) error {
	item.muVt.Lock()
	_, err := item.DeleteVote(sess.UserID)
	item.muVt.Unlock()
	if err != nil {
		return fmt.Errorf("DB err")
	}

	item.muVt.Lock()
	item.NewVote(-1, sess.UserID)
	item.muVt.Unlock()
	return nil
}

func (repo *ItemsRepo) AddComment(sess *session.Session, item *Item, message string) (string, error) {
	item.muCm.Lock()
	id, err := item.NewComment(sess.Username, sess.UserID, message)
	item.muCm.Unlock()
	return id, err
}

func (repo *ItemsRepo) DeleteComment(sess *session.Session, comID string, item *Item) error {
	item.muCm.Lock()
	defer item.muCm.Unlock()
	_, err := item.DeleteComment(sess.Username, comID)
	return err
}
