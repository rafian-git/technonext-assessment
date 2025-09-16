package model

import (
	"time"
)

type User struct {
	tableName struct{}  `pg:"users"`
	ID        int64     `pg:"id,pk"`
	Username  string    `pg:"username,unique,notnull"`
	Password  string    `pg:"password,notnull"` //hashed
	CreatedAt time.Time `pg:"created_at,default:now()"`
}
