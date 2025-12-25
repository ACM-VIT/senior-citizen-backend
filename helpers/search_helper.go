package helpers

import "strings"

func Trigrams(s string) []string {
	s = strings.ToLower(s)
	s = strings.ReplaceAll(s, " ", "")
	s = strings.TrimSpace(s)

	if len(s) < 3 {
		return nil
	}

	seen := make(map[string]struct{})
	var grams []string

	for i := 0; i <= len(s)-3; i++ {
		g := s[i : i+3]
		if _, ok := seen[g]; !ok {
			seen[g] = struct{}{}
			grams = append(grams, g)
		}
	}

	return grams
}
