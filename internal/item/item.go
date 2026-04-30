package item

type Service struct {
	repository *Repository
}

type Item struct {
	Id int
	Wishlist_id int
	Title string
	Description string
	URL string
	Priority int
	IsReserved int
}

