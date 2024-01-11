package models

import "time"

type Snip struct {
	ID      int
	Title   string
	Content string
	Created time.Time
	Expires time.Time
}
