package provider

import (
	"regexp"
	"strings"

	"github.com/mozillazg/go-pinyin"
)

var _pinyinArgs pinyin.Args = pinyin.NewArgs()

func init() {
	_pinyinArgs.Fallback = func(r rune, a pinyin.Args) []string {
		return []string{string(r)}
	}
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
