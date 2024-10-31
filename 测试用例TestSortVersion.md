# 测试用例TestSortVersion

首先我发现他将version专门做了一个字文件夹放在server文件夹下面，这样一眼看过去简单易懂容易阅读
和测试

### 比较版本大小

```go
func IsVersionGreaterOrEqualThan(version, target string) bool {
	return semver.Compare(fmt.Sprintf("v%s", version), fmt.Sprintf("v%s", target)) > -1
}
```

这个方法厉害的是简洁易懂，就是从方法名知道他的意图，这个意图比较长，从方法的实现他借助了go语言
的semver包，这个包会以人能阅读的产品版本，而且借助第三方稳定的包是偷懒的好办法


### 版本排序

```go
func TestSortVersion(t *testing.T) {
	tests := []struct {
		versionList []string
		want        []string
	}{
		{
			versionList: []string{"0.9.1", "0.10.0", "0.8.0"},
			want:        []string{"0.8.0", "0.9.1", "0.10.0"},
		},
		{
			versionList: []string{"1.9.1", "0.9.1", "0.10.0", "0.8.0"},
			want:        []string{"0.8.0", "0.9.1", "0.10.0", "1.9.1"},
		},
	}
	for _, test := range tests {
		sort.Sort(SortVersion(test.versionList))
		assert.Equal(t, test.versionList, test.want)
	}
}
```

实现代码很有趣

```go
type SortVersion []string

func (s SortVersion) Len() int {
	return len(s)
}

func (s SortVersion) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

func (s SortVersion) Less(i, j int) bool {
	v1 := fmt.Sprintf("v%s", s[i])
	v2 := fmt.Sprintf("v%s", s[j])
	return semver.Compare(v1, v2) == -1
}

```

这里面的SortVersion结构体相关方法（比如Len, Swap 和 Less） 是为了让版本字符串可以使用
Go标准库的sort.Sort 方法进行排序，这样方便的使用Go标准库对比，并保持代码整洁可读