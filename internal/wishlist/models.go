package wishlist

import "time"

type WishlistModel struct {
	Id int
	User_id int
	Event string
	Description string
	Date time.Time
}