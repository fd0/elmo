package gitolite

import (
	"path"
	"strings"
)

// Filter returns a list of repositories which match the pattern case
// insensitive.
func Filter(repos Repos, pattern string) Repos {
	if pattern == "" {
		return repos
	}

	res := make(Repos, len(repos))

	include := filter(repos.Names(), pattern)
	for _, name := range include {
		res[name] = repos[name]
	}

	return res
}

func filter(names []string, pattern string) (res []string) {
	pattern = strings.ToLower(pattern)
	patterns := strings.Split(pattern, "/")

	for _, name := range names {
		lowerName := strings.ToLower(name)
		if strings.Contains(lowerName, pattern) {
			res = append(res, name)
			continue
		}

		strs := strings.Split(lowerName, "/")
		if match(patterns, strs) {
			res = append(res, name)
			continue
		}

		// try each path component
		for _, pc := range strings.Split(lowerName, "/") {
			m, err := path.Match(pattern, pc)
			if err != nil {
				panic(err)
			}

			if m {
				res = append(res, name)
				continue
			}
		}
	}

	return res
}

func match(patterns, strs []string) bool {
	if len(patterns) == 0 && len(strs) == 0 {
		return true
	}

	if len(patterns) <= len(strs) {
	outer:
		for offset := len(strs) - len(patterns); offset >= 0; offset-- {

			for i := len(patterns) - 1; i >= 0; i-- {
				ok, err := path.Match(patterns[i], strs[offset+i])
				if err != nil {
					panic(err)
				}

				if !ok {
					continue outer
				}
			}

			return true
		}
	}

	return false
}

// ValidatePattern returns an error if the pattern is invalid.
func ValidatePattern(pattern string) error {
	_, err := path.Match(pattern, "")
	return err
}
