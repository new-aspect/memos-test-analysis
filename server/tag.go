package server

import (
	"regexp"
	"sort"
)

var tagRegexpList = []*regexp.Regexp{regexp.MustCompile(`^#([^\s#]+?) `), regexp.MustCompile(`[^\S#]?#([^\s#]+?) `)}

func findTagListFromMemosContent(memoContent string) []string {
	tagMapSet := make(map[string]bool)
	for _, tagRegexp := range tagRegexpList {
		for _, rawTag := range tagRegexp.FindAllString(memoContent, -1) {
			tag := tagRegexp.ReplaceAllString(rawTag, "$1")
			tagMapSet[tag] = true
		}
	}
	tagList := []string{}
	for tag := range tagMapSet {
		tagList = append(tagList, tag)
	}
	sort.Strings(tagList)
	return tagList
}
