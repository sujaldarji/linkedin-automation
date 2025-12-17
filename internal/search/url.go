package search

import (
	"net/url"
)

// BuildPeopleSearchURL builds a LinkedIn people search URL
func BuildPeopleSearchURL(input *Input) string {
	params := url.Values{}
	params.Set("keywords", input.Keywords)

	return "https://www.linkedin.com/search/results/people/?" + params.Encode()
}
