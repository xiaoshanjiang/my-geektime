package domain

import (
	"time"
)

// User 领域对象，是 DDD 中的 entity
// BO(business object)
type User struct {
	Id        int64
	Email     string
	Password  string
	Nickname  string
	Biography string
	Birthday  time.Time
	Ctime     time.Time
}

type UserEdit struct {
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Biography string    `json:"biography"`
	Birthday  time.Time `json:"birthday" time_format:"2006-01-02"`
}

type UserRead struct {
	Id        int64     `json:"id"`
	Email     string    `json:"email"`
	Nickname  string    `json:"nickname"`
	Biography string    `json:"biography"`
	Birthday  time.Time `json:"birthday" time_format:"2006-01-02"`
}
