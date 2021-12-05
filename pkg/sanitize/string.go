package sanitize

import "regexp"

var nonAlphaNum = regexp.MustCompile(`[^\w_-]+`)

func String(s string) string {
	return nonAlphaNum.ReplaceAllString(s, "-")
}
