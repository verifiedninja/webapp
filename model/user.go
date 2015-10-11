package model

import (
	"database/sql"
	"log"
	"net/http"
	"time"

	"github.com/verifiedninja/webapp/shared/database"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	Id         uint32    `db:"id"`
	First_name string    `db:"first_name"`
	Last_name  string    `db:"last_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Status_id  uint8     `db:"status_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// User_status table contains every possible user status (active/inactive)
type User_status struct {
	Id         uint8     `db:"id"`
	Status     string    `db:"status"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// Email_verification table contains all verification codes for emails
type Email_verification struct {
	Id         uint32    `db:"id"`
	User_id    uint32    `db:"user_id"`
	Token      string    `db:"token"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// User_verification table contains random tokens for photos
type User_verification struct {
	Id         uint8     `db:"id"`
	Token      string    `db:"token"`
	User_id    uint32    `db:"user_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	result := User{}
	err := database.DB.Get(&result, "SELECT id, password, status_id, first_name, last_name FROM user WHERE email = ? AND deleted = 0 LIMIT 1", email)
	return result, err
}

// UserEmailByUserId gets user email from id
func UserEmailByUserId(user_id int64) (User, error) {
	result := User{}
	err := database.DB.Get(&result, "SELECT email, first_name FROM user WHERE id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}

// UserStatusByUserId gets user status from id
func UserStatusByUserId(user_id int64) (User, error) {
	result := User{}
	err := database.DB.Get(&result, "SELECT status_id FROM user WHERE id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}

// UserIdByEmail gets user id from email
func UserIdByEmail(email string) (User, error) {
	result := User{}
	err := database.DB.Get(&result, "SELECT id FROM user WHERE email = ? AND deleted = 0 LIMIT 1", email)
	return result, err
}

// UserByUsername gets user id from username
func UserByUsername(username string, site string) (Username_info, error) {
	result := Username_info{}
	err := database.DB.Get(&result, "SELECT username.name as 'username', username.user_id as 'id', site.name as 'site', site.profile as 'profile', site.url as 'home' FROM username JOIN site ON site.id = username.site_id WHERE username.name = ? AND site.name = ? AND username.deleted = 0 LIMIT 1", username, site)
	if err != nil {
		return result, err
	}
	return result, nil
}

// UserNameById gets user first_name and last_name from id
func UserNameById(id int) (User, error) {
	result := User{}
	err := database.DB.Get(&result, "SELECT first_name, last_name FROM user WHERE id = ? AND deleted = 0 LIMIT 1", id)
	return result, err
}

// UserCreate creates user
func UserCreate(first_name, last_name, email, password string) (sql.Result, error) {
	res, err := database.DB.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", first_name, last_name, email, password)
	return res, err
}

// UserEmailUpdate updates the user's email address
func UserEmailUpdate(user_id int64, email string) error {
	_, err := database.DB.Exec("UPDATE user SET email = ? WHERE id = ? AND deleted = 0", email, user_id)
	return err
}

// UserPasswordUpdate updates the current user's password
func UserPasswordUpdate(user_id int64, password string) error {
	_, err := database.DB.Exec("UPDATE user SET password = ? WHERE id = ? AND deleted = 0", password, user_id)
	return err
}

// UserEmailVerificationCreate adds a token to the database for email verification
func UserEmailVerificationCreate(user_id int64, hash string) error {
	_, err := database.DB.Exec("UPDATE email_verification SET deleted = 1 WHERE deleted = 0 AND user_id = ?", user_id)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("INSERT INTO email_verification (user_id, token) VALUES (?,?)", user_id, hash)
	return err
}

// UserUnverifiy updates the user so email verification is required to login again
func UserReverify(user_id int64) error {
	_, err := database.DB.Exec("UPDATE user SET status_id = 4 WHERE id = ? AND deleted = 0 LIMIT 1", user_id)
	return err
}

// UserEmailVerified gets if a token exists in the table, removes the token, and updates the user
func UserEmailVerified(token string) (bool, error) {
	result := Email_verification{}
	err := database.DB.Get(&result, "SELECT id, token, user_id FROM email_verification WHERE token = ? AND deleted = 0 LIMIT 1", token)

	if result.Token == token {
		_, err := database.DB.Exec("UPDATE email_verification SET deleted = 1 WHERE token = ? AND deleted = 0 LIMIT 1", token)
		if err != nil {
			return false, err
		} else {
			_, err := database.DB.Exec("UPDATE user SET status_id = 1 WHERE id = ? AND deleted = 0 LIMIT 1", result.User_id)
			return true, err
		}
	} else {
		return false, err
	}
}

func UserTokenCreate(user_id uint64, token string) error {
	_, err := database.DB.Exec("INSERT INTO user_verification (user_id, token) VALUES (?,?)", user_id, token)
	return err
}

func UserTokenByUserId(user_id uint64) (User_verification, error) {
	result := User_verification{}
	err := database.DB.Get(&result, "SELECT token FROM user_verification WHERE user_id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}

type EmailExpire struct {
	Id         uint32    `db:"id"`
	First_name string    `db:"first_name"`
	Email      string    `db:"email"`
	Expiring   bool      `db:"expiring"`
	Expired    bool      `db:"expired"`
	Updated_at time.Time `db:"updated_at"`
}

// UserIdByEmail gets user id from email
func EmailsWithVerificationIn30Days() ([]EmailExpire, error) {
	result := []EmailExpire{}
	err := database.DB.Select(&result, "SELECT user.id, user.first_name, user.email, "+
		"TIMESTAMPDIFF(DAY, MAX(email_verification.updated_at), NOW()) = 25 as 'expiring',"+
		"TIMESTAMPDIFF(DAY, MAX(email_verification.updated_at), NOW()) = 31 as 'expired',"+
		"MAX(email_verification.updated_at) as 'updated_at'"+
		"FROM user JOIN email_verification ON email_verification.user_id = user.id "+
		"WHERE email_verification.deleted = 1 AND user.deleted = 0 GROUP BY user.id"+
		" ORDER BY email_verification.updated_at DESC")
	return result, err
}

func EmailVerificationTokenByUserId(user_id uint64) (Email_verification, error) {
	result := Email_verification{}
	err := database.DB.Get(&result, "SELECT token FROM email_verification WHERE user_id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}

// UserLoginCreate logs the successful login and IP
func UserLoginCreate(user_id int64, r *http.Request) error {
	ip, err := GetRemoteIP(r)
	if err != nil {
		log.Println(err)
	}

	_, err = database.DB.Exec("UPDATE user_login SET deleted = 1 WHERE deleted = 0 AND user_id = ?", user_id)
	if err != nil {
		return err
	}
	_, err = database.DB.Exec("INSERT INTO user_login (user_id, remote_address) VALUES (?,?)", user_id, ip)
	return err
}
