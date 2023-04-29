package url

import "net/url"

func MustParse(rawURL string) *url.URL {
	value, err := url.Parse(rawURL)
	if err != nil {
		panic(err)
	}

	return value
}

func MustJoinPath(base string, elem ...string) string {
	value, err := url.JoinPath(base, elem...)
	if err != nil {
		panic(err)
	}

	return value
}
