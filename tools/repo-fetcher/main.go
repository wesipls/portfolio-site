package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
)

type Repo struct {
	Name        string         `json:"name"`
	URL         string         `json:"html_url"`
	Description string         `json:"description"`
	PushedAt    string         `json:"pushed_at"`
	Fork        bool           `json:"fork"`
	Languages   map[string]int `json:"languages"`
}

type CacheEntry struct {
	PushedAt string         `json:"pushed_at"`
	Langs    map[string]int `json:"languages"`
}

func loadCache() map[string]CacheEntry {
	file, err := os.Open("cache.json")
	if err != nil {
		return make(map[string]CacheEntry)
	}
	defer file.Close()

	var cache map[string]CacheEntry
	json.NewDecoder(file).Decode(&cache)
	return cache
}

func saveCache(cache map[string]CacheEntry) {
	file, _ := os.Create("cache.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(cache)
}

func fetchLanguages(repoName string) map[string]int {
	url := "https://api.github.com/repos/wesipls/" + repoName + "/languages"

	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	var langs map[string]int
	json.NewDecoder(res.Body).Decode(&langs)
	return langs
}

func main() {
	cache := loadCache()

	res, _ := http.Get("https://api.github.com/users/wesipls/repos")
	defer res.Body.Close()

	var repos []Repo
	json.NewDecoder(res.Body).Decode(&repos)

	// filter out forks
	var filtered []Repo
	for _, repo := range repos {
		if repo.Fork {
			continue
		}
		filtered = append(filtered, repo)
	}

	// sort by latest commit
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PushedAt > filtered[j].PushedAt
	})

	// limit to top 10
	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	// attach languages with caching
	for i, repo := range filtered {
		cached, exists := cache[repo.Name]

		if exists && cached.PushedAt == repo.PushedAt {
			// use cached languages
			filtered[i].Languages = cached.Langs
			continue
		}

		// fetch fresh languages
		langs := fetchLanguages(repo.Name)
		filtered[i].Languages = langs

		// update cache
		cache[repo.Name] = CacheEntry{
			PushedAt: repo.PushedAt,
			Langs:    langs,
		}
	}

	saveCache(cache)

	file, _ := os.Create("projects.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(filtered)
}
