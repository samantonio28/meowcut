package utils

import "regexp"

var urlRegex = regexp.MustCompile(`^(https?|ftp)://[^\s/$.?#].[^\s]*$`)

// IsValidURL проверяет, является ли строка валидным URL.
func IsValidURL(url string) bool {
	return urlRegex.MatchString(url)
}

// IsValidShortID проверяет, что shortID состоит из 10 допустимых символов.
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