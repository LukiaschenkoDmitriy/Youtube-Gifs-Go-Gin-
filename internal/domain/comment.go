package domain

import "time"

type Comment struct {
	ID           string     `json:"id" db:"id"`
	UserId       string     `json:"user_id" db:"user_id"`
	VideoId      string     `json:"video_id" db:"video_id"`
	GifUrl       string     `json:"gif_url" db:"gif_url"`
	Text         string     `json:"text" db:"text"`
	Likes        int        `json:"likes" db:"likes"`
	Dislikes     int        `json:"dislikes" db:"dislikes"`
	AnswerTo     *string    `json:"answer_to" db:"answer_to"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	LikedByMe    bool       `json:"liked_by_me" db:"-"`
	DislikedByMe bool       `json:"disliked_by_me" db:"-"`
	User         User       `json:"user" db:"-"`
	Answers      []*Comment `json:"answers" db:"-"`
}

type CreateCommentRequest struct {
	VideoId  string  `json:"video_id" binding:"required" db:"video_id"`
	GifUrl   string  `json:"gif_url" db:"gif_url"`
	Text     string  `json:"text" db:"text"`
	AnswerTo *string `json:"answer_to" binding:"omitempty,uuid" db:"answer_to"`
}
