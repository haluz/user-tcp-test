package server

type (
	loginRequest struct {
		UserID int `json:"user_id"`
	}

	userStatus struct {
		UserId int  `json:"user_id"`
		Online bool `json:"online"`
	}
)
