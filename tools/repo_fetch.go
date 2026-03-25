package main

import (
	"encoding/json"
	"net/http"
	"os"
	"sort"
)

type Repo struct {
	Name        string `json:"name"`
	URL         string `json:"html_url"`
	Description string `json:"description"`
	PushedAt    string `json:"pushed_at"`
	Fork        bool   `json:"fork"`
}

func main() {
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

	// optional: limit to top 10
	if len(filtered) > 10 {
		filtered = filtered[:10]
	}

	file, _ := os.Create("projects.json")
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	encoder.Encode(filtered)
}
