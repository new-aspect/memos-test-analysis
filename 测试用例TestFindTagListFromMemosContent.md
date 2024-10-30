# 测试用例TestFindTagListFromMemosContent

这是2022年11月18日0.7.3在server/tag.go的测试用例TestFindTagListFromMemosContent

```go
func TestFindTagListFromMemosContent(t *testing.T) {
	tests := []struct {
		memoContent string
		want        []string
	}{
		{
			memoContent: "#tag1 ",
			want:        []string{"tag1"},
		},
		{
			memoContent: "#tag1 #tag2 ",
			want:        []string{"tag1", "tag2"},
		},
		{
			memoContent: "#tag1 #tag2",
			want:        []string{"tag1"}, // tag后面必须有空格才能被识别为tag
		},
		{
			memoContent: "#tag1 #tag2 \n#tag3 ",
			want:        []string{"tag1", "tag2", "tag3"},
		},
		{
			memoContent: "#tag1 #tag2 \n#tag3 #tag4 ",
			want:        []string{"tag1", "tag2", "tag3", "tag4"},
		},
	}

	for _, test := range tests {
		result := findTagListFromMemosContent(test.memoContent)
		if len(result) != len(test.want) {
			t.Errorf("Find tag list %s: got result %v, want %v.", test.memoContent, result, test.want)
		}
	}
}
```
我发现涉及到处理字符串的逻辑，可以写一个测试简单的测试用例，然后对特殊要求加注释，比如tag后面必须要有空格才能被识别，这样你在测试运行不通过的时候，也能很方便的意识到为什么通不过

```go
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
```