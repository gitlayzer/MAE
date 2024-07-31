package tools

import "regexp"

func RegexMatched(image string) bool {
	r := regexp.MustCompile(`^\w+:\w+$|^\w+/\w+:\w+$`)
	return r.MatchString(image)
}
