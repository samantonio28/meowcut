package utils

import "regexp"

var urlRegex = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)

func IsValidURL(url string) bool {
	return urlRegex.MatchString(url)
}

func IsValidShortID(shortID string) bool {
	if len(shortID) != 10 {
		return false
	}
	for _, ch := range shortID {
		if !((ch >= 'A' && ch <= 'Z') || (ch >= 'a' && ch <= 'z') || (ch >= '0' && ch <= '9') || ch == '_') {
			return false
		}
	}
	return true
}