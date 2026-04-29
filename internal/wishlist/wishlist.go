package wishlist

import "time"

type Service struct {
	repository *Repository
}

type Wishlist struct {
	User_id int
	Event string
	Description string
	Date time.Time
}



