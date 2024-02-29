package forum

import (
	"forum/forum/domain"

	"github.com/gofrs/uuid"
)

type Business interface {
	Login(username, password string) (uuid.UUID, error)
	Registration(username, password, email string) error
	Session(sessionID string) (*domain.Session, error)
	GetUserActivity(userID int) (domain.UserActivity, error)
	Post(domain.Posts) error
	GetAllPosts() ([]domain.Posts, error)
	GetMyPosts(userId int) ([]domain.Posts, error)
	DeletePost(postId int) error

	GetPostByID(postId int) (domain.Posts, error)
	AddComment(comment domain.Comments, userID, postID int) error
	GetComments(postId int) ([]domain.Comments, error)
	DeleteComment(comment_id int) error
	GetUserById(userId int) ([]domain.User, error)
	LikePost(postID, userID, ownerID int, activity string, username string) error
	DislikePost(postID, userID, ownerID int, activity string, username string) error
	GetLikedPosts(userID int) ([]domain.Posts, error)
	GetDislikedPosts(userID int) ([]domain.Posts, error)
	GetPostsByCategories(categories []string) ([]domain.Posts, error)
	LikeComment(commentID int, userID int, ownerID int, activity string, username string, postID int) error
	DislikeComment(commentID int, userID int, ownerID int, activity string, username string, postID int) error
	GetAllNotifications(ownerID int) ([]domain.Notification, error)
	GetAllNotificationsComment(ownerID int) ([]domain.Notification_comments, error)
	EditPost(postId int, post domain.Posts) error
	EditComment(commentId int, comment domain.Comments) error
}
