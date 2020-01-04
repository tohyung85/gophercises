package hnhandler

import (
	"fmt"
	"html/template"
	"net/http"
	"path/filepath"
	"sort"
	"sync"
	"time"

	"github.com/tohyung85/gophercises/exercise-13-quiet-hn/hn"
)

const (
	refreshTimeInMins = 10
	cacheValidityMins = 15
)

type Page struct {
	Time    float64
	Stories []hn.HNStory
}

type HNHandler struct {
	templ       *template.Template
	activeCache *Cache
	caches      []*Cache
	mutex       sync.Mutex
}

type SafeStories struct {
	stories []hn.HNStory
	mux     sync.Mutex
}

type Cache struct {
	index     int
	stories   []hn.HNStory
	lastSaved time.Time
}

func NewHandler() *HNHandler {
	filePath, _ := filepath.Abs("/Users/joshuatan/Go/src/github.com/tohyung85/gophercises/exercise-13-quiet-hn/templates/home.html")
	tmpl, err := template.ParseFiles(filePath)
	if err != nil {
		panic(err)
	}
	cache1 := &Cache{index: 0}
	cache2 := &Cache{index: 1}
	handler := &HNHandler{tmpl, cache1, []*Cache{cache1, cache2}, sync.Mutex{}}
	err = handler.refreshCache(cache1)
	if err != nil {
		fmt.Println("Error creating first api call!")
	}
	go handler.runCacheRefreshProgram()
	go handler.runActiveCacheSwitchProgram()
	return handler
}

func (handler *HNHandler) runActiveCacheSwitchProgram() {
	for {
		time.Sleep(refreshTimeInMins * time.Minute)
		handler.mutex.Lock()
		nextCacheIdx := handler.activeCache.index + 1
		if nextCacheIdx == len(handler.caches) {
			nextCacheIdx = 0
		}
		fmt.Printf("Switching to cache at index: %d\n", nextCacheIdx)
		handler.activeCache = handler.caches[nextCacheIdx]
		handler.mutex.Unlock()
	}
}

func (handler *HNHandler) runCacheRefreshProgram() {
	for {
		for _, cache := range handler.caches {
			err := handler.refreshCache(cache)
			if err != nil {
				fmt.Println("Error refreshing cache at index", cache.index)
			}
			fmt.Printf("Refreshed cache at index: %d\n", cache.index)
			time.Sleep(refreshTimeInMins * time.Minute)
		}
	}
}

func (handler *HNHandler) refreshCache(cache *Cache) error {
	hnClient := hn.NewClient()
	topStoriesIds, err := hnClient.GetTopStories()
	if err != nil {
		return err
	}

	stories, err := getStoriesList(topStoriesIds, hnClient)
	if err != nil {
		return err
	}
	handler.mutex.Lock()
	cache.stories = stories
	cache.lastSaved = time.Now()
	handler.mutex.Unlock()
	return nil
}

func (handler *HNHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	if startTime.Sub(handler.activeCache.lastSaved).Minutes() > cacheValidityMins {
		fmt.Println("Getting from api path")
		err := handler.refreshCache(handler.activeCache)
		if err != nil {
			fmt.Println("Issue with getting from api path, fallback to previous cache version")
		}
	} else {
		fmt.Println("Pull from cache")
	}
	handler.mutex.Lock()
	stories := handler.activeCache.stories
	handler.mutex.Unlock()
	endTime := time.Now()
	timeElapsed := endTime.Sub(startTime)
	page := &Page{timeElapsed.Seconds(), stories}

	handler.templ.Execute(w, page)

	return
}

func getStoriesList(ids []int, hnClient *hn.HNClient) ([]hn.HNStory, error) {
	fmt.Println("Getting stories")
	stories := make([]hn.HNStory, 0)
	ss := &SafeStories{stories: stories}
	completed := make(chan bool)
	toProcess := 30
	currentIdx := 0
	var wg sync.WaitGroup
	for toProcess > 0 {
		for i := currentIdx; i < currentIdx+toProcess; i++ {
			wg.Add(1)
			id := ids[i]
			go processId(id, i, completed, ss, hnClient, &wg)
		}
		wg.Wait()
		currentIdx = currentIdx + toProcess
		toProcess = toProcess - ss.numberOfStories()
	}
	ss.sortStoriesByRank()
	return ss.stories, nil
}

func processId(id int, rank int, c chan bool, ss *SafeStories, hnClient *hn.HNClient, wg *sync.WaitGroup) {
	defer wg.Done()
	story, err := hnClient.GetItemWithID(id)
	if err != nil {
		fmt.Printf("Error processing id: %d", id)
		return
	}
	if story.Type != "story" || story.URL == "" {
		return
	}
	story.Rank = rank
	ss.addStory(story)
}

func (ss *SafeStories) sortStoriesByRank() {
	sort.Slice(ss.stories, func(i, j int) bool {
		if ss.stories[i].Rank < ss.stories[j].Rank {
			return true
		}
		return false
	})
}

func (ss *SafeStories) numberOfStories() int {
	ss.mux.Lock()
	defer ss.mux.Unlock()

	return len(ss.stories)
}

func (ss *SafeStories) addStory(story hn.HNStory) {
	ss.mux.Lock()
	defer ss.mux.Unlock()
	ss.stories = append(ss.stories, story)
}
