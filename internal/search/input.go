package search

import "fmt"

// Input represents search criteria for LinkedIn people search
type Input struct {
	Keywords string
	PageLimit int
}

// Validate ensures search input is usable
func (in *Input) Validate() error {
	if in.Keywords == "" {
		return fmt.Errorf("search keywords cannot be empty")
	}

	if in.PageLimit <= 0 {
		in.PageLimit = 1 // sensible default
	}

	return nil
}
