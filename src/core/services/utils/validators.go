package utils

import "regexp"

func IsUrlValid(url string) bool {
	regex := regexp.MustCompile(`^(http|https)://[a-zA-Z0-9.-]+(:[0-9]+)?(/.*)?$`)
	return regex.MatchString(url)
}
