package scanner

type byPriority []string

var (
	weighted = map[string]int{
		"unfork-action":       0,
		"unstable-docker-tag": 1,
		"unstable-github-ref": 2,
		"outdated-action":     3,
	}
)

func (s byPriority) Len() int {
	return len(s)
}
func (s byPriority) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}
func (s byPriority) Less(i, j int) bool {
	return weighted[s[i]] < weighted[s[j]]
}
