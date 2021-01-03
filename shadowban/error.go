package shadowban

import "fmt"

// RedditError is a specialized error if Reddit falls over
type RedditError struct {
	message    string
	StatusCode int
}

func (r *RedditError) Error() string {
	return fmt.Sprintf("Reddit is down. Status code: %d. %s", r.StatusCode, r.message)
}
