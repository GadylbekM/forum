package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"forum/forum/domain"

	_ "github.com/mattn/go-sqlite3"
)

type RepoSqlLite struct {
	db *sql.DB
}

func NewDatabase() (*RepoSqlLite, error) {
	db, err := sql.Open("sqlite3", "forum.db")
	if err != nil {
		return nil, err
	}
	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			user_id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT NOT NULL UNIQUE,
			email TEXT NOT NULL UNIQUE,
			password TEXT NOT NULL 
		);
		CREATE TABLE IF NOT EXISTS posts (
			post_id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			username TEXT NOT NULL, 
			category TEXT NOT NULL,
			title TEXT NOT NULL,
			content TEXT NOT NULL,
			imagefield TEXT,

			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP
		);
		CREATE TABLE IF NOT EXISTS comments (
			comment_id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			username TEXT NOT NULL,

			content TEXT NOT NULL,
			creation_date DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
		CREATE TABLE IF NOT EXISTS likes (
			like_id INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id INTEGER NOT NULL,
			user_id INTEGER NOT NULL,
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (post_id) REFERENCES posts(post_id),
			FOREIGN KEY (user_id) REFERENCES users(user_id)
		);
		CREATE TABLE IF NOT EXISTS session (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id INTEGER NOT NULL,
			username TEXT NOT NULL,
			session_id TEXT NOT NULL
		);
	`)
	return &RepoSqlLite{
		db: db,
	}, err
}

func (r *RepoSqlLite) GetUser(username string) (domain.User, error) {
	var user domain.User
	err := r.db.QueryRow("SELECT user_id, username, email, password, registration_date FROM users WHERE username = ?", username).Scan(
		&user.UserId,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.RegistrationDate,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.User{}, domain.ErrInvalidUser
	}
	return user, nil
}

func (r *RepoSqlLite) GetUserByEmail(email string) (domain.User, error) {
	var user domain.User

	err := r.db.QueryRow("SELECT username, email, password FROM users WHERE email = ?", email).
		Scan(&user.Username, &user.Email, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return domain.User{}, domain.ErrInvalidUser
		}
		return domain.User{}, err
	}

	return user, nil
}

func (r *RepoSqlLite) SaveSession(session domain.Session) error {
	_, err := r.db.Exec("INSERT OR REPLACE INTO session (user_id, username, session_id, creation_date , expiration_date) VALUES (?,?,?,? ,?)", session.UserId, session.Username, session.SessionId, session.CreationDate, session.ExpiritionDate)
	return err
}

func (r *RepoSqlLite) GetSession(sessionID string) (domain.Session, error) {
	var session domain.Session
	err := r.db.QueryRow("SELECT user_id, username, session_id  FROM session WHERE session_id = ?", sessionID).Scan(
		&session.UserId,
		&session.Username,
		&session.SessionId,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return domain.Session{}, domain.ErrSessionNotFound
	}
	return session, err
}

func (r *RepoSqlLite) SaveUser(user domain.User) error {
	_, err := r.db.Exec("INSERT INTO users (username, email, password) VALUES (?,?,?)", user.Username, user.Email, user.Password)
	return err
}

func (r *RepoSqlLite) SavePosts(posts domain.Posts) error {
	_, err := r.db.Exec("INSERT INTO posts (user_id, username, category, title, content, category_id ,imagefield, creation_date) VALUES (?,?,?,?,?,?,?,?)", posts.UserId, posts.Username, posts.Category, posts.Title, posts.Content, posts.CategoryId, posts.ImageField, posts.CreationDate)
	return err
}

func (r *RepoSqlLite) EditPost(postId int, post domain.Posts) error {
	_, err := r.db.Exec("UPDATE posts SET title=?, content=?, imagefield=? WHERE post_id=?", post.Title, post.Content, post.ImageField, postId)

	return err
}

func (r *RepoSqlLite) EditComment(commentId int, comment domain.Comments) error {
	_, err := r.db.Exec("UPDATE comments SET content=? WHERE comment_id=?", comment.Content, commentId)

	return err
}

func (r *RepoSqlLite) GetPosts() ([]domain.Posts, error) {
	var posts []domain.Posts
	rows, err := r.db.Query("SELECT post_id, user_id, username, category, title, content, imagefield, creation_date, likes, dislikes FROM posts")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Posts
		err := rows.Scan(&p.PostId, &p.UserId, &p.Username, &p.Category, &p.Title, &p.Content, &p.ImageField, &p.CreationDate, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *RepoSqlLite) GetUserPosts(userID int) ([]domain.Posts, error) {
	var posts []domain.Posts
	rows, err := r.db.Query("SELECT post_id, user_id, username, category, title, content, imagefield, category_id, creation_date ,likes, dislikes FROM posts WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var p domain.Posts
		err := rows.Scan(&p.PostId, &p.UserId, &p.Username, &p.Category, &p.Title, &p.Content, &p.ImageField, &p.CategoryId, &p.CreationDate, &p.Likes, &p.Dislikes)
		if err != nil {
			return nil, err
		}

		posts = append(posts, p)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *RepoSqlLite) DeletePost(postID int) error {
	_, err := r.db.Exec("DELETE FROM posts WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM likes WHERE post_id = ?", postID)
	if err != nil {
		return err
	}

	_, err = r.db.Exec("DELETE FROM dislikes WHERE post_id = ?", postID)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoSqlLite) GetPostByID(postID int) (domain.Posts, error) {
	var p domain.Posts
	err := r.db.QueryRow("SELECT post_id, user_id, username, category, title, content, imagefield, creation_date, likes, dislikes FROM posts WHERE post_id = ?", postID).
		Scan(&p.PostId, &p.UserId, &p.Username, &p.Category, &p.Title, &p.Content, &p.ImageField, &p.CreationDate, &p.Likes, &p.Dislikes)
	if err != nil {

		if err == sql.ErrNoRows {
			return domain.Posts{}, fmt.Errorf("post not found")
		}
		return domain.Posts{}, err
	}

	return p, nil
}

func (r *RepoSqlLite) AddComment(comments domain.Comments) error {
	_, err := r.db.Exec("INSERT INTO comments ( post_id, user_id, content, creation_date, username) VALUES (?,?,?,?,?)", comments.PostId, comments.UserId, comments.Content, comments.CreationDate, comments.Username)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func (r *RepoSqlLite) GetComments(postId int) ([]domain.Comments, error) {
	var comments []domain.Comments
	rows, err := r.db.Query("SELECT comment_id, post_id, user_id, content, creation_date, username, likes , dislikes FROM comments WHERE post_id = ?", postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.Comments
		err := rows.Scan(&c.CommentId, &c.PostId, &c.UserId, &c.Content, &c.CreationDate, &c.Username, &c.Likes, &c.Dislikes)
		if err != nil {
			return nil, err
		}

		comments = append(comments, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

func (r *RepoSqlLite) DeleteComment(commentId int) error {
	_, err := r.db.Exec("DELETE FROM comments WHERE comment_id = ?", commentId)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoSqlLite) GetUserById(userId int) ([]domain.User, error) {
	var users []domain.User
	rows, err := r.db.Query("SELECT user_id, username, email, password, registration_date FROM users WHERE post_id = ?", userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var c domain.User
		err := rows.Scan(&c.UserId, &c.Username, &c.Email, &c.Password, &c.RegistrationDate)
		if err != nil {
			return nil, err
		}
		users = append(users, c)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}

func (r *RepoSqlLite) DeleteNotification(userId, postId int) error {
	_, err := r.db.Exec("DELETE FROM notifications WHERE post_id = ? AND user_id = ?", postId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoSqlLite) DeleteNotificationComments(userId, commentId int) error {
	_, err := r.db.Exec("DELETE FROM notifications_comments WHERE comment_id = ? AND user_id = ?", commentId, userId)
	if err != nil {
		return err
	}
	return nil
}

func (r *RepoSqlLite) LikePost(postID, userID int, notification domain.Notification) error {
	disliked, err := r.HasDislikedPost(postID, userID)
	if err != nil {
		return err
	}
	liked, err := r.HasLikedPost(postID, userID)
	if err != nil {
		return err
	}
	dislikeCount, err := r.GetDislikesCount(postID)
	if err != nil {
		return err
	}

	likesCount, err := r.GetLikesCount(postID)
	if err != nil {
		return err
	}
	if liked {
		err := r.DeleteNotification(userID, postID)
		_, err = r.db.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			return err
		}
		if likesCount > 0 {
			_, err = r.db.Exec("UPDATE posts SET likes = likes - 1 WHERE post_id = ?", postID)
			if err != nil {
				return err
			}
		}

	}
	if disliked {
		err := r.DeleteNotification(userID, postID)

		_, err = r.db.Exec("DELETE FROM dislikes WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			return err
		}
		if dislikeCount > 0 {

			_, err = r.db.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE post_id = ?", postID)
			if err != nil {
				return err
			}
		}
	}

	if !liked {
		err := r.CreateNotification(notification)

		_, err = r.db.Exec("INSERT INTO likes (post_id, user_id) VALUES (?, ?)", postID, userID)
		if err != nil {
			return err
		}

		_, err = r.db.Exec("UPDATE posts SET likes = likes + 1 WHERE post_id = ?", postID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoSqlLite) DislikePost(postID, userID int, notification domain.Notification) error {
	disliked, err := r.HasDislikedPost(postID, userID)
	if err != nil {
		return err
	}
	liked, err := r.HasLikedPost(postID, userID)
	if err != nil {
		return err
	}
	dislikeCount, err := r.GetDislikesCount(postID)
	if err != nil {
		return err
	}

	likesCount, err := r.GetLikesCount(postID)
	if err != nil {
		return err
	}
	if disliked {
		err := r.DeleteNotification(userID, postID)
		_, err = r.db.Exec("DELETE FROM dislikes WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			return err
		}
		if dislikeCount > 0 {
			_, err = r.db.Exec("UPDATE posts SET dislikes = dislikes - 1 WHERE post_id = ?", postID)
			if err != nil {
				return err
			}
		}

	}
	if liked {
		err := r.DeleteNotification(userID, postID)
		_, err = r.db.Exec("DELETE FROM likes WHERE post_id = ? AND user_id = ?", postID, userID)
		if err != nil {
			return err
		}
		if likesCount > 0 {

			_, err = r.db.Exec("UPDATE posts SET likes = likes - 1 WHERE post_id = ?", postID)
			if err != nil {
				return err
			}
		}

	}

	if !disliked {
		err := r.CreateNotification(notification)
		_, err = r.db.Exec("INSERT INTO dislikes (post_id, user_id) VALUES (?, ?)", postID, userID)
		if err != nil {
			return err
		}

		_, err = r.db.Exec("UPDATE posts SET dislikes = dislikes + 1 WHERE post_id = ?", postID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoSqlLite) HasLikedPost(postID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM likes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RepoSqlLite) HasDislikedPost(postID, userID int) (bool, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM dislikes WHERE post_id = ? AND user_id = ?", postID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RepoSqlLite) GetLikesCount(postID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT likes FROM posts WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepoSqlLite) GetDislikesCount(postID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT dislikes FROM posts WHERE post_id = ?", postID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepoSqlLite) GetLikedPostIDs(userID int) ([]int, error) {
	var likedPostIDs []int

	rows, err := r.db.Query("SELECT post_id FROM likes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID int

		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		likedPostIDs = append(likedPostIDs, postID)
	}
	return likedPostIDs, nil
}

func (r *RepoSqlLite) GetDislikedPostIDs(userID int) ([]int, error) {
	var dislikedPostIDs []int

	rows, err := r.db.Query("SELECT post_id FROM dislikes WHERE user_id = ?", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var postID int

		err := rows.Scan(&postID)
		if err != nil {
			return nil, err
		}
		dislikedPostIDs = append(dislikedPostIDs, postID)
	}
	return dislikedPostIDs, nil
}

func (r *RepoSqlLite) GetPostsByCategories(category []string) ([]domain.Posts, error) {
	var posts []domain.Posts

	for _, c := range category {
		rows, err := r.db.Query("SELECT post_id, user_id, username, category, title, content, imagefield, creation_date, likes, dislikes FROM posts WHERE category = ?", c)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var post domain.Posts
			err := rows.Scan(&post.PostId, &post.UserId, &post.Username, &post.Category, &post.Title, &post.Content, &post.ImageField, &post.CreationDate, &post.Likes, &post.Dislikes)
			if err != nil {
				return nil, err
			}

			posts = append(posts, post)
		}
	}

	return posts, nil
}

func (r *RepoSqlLite) LikeComment(commentID int, userID int, notification_comments domain.Notification_comments) error {
	disliked, err := r.HasDislikedComment(commentID, userID)
	if err != nil {
		return err
	}
	liked, err := r.HasLikedComment(commentID, userID)
	if err != nil {
		return err
	}
	dislikeCount, err := r.GetDislikesCountForComment(commentID)
	if err != nil {
		return err
	}

	likesCount, err := r.GetLikesCountForComment(commentID)
	if err != nil {
		return err
	}
	if liked {
		err := r.DeleteNotificationComments(userID, commentID)
		_, err = r.db.Exec("DELETE FROM likesforcomments WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return err
		}
		if likesCount > 0 {
			_, err = r.db.Exec("UPDATE comments SET likes = likes - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				return err
			}
		}

	}
	if disliked {
		err := r.DeleteNotificationComments(userID, commentID)

		_, err = r.db.Exec("DELETE FROM dislikesforcomments WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return err
		}
		if dislikeCount > 0 {

			_, err = r.db.Exec("UPDATE comments SET dislikes = dislikes - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				return err
			}
		}
	}

	if !liked {
		err := r.CreateNotificationComments(notification_comments)

		_, err = r.db.Exec("INSERT INTO likesforcomments (comment_id, user_id) VALUES (?, ?)", commentID, userID)
		if err != nil {
			return err
		}

		_, err = r.db.Exec("UPDATE comments SET likes = likes + 1 WHERE comment_id = ?", commentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoSqlLite) DislikeComment(commentID int, userID int, notification_comments domain.Notification_comments) error {
	disliked, err := r.HasDislikedComment(commentID, userID)
	if err != nil {
		return err
	}
	liked, err := r.HasLikedComment(commentID, userID)
	if err != nil {
		return err
	}
	dislikeCount, err := r.GetDislikesCountForComment(commentID)
	if err != nil {
		fmt.Println(err)

		return err
	}

	likesCount, err := r.GetLikesCountForComment(commentID)
	if err != nil {
		return err
	}
	if disliked {
		err := r.DeleteNotificationComments(userID, commentID)

		_, err = r.db.Exec("DELETE FROM dislikesforcomments WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return err
		}
		if dislikeCount > 0 {
			_, err = r.db.Exec("UPDATE comments SET dislikes = dislikes - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				return err
			}
		}

	}
	if liked {
		err := r.DeleteNotificationComments(userID, commentID)

		_, err = r.db.Exec("DELETE FROM likesforcomments WHERE comment_id = ? AND user_id = ?", commentID, userID)
		if err != nil {
			return err
		}
		if likesCount > 0 {

			_, err = r.db.Exec("UPDATE comments SET likes = likes - 1 WHERE comment_id = ?", commentID)
			if err != nil {
				return err
			}
		}

	}

	if !disliked {
		err := r.CreateNotificationComments(notification_comments)

		_, err = r.db.Exec("INSERT INTO dislikesforcomments (comment_id, user_id) VALUES (?, ?)", commentID, userID)
		if err != nil {
			return err
		}

		_, err = r.db.Exec("UPDATE comments SET dislikes = dislikes + 1 WHERE comment_id = ?", commentID)
		if err != nil {
			return err
		}
	}

	return nil
}

func (r *RepoSqlLite) HasLikedComment(commentID int, userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM likesforcomments WHERE comment_id = ? AND user_id = ?"
	var count int
	err := r.db.QueryRow(query, commentID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RepoSqlLite) HasDislikedComment(commentID int, userID int) (bool, error) {
	query := "SELECT COUNT(*) FROM dislikesforcomments WHERE comment_id = ? AND user_id = ?"
	var count int
	err := r.db.QueryRow(query, commentID, userID).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *RepoSqlLite) GetDislikesCountForComment(commentID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT dislikes FROM comments WHERE comment_id = ?", commentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepoSqlLite) GetLikesCountForComment(commentID int) (int, error) {
	var count int
	err := r.db.QueryRow("SELECT likes FROM comments WHERE comment_id = ?", commentID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}

func (r *RepoSqlLite) InvalidateSessions(userID int) error {
	_, err := r.db.Exec("DELETE FROM session WHERE user_id = ?", userID)
	return err
}

// GetCommentsByUser retrieves comments left by a user.
func (r *RepoSqlLite) GetCommentsByUser(userID int) ([]domain.Comments, error) {
	comments := []domain.Comments{}
	query := "SELECT comment_id, post_id, user_id, content, username FROM comments WHERE user_id = ?"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var comment domain.Comments
		if err := rows.Scan(&comment.CommentId, &comment.PostId, &comment.UserId, &comment.Content, &comment.Username); err != nil {
			return nil, err
		}
		comments = append(comments, comment)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return comments, nil
}

// GetCreatedPosts retrieves posts created by a user.
func (r *RepoSqlLite) GetCreatedPosts(userID int) ([]domain.Posts, error) {
	posts := []domain.Posts{}
	query := "SELECT post_id, title, content, user_id, creation_date, category FROM posts WHERE user_id = ?"

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var post domain.Posts
		if err := rows.Scan(&post.PostId, &post.Title, &post.Content, &post.UserId, &post.CreationDate, &post.Category); err != nil {
			return nil, err
		}
		posts = append(posts, post)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return posts, nil
}

func (r *RepoSqlLite) CreateNotification(notification domain.Notification) error {
	query := "INSERT INTO notifications (user_id, activity, post_id, owner_id, creation_date,username ) VALUES (?, ?, ?, ?, ?,?)"
	_, err := r.db.Exec(query, notification.UserId, notification.Type, notification.PostId, notification.OwnerId, notification.Timestamp, notification.Username)
	return err
}

func (r *RepoSqlLite) CreateNotificationComments(notification domain.Notification_comments) error {
	query := "INSERT INTO notifications_comments (user_id, activity, comment_id, owner_id, creation_date,username, post_id) VALUES (?, ?, ?, ?, ?,?,?)"
	_, err := r.db.Exec(query, notification.UserId, notification.Type, notification.CommentId, notification.OwnerId, notification.Timestamp, notification.Username, notification.PostId)
	return err
}

// GetAllNotifications retrieves general notifications for an owner.
func (r *RepoSqlLite) GetAllNotifications(ownerID int) ([]domain.Notification, error) {
	notifications := []domain.Notification{}
	query := "SELECT id, user_id, activity, post_id, owner_id, creation_date,username FROM notifications WHERE owner_id = ?"
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification domain.Notification
		if err := rows.Scan(&notification.Id, &notification.UserId, &notification.Type, &notification.PostId, &notification.OwnerId, &notification.Timestamp, &notification.Username); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}

// GetAllNotificationsComment retrieves comment-specific notifications for an owner.
func (r *RepoSqlLite) GetAllNotificationsComment(ownerID int) ([]domain.Notification_comments, error) {
	notifications := []domain.Notification_comments{}
	query := "SELECT id, user_id, activity, comment_id, owner_id, creation_date,username , post_id FROM notifications_comments WHERE owner_id = ?"
	rows, err := r.db.Query(query, ownerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var notification domain.Notification_comments
		if err := rows.Scan(&notification.Id, &notification.UserId, &notification.Type, &notification.CommentId, &notification.OwnerId, &notification.Timestamp, &notification.Username, &notification.PostId); err != nil {
			return nil, err
		}
		notifications = append(notifications, notification)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return notifications, nil
}
