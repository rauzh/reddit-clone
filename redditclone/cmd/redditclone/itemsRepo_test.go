//nolint:all
package main

import (
	"redditclone/pkg/items"
	"redditclone/pkg/session"
	"sync"
	"testing"
)

var itemTest items.Item = items.Item{
	Category: "music",
	Text:     "text",
	Title:    "lkaej",
	Type:     "string `json",

	ID:          "1",
	Description: "-",
	URL:         "url",
	Views:       324,
	Created:     "dsklfjsldkfjsldkf",
	Author:      items.Author{Username: "CHAPA", ID: "228"},

	Score:            22,
	UpvotePercentage: 78,
	Votes:            []*items.Vote{{User: "bibp", Vote: 1}, {User: "boba", Vote: -1}},

	Comments: []*items.Comment{{Author: items.Author{Username: "CHAPA", ID: "228"}, Body: "lksdjf", Created: "dsklfjsldkfjsldkf", ID: "1"}},
}

func TestAddAsyncRace(t *testing.T) {
	ItemsRepo := items.NewRepo()
	sess := session.Session{UserID: "228", Username: "CHAPA", Token: "kdfjklsdjf"}
	item_1 := items.NewItem()
	item_2 := items.NewItem()

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ItemsRepo.Add(&sess, item_1)
	}()
	ItemsRepo.Add(&sess, item_2)
	wg.Wait()
}

func TestGetAllAsyncRace(t *testing.T) {
	ItemsRepo := items.NewRepo()
	sess := session.Session{UserID: "228", Username: "CHAPA", Token: "kdfjklsdjf"}
	// item_1 := items.NewItem()
	// item_2 := items.NewItem()
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ItemsRepo.GetAll()
	}()
	ItemsRepo.GetAll()
	wg.Wait()
}

func TestGetByItemIDCategoryAsyncRace(t *testing.T) {
	ItemsRepo := items.NewRepo()
	sess := session.Session{UserID: "228", Username: "CHAPA", Token: "kdfjklsdjf"}
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ItemsRepo.GetByItemID("1")
	}()
	ItemsRepo.GetByItemID("1")
	wg.Wait()
}

func TestGetByCategoryCategoryAsyncRace(t *testing.T) {
	ItemsRepo := items.NewRepo()
	sess := session.Session{UserID: "228", Username: "CHAPA", Token: "kdfjklsdjf"}
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ItemsRepo.GetByCategory("music")
	}()
	ItemsRepo.GetByCategory("music")
	wg.Wait()
}

func TestGetByUserIDAsyncRace(t *testing.T) {
	ItemsRepo := items.NewRepo()
	sess := session.Session{UserID: "228", Username: "CHAPA", Token: "kdfjklsdjf"}
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ItemsRepo.GetByUserID("1")
	}()
	ItemsRepo.GetByUserID("1")
	wg.Wait()
}

func TestDeleteAsyncRace(t *testing.T) {
	ItemsRepo := items.NewRepo()
	sess := session.Session{UserID: "228", Username: "CHAPA", Token: "kdfjklsdjf"}
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)
	ItemsRepo.Add(&sess, &itemTest)

	wg := sync.WaitGroup{}

	wg.Add(1)
	go func() {
		defer wg.Done()
		ItemsRepo.Delete(&sess, "1")
	}()
	ItemsRepo.Delete(&sess, "1")
	wg.Wait()
}
