package model

import (
	"database/sql"
	"time"

	"github.com/verifiedninja/webapp/shared/database"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the role of each user
type Role struct {
	Id         uint32    `db:"id"`
	User_id    uint32    `db:"user_id"`
	Level_id   uint8     `db:"level_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// Role_level table contains levels of access (Administrator, User, etc)
type Role_level struct {
	Id         uint8     `db:"id"`
	Name       string    `db:"name"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

const (
	Role_level_Administrator = uint8(1)
	Role_level_User          = uint8(2)
)

// RoleCreate creates user with a role
func RoleCreate(user_id int64, level_id uint8) (sql.Result, error) {
	res, err := database.DB.Exec("INSERT INTO role (user_id, level_id) VALUES (?,?)", user_id, level_id)
	return res, err
}

// RoleByUserId gets the role by user_id
func RoleByUserId(user_id int64) (Role, error) {
	result := Role{}
	err := database.DB.Get(&result, "SELECT level_id FROM role WHERE user_id = ? LIMIT 1", user_id)
	return result, err
}
