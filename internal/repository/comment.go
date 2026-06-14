package repository

import (
	"context"
	"fmt"

	"github.com/dmytrii/youtube-gifs-chat/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

const (
	ActionLike    = 1
	ActionDislike = 2
)

type CommentRepository struct {
	pool *pgxpool.Pool
}

func NewCommentRepository(pool *pgxpool.Pool) *CommentRepository {
	return &CommentRepository{
		pool: pool,
	}
}

func (cr *CommentRepository) CreateComment(ctx context.Context, comment *domain.CreateCommentRequest, userId string) (*domain.Comment, error) {
	query := `INSERT INTO comments (user_id, gif_url, text, video_id, answer_to) VALUES ($1, $2, $3, $4, $5) RETURNING id`

	var commentId string

	err := cr.pool.QueryRow(ctx, query, userId, comment.GifUrl, comment.Text, comment.VideoId, comment.AnswerTo).Scan(&commentId)

	if err != nil {
		return nil, err
	}

	return cr.GetCommentById(ctx, commentId, userId)
}

func (cr *CommentRepository) DeleteComment(ctx context.Context, commentId string, userId string) error {
	query := `DELETE FROM comments WHERE id = $1 AND user_id = $2`

	_, err := cr.pool.Exec(ctx, query, commentId, userId)

	return err
}

func (cr *CommentRepository) ToggleAction(ctx context.Context, userId string, commentId string, actionType int) error {
	var existingId string
	checkQuery := `SELECT id FROM comment_actions WHERE user_id = $1 AND comment_id = $2 AND type = $3`
	err := cr.pool.QueryRow(ctx, checkQuery, userId, commentId, actionType).Scan(&existingId)

	if err == nil {
		_, err = cr.pool.Exec(ctx,
			`DELETE FROM comment_actions WHERE id = $1`, existingId)
		return err
	}

	_, _ = cr.pool.Exec(ctx,
		`DELETE FROM comment_actions WHERE user_id = $1 AND comment_id = $2`, userId, commentId)

	_, err = cr.pool.Exec(ctx,
		`INSERT INTO comment_actions (user_id, comment_id, type) VALUES ($1, $2, $3)`,
		userId, commentId, actionType)

	return err
}

func (cr *CommentRepository) LikeComment(ctx context.Context, userId string, commentId string) error {
	return cr.ToggleAction(ctx, userId, commentId, ActionLike)
}

func (cr *CommentRepository) DislikeComment(ctx context.Context, userId string, commentId string) error {
	return cr.ToggleAction(ctx, userId, commentId, ActionDislike)
}

func (cr *CommentRepository) GetCommentById(ctx context.Context, id string, userId string) (*domain.Comment, error) {
	query := `
		WITH RECURSIVE comment_tree AS (
			SELECT c.* FROM comments c WHERE c.id = $1
			UNION ALL
			SELECT c.* FROM comments c
			INNER JOIN comment_tree ct ON c.answer_to = ct.id
		)
		SELECT 
			c.id, c.user_id, c.video_id, c.gif_url, c.text, c.answer_to, c.created_at, c.position,
			u.id, u.name, u.picture, u.custom_url,
			COUNT(CASE WHEN ca.type = 1 THEN 1 END) AS likes,
			COUNT(CASE WHEN ca.type = 2 THEN 1 END) AS dislikes,
			COALESCE(BOOL_OR(ca.type = 1 AND ca.user_id = $2::uuid), false) AS liked_by_me,
			COALESCE(BOOL_OR(ca.type = 2 AND ca.user_id = $2::uuid), false) AS disliked_by_me
		FROM comment_tree c
		JOIN users u ON u.id = c.user_id
		LEFT JOIN comment_actions ca ON ca.comment_id = c.id
		GROUP BY c.id, c.user_id, c.video_id, c.gif_url, c.text, c.answer_to, c.created_at, c.position,
				u.id, u.name, u.picture, u.custom_url
		ORDER BY c.created_at ASC
	`

	rows, err := cr.pool.Query(ctx, query, id, userId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c, err := scanCommentFields(rows.Scan)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	if len(comments) == 0 {
		return nil, fmt.Errorf("comment not found")
	}

	roots := cr.buildTree(comments)
	return roots[0], nil
}

func (cr *CommentRepository) GetCommentByVideoId(ctx context.Context, videoId string, userId *string, cursor int, limit int) (*domain.Pagination, error) {
	query := `
        WITH RECURSIVE roots AS (
			SELECT id, position
			FROM comments
			WHERE video_id = $1 AND answer_to IS NULL AND position > $3
			ORDER BY position ASC
			LIMIT $4 + 1
		),
		tree AS (
			SELECT c.* FROM comments c
			JOIN roots r ON r.id = c.id
			UNION ALL
			SELECT c.* FROM comments c
			JOIN tree t ON c.answer_to = t.id
		)
		SELECT
			t.id, t.user_id, t.video_id, t.gif_url, t.text, t.answer_to, t.created_at, t.position,
			u.id AS commenter_id, u.name AS user_name, u.picture AS user_picture, u.custom_url AS user_custom_url,
			COUNT(CASE WHEN ca.type = 1 THEN 1 END) AS likes,
			COUNT(CASE WHEN ca.type = 2 THEN 1 END) AS dislikes,
			COALESCE(BOOL_OR(ca.type = 1 AND ca.user_id = $2::uuid), false) AS liked_by_me,
			COALESCE(BOOL_OR(ca.type = 2 AND ca.user_id = $2::uuid), false) AS disliked_by_me
		FROM tree t
		JOIN users u ON u.id = t.user_id
		LEFT JOIN comment_actions ca ON ca.comment_id = t.id
		GROUP BY t.id, t.user_id, t.video_id, t.gif_url, t.text, t.answer_to, t.created_at, t.position, u.id, u.name, u.picture, u.custom_url
		ORDER BY t.position ASC
    `

	allCommentsQuery := `
		SELECT COUNT(*) as comments_count FROM comments WHERE video_id = $1 AND answer_to IS NULL
	`

	var userIdParam interface{}
	if userId != nil {
		userIdParam = *userId
	}

	rows, err := cr.pool.Query(ctx, query, videoId, userIdParam, cursor, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*domain.Comment
	for rows.Next() {
		c, err := scanCommentFields(rows.Scan)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}

	comments = cr.buildTree(comments)

	commentsPagination := &domain.Pagination{
		Meta: &domain.Meta{
			HasCrusor:  false,
			NextCursor: cursor,
		},
	}

	commentsCountRow := cr.pool.QueryRow(ctx, allCommentsQuery, videoId)
	commentsCountRow.Scan(&commentsPagination.Meta.CommentsCount)

	if len(comments) > limit {
		commentsPagination.Meta.HasCrusor = true
		commentsPagination.Meta.NextCursor = comments[limit-1].Position
		comments = comments[:limit]
	}

	if comments == nil {
		commentsPagination.Entities = []*domain.Comment{}
	} else {
		commentsPagination.Entities = comments
	}

	return commentsPagination, nil

}

func (cr *CommentRepository) scanComment(row pgx.Row) (*domain.Comment, error) {
	c := &domain.Comment{}
	err := row.Scan(
		&c.ID, &c.UserId, &c.VideoId, &c.GifUrl, &c.Text, &c.AnswerTo, &c.CreatedAt, &c.Position,
		&c.User.ID, &c.User.Name, &c.User.Picture, &c.User.CustomUrl,
		&c.Likes, &c.Dislikes,
		&c.LikedByMe, &c.DislikedByMe,
	)
	return c, err
}

func (cr *CommentRepository) scanComments(rows pgx.Rows) ([]*domain.Comment, error) {
	var comments []*domain.Comment
	for rows.Next() {
		c, err := cr.scanComment(rows)
		if err != nil {
			return nil, err
		}
		comments = append(comments, c)
	}
	return comments, nil
}

func (cr *CommentRepository) buildTree(comments []*domain.Comment) []*domain.Comment {
	commentMap := make(map[string]*domain.Comment, len(comments))
	for _, c := range comments {
		c.Answers = []*domain.Comment{}
		commentMap[c.ID] = c
	}

	var roots []*domain.Comment
	for _, c := range comments {
		if c.AnswerTo == nil {
			roots = append(roots, c)
		} else {
			parent, ok := commentMap[*c.AnswerTo]
			if ok {
				parent.Answers = append(parent.Answers, c)
			} else {
				roots = append(roots, c)
			}
		}
	}
	return roots
}

func scanCommentFields(scan func(...any) error) (*domain.Comment, error) {
	c := &domain.Comment{}
	err := scan(
		&c.ID, &c.UserId, &c.VideoId, &c.GifUrl, &c.Text, &c.AnswerTo, &c.CreatedAt, &c.Position,
		&c.User.ID, &c.User.Name, &c.User.Picture, &c.User.CustomUrl,
		&c.Likes, &c.Dislikes,
		&c.LikedByMe, &c.DislikedByMe,
	)
	return c, err
}
