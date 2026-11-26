package stringutils

func Reverse(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func LongestString(strs []string) string {
	if len(strs) == 0 {
		return ""
	}
	longest := strs[0]
	for _, s := range strs[1:] {
		if len([]rune(s)) > len([]rune(longest)) {
			longest = s
		}
	}
	return longest
}
