package domain

import (
	"time"
)

// User 领域对象，是 DDD 中的 entity
// BO(business object)
type User struct {
	Id       int64
	Email    string
	Password string
	Nickname string
	Phone    string
	AboutMe  string
	Ctime    time.Time
	Birthday time.Time
}
