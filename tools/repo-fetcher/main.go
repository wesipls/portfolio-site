package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
)

type CommitInfo struct {
	SHA     string `json:"sha"`
	Message string `json:"message"`
	Date    string `json:"date"`
	URL     string `json:"url"`
}

type Repo struct {
	Name        string         `json:"name"`
	URL         string         `json:"html_url"`
	Description string         `json:"description"`
	PushedAt    string         `json:"pushed_at"`
	Fork        bool           `json:"fork"`
	Languages   map[string]int `json:"languages"`
	LastCommit  *CommitInfo    `json:"last_commit,omitempty"`
}

type CacheEntry struct {
	PushedAt   string         `json:"pushed_at"`
	Langs      map[string]int `json:"languages"`
	LastCommit *CommitInfo    `json:"last_commit,omitempty"`
}

type githubCommitResponse struct {
	SHA     string `json:"sha"`
	HTMLURL string `json:"html_url"`
	Commit  struct {
		Message string `json:"message"`
		Author  struct {
			Date string `json:"date"`
		} `json:"author"`
	} `json:"commit"`
}

func loadCache() map[string]CacheEntry {
	file, err := os.Open("cache.json")
	if err != nil {
		return make(map[string]CacheEntry)
	}
	defer file.Close()

	var cache map[string]CacheEntry
	if err := json.NewDecoder(file).Decode(&cache); err != nil {
		return make(map[string]CacheEntry)
	}
	return cache
}

func saveCache(cache map[string]CacheEntry) {
	file, err := os.Create("cache.json")
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(cache)
}

func fetchLanguages(repoName string) map[string]int {
	url := "https://api.github.com/repos/wesipls/" + repoName + "/languages"

	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil
	}

	var langs map[string]int
	if err := json.NewDecoder(res.Body).Decode(&langs); err != nil {
		return nil
	}
	return langs
}

func fetchLastCommit(repoName string) *CommitInfo {
	url := "https://api.github.com/repos/wesipls/" + repoName + "/commits?per_page=1"

	res, err := http.Get(url)
	if err != nil {
		return nil
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil
	}

	var commits []githubCommitResponse
	if err := json.NewDecoder(res.Body).Decode(&commits); err != nil {
		return nil
	}

	if len(commits) == 0 {
		return nil
	}

	return &CommitInfo{
		SHA:     commits[0].SHA,
		Message: commits[0].Commit.Message,
		Date:    commits[0].Commit.Author.Date,
		URL:     commits[0].HTMLURL,
	}
}

func main() {
	cache := loadCache()

	res, err := http.Get("https://api.github.com/users/wesipls/repos")
	if err != nil {
		return
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return
	}

	var repos []Repo
	if err := json.NewDecoder(res.Body).Decode(&repos); err != nil {
		return
	}

	var filtered []Repo
	for _, repo := range repos {
		if repo.Fork {
			continue
		}
		filtered = append(filtered, repo)
	}

	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].PushedAt > filtered[j].PushedAt
	})

	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	for i, repo := range filtered {
		cached, exists := cache[repo.Name]

		if exists && cached.PushedAt == repo.PushedAt {
			filtered[i].Languages = cached.Langs
			filtered[i].LastCommit = cached.LastCommit
			continue
		}

		langs := fetchLanguages(repo.Name)
		lastCommit := fetchLastCommit(repo.Name)

		filtered[i].Languages = langs
		filtered[i].LastCommit = lastCommit

		cache[repo.Name] = CacheEntry{
			PushedAt:   repo.PushedAt,
			Langs:      langs,
			LastCommit: lastCommit,
		}
	}

	saveCache(cache)

	file, err := os.Create("projects.json")
	if err != nil {
		return
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	_ = encoder.Encode(filtered)
}
