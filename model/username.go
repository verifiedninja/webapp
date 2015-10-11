package model

import (
	"database/sql"
	//"database/sql"
	"time"

	"github.com/verifiedninja/webapp/shared/database"
)

// *****************************************************************************
// Username
// *****************************************************************************

// Username table contains the usernames for each user
type Username struct {
	Id         uint32    `db:"id"`
	Name       string    `db:"name"`
	User_id    uint32    `db:"user_id"`
	Site_id    uint32    `db:"site_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// Site table contains every possible website were usernames can be verified
type Site struct {
	Id         uint32    `db:"id"`
	Name       string    `db:"name"`
	Url        string    `db:"url"`
	Profile    string    `db:"profile"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

func SiteList() ([]Site, error) {
	result := []Site{}
	err := database.DB.Select(&result, "SELECT id, name FROM site WHERE deleted = 0")
	return result, err
}

func UsernameRemove(user_id uint64, site_id uint64) (sql.Result, error) {
	result, err := database.DB.Exec("UPDATE username SET deleted = 1 WHERE user_id = ? AND site_id = ? AND deleted = 0 LIMIT 1", user_id, site_id)

	if err != nil {
		return nil, err
	}

	return result, nil
}

func UsernameAdd(user_id uint64, username string, site_id uint64) error {
	_, err := database.DB.Exec("INSERT INTO username (user_id, name, site_id) VALUES (?,?,?) ON DUPLICATE KEY UPDATE name = ?", user_id, username, site_id, username)

	if err != nil {
		return err
	}

	return nil
}

func UsernamesByUserId(user_id uint64) ([]Username, error) {
	result := []Username{}
	err := database.DB.Select(&result, "SELECT name, site_id FROM username WHERE deleted = 0 AND user_id = ?", user_id)
	return result, err
}

type Username_info struct {
	Id       uint32 `db:"id"`
	Username string `db:"username"`
	Site     string `db:"site"`
	Profile  string `db:"profile"`
	Home     string `db:"home"`
}

func UserinfoByUserId(user_id uint64) ([]Username_info, error) {
	result := []Username_info{}
	err := database.DB.Select(&result, "SELECT username.name as 'username', site.name as 'site', site.profile as 'profile', site.url as 'home' FROM username JOIN site ON site.id = username.site_id WHERE username.user_id = ? AND username.deleted = 0", user_id)
	return result, err
}
