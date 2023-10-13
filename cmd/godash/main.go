package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/callmemhz/godash/internal/provider"
	"github.com/callmemhz/godash/model"
	"github.com/callmemhz/godash/model/alfred"
	"github.com/samber/lo"
)

func main() {

	repo := os.Getenv("REPO")
	if repo == "" {
		fmt.Printf("REPO environment variable is required")
		return
	}
	providerFactory := &provider.TxtProviderFactory{RepoPath: repo}
	provider, err := providerFactory.NewProvider()
	if err != nil {
		fmt.Println(err)
		return
	}

	var namespace, keyword string
	if len(os.Args) <= 2 {
		keyword = os.Args[1]
	} else {
		namespace, keyword = os.Args[1], os.Args[2]
	}

	entries, err := provider.Search(namespace, keyword)
	if err != nil {
		fmt.Println(err)
		return
	}

	items := lo.Map(entries, func(entry model.Entry, _ int) alfred.Item {
		return alfred.Item{
			Title:    entry.Title(),
			Subtitle: entry.Subtitle(),
			Arg:      entry.Content(),
		}
	})

	debugMsg := fmt.Sprintf("args: %v, queryNamespace: %v, keyword: %v", os.Args, namespace, keyword)
	res, _ := json.Marshal(alfred.Items{Items: items, DebugMsg: debugMsg})
	fmt.Print(string(res))
}
