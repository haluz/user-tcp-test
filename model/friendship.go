package model

type Friendship struct {
	UserID  int   `json:"user_id"`
	Friends []int `json:"friends"`
}
