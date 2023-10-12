package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"cheatsheet/model/alfred"

	"github.com/mozillazg/go-pinyin"
)

var _pinyinArgs pinyin.Args = pinyin.NewArgs()

func init() {
	_pinyinArgs.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
}

func main() {
	repo := os.Getenv("REPO")
	if repo == "" {
		fmt.Printf("REPO environment variable is required")
		return
	}

	var queryNamespace, keyword string
	if len(os.Args) <= 2 {
		keyword = os.Args[1]
	} else {
		queryNamespace, keyword = os.Args[1], os.Args[2]
	}

	matched := []alfred.Item{}

	files, err := os.ReadDir(repo)
	if err != nil {
		fmt.Printf("tranversal %s error: %v", repo, err)
		return
	}
	for _, f := range files {
		if f.IsDir() {
			continue
		}
		if !strings.HasSuffix(f.Name(), ".txt") {
			continue
		}
		if !strings.HasPrefix(f.Name(), queryNamespace) {
			continue
		}

		filename := filepath.Join(repo, f.Name())
		namespace := f.Name()[:len(f.Name())-len(filepath.Ext(f.Name()))]
		data, err := os.ReadFile(filename)
		if err != nil {
			fmt.Printf("cannot read %s: %v\n", filename, err)
			return
		}

		entries := strings.Split(string(data), ";;;")
		for _, entry := range entries {
			if entry = strings.TrimSpace(entry); entry == "" {
				continue
			}
			lines := strings.Split(entry, "\n")
			title := strings.TrimSpace(lines[0])
			if (queryNamespace != "" && keyword == "-") || fuzzyMatch(title, keyword) {
				content := strings.TrimSpace(strings.Join(lines[1:], "\n"))
				desc := "@" + namespace
				item := alfred.Item{
					Title:    title,
					Subtitle: desc,
					Arg:      content,
				}
				matched = append(matched, item)
			}
		}
	}

	debugMsg := fmt.Sprintf("args: %v, queryNamespace: %v, keyword: %v", os.Args, queryNamespace, keyword)
	res, _ := json.Marshal(alfred.Items{Items: matched, DebugMsg: debugMsg})
	fmt.Print(string(res))
}

func fuzzyMatch(content, keyword string) bool {
	if keyword == "" {
		return false
	}

	pinyin := pinyin.LazyPinyin(content, _pinyinArgs)
	pattern := "(?i)" + strings.Join(strings.Split(keyword, ""), ".*?")
	match, err := regexp.MatchString(pattern, strings.Join(pinyin, ""))
	if err != nil {
		panic(err)
	}
	return match
}
