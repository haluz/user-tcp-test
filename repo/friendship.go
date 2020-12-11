package repo

import (
	"encoding/json"
	"github.com/haluz/user-tcp-test/model"
)

type FriendshipRepo struct {
	friends map[int][]int
}

const friendshipData = `[
	{"user_id": 1, "friends": [2, 3, 4]},
	{"user_id": 2, "friends": [1]},
	{"user_id": 3, "friends": [1, 4]},
	{"user_id": 4, "friends": [1, 3]}
]`

func NewFriendshipRepo() (*FriendshipRepo, error) {
	var friendships []model.Friendship
	if err := json.Unmarshal([]byte(friendshipData), &friendships); err != nil {
		return nil, err
	}

	friends := make(map[int][]int)
	for _, v := range friendships {
		for _, ff := range v.Friends {
			friends[v.UserID] = append(friends[v.UserID], ff)
		}
	}

	return &FriendshipRepo{
		friends: friends,
	}, nil

}

func (r *FriendshipRepo) Friends(id int) ([]int, error) {
	return r.friends[id], nil
}
