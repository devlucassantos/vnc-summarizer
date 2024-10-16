package validators

import "regexp"

func IsUrlValid(url string) bool {
	regex := regexp.MustCompile(`^(?:https?://)?(w{3}\.)?[\w_-]+((\.\w{2,}){1,3})(/([^/\n]+/?)*(\?[\w_-]+=[^?/&]*(&[\w_-]+=[^?/&]*)*)?)?$`)
	return regex.MatchString(url)
}
