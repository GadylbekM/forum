package forum

import "forum/forum/domain"

type Repo interface {
	GetUser(username string) (domain.User, error)
	SaveSession(session domain.Session) error
	GetSession(sessionID string) (domain.Session, error)
	SaveUser(domain.User) error
	GetCommentsByUser(userID int) ([]domain.Comments, error)
	GetCreatedPosts(userID int) ([]domain.Posts, error)
	SavePosts(domain.Posts) error
	GetPosts() ([]domain.Posts, error)
	GetUserPosts(userId int) ([]domain.Posts, error)
	DeletePost(postId int) error
	DeleteComment(comment_id int) error
	GetPostByID(postID int) (domain.Posts, error)
	AddComment(domain.Comments) error
	GetComments(postId int) ([]domain.Comments, error)
	GetUserById(userId int) ([]domain.User, error)
	LikePost(postID, userID int, notification domain.Notification) error
	DislikePost(postID, userID int, notification domain.Notification) error
	GetLikedPostIDs(userID int) ([]int, error)
	GetDislikedPostIDs(userID int) ([]int, error)
	GetPostsByCategories(categories []string) ([]domain.Posts, error)
	GetUserByEmail(email string) (domain.User, error)
	LikeComment(commentID int, userID int, notification_comments domain.Notification_comments) error
	DislikeComment(commentID int, userID int, notification_comments domain.Notification_comments) error
	InvalidateSessions(userID int) error
	CreateNotification(notification domain.Notification) error
	CreateNotificationComments(notification domain.Notification_comments) error
	GetAllNotificationsComment(ownerID int) ([]domain.Notification_comments, error)
	GetAllNotifications(ownerID int) ([]domain.Notification, error)
	EditPost(postId int, post domain.Posts) error
	EditComment(commentId int, comment domain.Comments) error
}
