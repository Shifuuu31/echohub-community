package models

import "time"

type Post struct {
	ID            int
	user_id       int
	title         string
	post_content  string
	category_id   int
	creation_date time.Time
}
