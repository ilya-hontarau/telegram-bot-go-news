package cache

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/illfate/telegram-bot-go-news/pkg/config"

	"github.com/stretchr/testify/assert"
)

func habr(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/xml")
	file, err := os.Open("test.xml")
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	defer file.Close()

	_, err = io.Copy(w, file)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

func TestScrapePosts(t *testing.T) {
	tt := []struct {
		Name     string
		Category string
		Link     string
	}{
		{
			Name:     "go",
			Category: "go",
			Link:     "https://habr.com/ru/post/461723/",
		},
		{
			Name:     "algorithm",
			Category: "алгоритмы",
			Link:     "https://habr.com/ru/post/461767/",
		},
	}
	s := httptest.NewServer(http.HandlerFunc(habr))
	conf, err := config.New("test-config.yml")
	if err != nil {
		t.Fatal("Couldn't get config file")
	}
	cache := New(*conf)
	cache.ScrapePosts(s.URL)
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			result := cache.postsCache[tc.Category][0]
			assert.Equal(t, tc.Link, result.Link)
		})
	}
}

func TestGetLink(t *testing.T) {
	tt := []struct {
		Name     string
		Result   string
		Category string
		UserName string
	}{
		{
			Name:     "correct_test",
			Result:   "https://habr.com/ru/post/461467/",
			Category: "Go",
			UserName: "Ilya",
		},
		{
			Name:     "incorrect_link",
			Category: "go",
			UserName: "Sasha",
		},
		{
			Name:     "synonym",
			Result:   "https://habr.com/ru/post/461723/",
			Category: "Golang",
			UserName: "Max",
		},
		{
			Name:     "synonym_in_russian",
			Result:   "https://habr.com/ru/post/461723/",
			Category: "cryptography",
			UserName: "Vlad",
		},
	}
	conf, err := config.New("test-config.yml")
	if err != nil {
		t.Fatal("Couldn't get config file")
	}
	cache := New(*conf)
	cache.postsCache["go"] = []Post{
		{
			Link: "https://habr.com/ru/post/461723/",
		},
		{
			Link: "https://habr.com/ru/post/461467/",
		},
	}
	cache.postsCache["криптография"] = []Post{
		{
			Link: "https://habr.com/ru/post/461723/",
		},
	}
	cache.userUrls["Ilya"] = []string{
		"https://habr.com/ru/post/461723/",
		"https://habr.com/ru/post/461545/",
	}
	cache.userUrls["Sasha"] = []string{
		"https://habr.com/ru/post/461723/",
		"https://habr.com/ru/post/461467/",
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			result := cache.GetLink(tc.Category, tc.UserName)
			assert.Equal(t, tc.Result, result)
		})
	}
}

func TestDeleteOldPosts(t *testing.T) {
	tt := []struct {
		Name     string
		LifeTime time.Duration
		Link     string
		Deleted  bool
	}{
		{
			Name:     "delete_test",
			LifeTime: 10 * time.Minute,
			Link:     "https://habr.com/ru/post/461467/",
			Deleted:  true,
		},
		{
			Name:     "non_delete_test",
			LifeTime: 5 * time.Minute,
			Link:     "https://habr.com/ru/post/311467/",
			Deleted:  false,
		},
	}
	conf, err := config.New("test-config.yml")
	if err != nil {
		t.Fatal("Couldn't get config file")
	}
	cache := New(*conf)
	cache.postsCache["go"] = []Post{
		{
			Link:    "https://habr.com/ru/post/501723/",
			AddedAt: time.Now(),
		},
		{
			Link:    "https://habr.com/ru/post/461467/",
			AddedAt: time.Now().Add(-10 * time.Minute),
		},
		{
			Link:    "https://habr.com/ru/post/311467/",
			AddedAt: time.Now().Add(-4 * time.Minute),
		},
	}
	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			cache.deleteOldPosts(tc.LifeTime)
			find := postsHasLink(cache.postsCache["go"], tc.Link)
			assert.NotEqual(t, tc.Deleted, find)
		})
	}
}
