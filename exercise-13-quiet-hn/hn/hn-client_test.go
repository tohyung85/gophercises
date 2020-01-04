package hn

import (
	"testing"
)

func TestGetTopStories(t *testing.T) {
	hnClient := &HNClient{}
	topStoriesIds, err := hnClient.GetTopStories()
	if err != nil {
		t.Errorf("Error occured: %s\n", err)
	}
	t.Logf("Success: Got %d ids\n", len(topStoriesIds))
}

func TestGetItem(t *testing.T) {
	hnClient := &HNClient{}
	item, err := hnClient.GetItemWithID(8863)
	if err != nil {
		t.Errorf("Error occured: %s\n", err)
	}
	expectedTitle := "My YC app: Dropbox - Throw away your USB drive"
	if item.Title != expectedTitle {
		t.Errorf("Failed: got title of %s vs %s\n", item.Title, expectedTitle)
	}
	expectedUrl := "http://www.getdropbox.com/u/2/screencast.html"
	if item.URL != expectedUrl {
		t.Errorf("Failed: got url of %s vs %s\n", item.URL, expectedUrl)
	}

	expectedDomain := "getdropbox.com"
	if item.Host != expectedDomain {
		t.Errorf("Failed: got host of %s vs %s\n", item.Host, expectedDomain)
	}

	t.Logf("Success: Got story:%s", item)
}
