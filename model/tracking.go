package model

import (
	"log"
	"net"
	"net/http"

	"time"

	"github.com/verifiedninja/webapp/shared/database"
)

// *****************************************************************************
// Tracking
// *****************************************************************************

// User table contains the request tracking log
type Tracking_url struct {
	Id            uint32    `db:"id"`
	User_id       uint32    `db:"user_id"`
	Method        string    `db:"method"`
	Url           string    `db:"url"`
	RemoteAddress string    `db:"remote_address"`
	Referer       string    `db:"referer"`
	UserAgent     string    `db:"user_agent"`
	Created_at    time.Time `db:"created_at"`
	Updated_at    time.Time `db:"updated_at"`
}

// User table contains the request tracking log for API
type Tracking_API struct {
	Id             uint32    `db:"id"`
	User_id        uint32    `db:"user_id"`
	Method         string    `db:"method"`
	Url            string    `db:"url"`
	RemoteAddress  string    `db:"remote_address"`
	Referer        string    `db:"referer"`
	UserAgent      string    `db:"user_agent"`
	Lookup_user_id uint32    `db:"lookup_user_id"`
	Created_at     time.Time `db:"created_at"`
	Updated_at     time.Time `db:"updated_at"`
}

// Api_authentication
type Api_authentication struct {
	Id         uint32    `db:"id"`
	User_id    uint32    `db:"user_id"`
	Userkey    string    `db:"userkey"`
	Token      string    `db:"token"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// Walker tables contains usernames found and viewed
type Walker struct {
	Id         uint32    `db:"id"`
	User       string    `db:"user"`
	Site_id    uint32    `db:"site_id"`
	Viewed     uint8     `db:"viewed"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

func GetRemoteIP(r *http.Request) (string, error) {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if add := r.Header.Get("X-FORWARDED-FOR"); add != "" {
		ip = add
	} else if add := r.Header.Get("x-forwarded-for"); add != "" {
		ip = add
	} else if add := r.Header.Get("X-Forwarded-For"); add != "" {
		ip = add
	}

	return ip, err
}

func TrackRequestURL(user_id uint64, r *http.Request) error {
	ip, err := GetRemoteIP(r)
	if err != nil {
		log.Println(err)
	}

	_, err = database.DB.Exec("INSERT INTO tracking_url (user_id, method, url, remote_address, referer, user_agent) "+
		"VALUES (?,?,?,?,?,?)", user_id, r.Method, r.URL.RequestURI(), ip, r.Referer(), r.UserAgent())
	return err
}

func TrackRequestAPI(user_id uint64, r *http.Request, other_user_id uint64, verified bool) error {
	ip, err := GetRemoteIP(r)
	if err != nil {
		log.Println(err)
	}

	_, err = database.DB.Exec("INSERT INTO tracking_api (user_id, method, url, remote_address, referer, user_agent, lookup_user_id, verified) "+
		"VALUES (?,?,?,?,?,?,?,?)", user_id, r.Method, r.URL.Path, ip, r.Referer(), r.UserAgent(),
		other_user_id, verified)
	return err
}

func ApiAuthenticationByUserId(user_id uint64) (Api_authentication, error) {
	result := Api_authentication{}
	err := database.DB.Get(&result, "SELECT userkey, token FROM api_authentication WHERE user_id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}

func ApiAuthenticationCreate(user_id uint64, userkey, token string) error {
	_, err := database.DB.Exec("INSERT INTO api_authentication (user_id, userkey, token) VALUES (?,?,?)", user_id, userkey, token)
	return err
}

func ApiAuthenticationByKeys(userkey, token string) (Api_authentication, error) {
	result := Api_authentication{}
	err := database.DB.Get(&result, "SELECT user_id FROM api_authentication WHERE userkey = ? AND token = ? AND deleted = 0 LIMIT 1", userkey, token)
	return result, err
}
