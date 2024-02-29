package domain

type UserActivity struct {
	Username      string
	CreatedPosts  []Posts
	LikedItems    []Posts
	DislikedItems []Posts
	Comments      []Comments
}
