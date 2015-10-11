package model

import (
	//"database/sql"
	"time"

	"github.com/verifiedninja/webapp/shared/database"
)

// *****************************************************************************
// Photo
// *****************************************************************************

// Photo table contains the information for each photo
type Photo struct {
	Id         uint32    `db:"id"`
	Path       string    `db:"path"`
	User_id    uint32    `db:"user_id"`
	Note       string    `db:"note"`
	Initial    uint8     `db:"initial"`
	Status_id  uint8     `db:"status_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

// User_status table contains every possible user status (active/inactive)
type Photo_status struct {
	Id         uint8     `db:"id"`
	Status     string    `db:"status"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

func PhotoCreate(user_id uint64, path string, initial bool) error {
	_, err := database.DB.Exec("INSERT INTO photo (path, user_id, initial) VALUES (?,?,?)", path, user_id, initial)
	return err
}

func PhotoDelete(user_id uint64, path string) error {
	_, err := database.DB.Exec("UPDATE photo SET deleted = 1 WHERE path = ? AND user_id = ? LIMIT 1", path, user_id)
	return err
}

func PhotosByUserId(user_id uint64) ([]Photo, error) {
	result := []Photo{}
	err := database.DB.Select(&result, "SELECT path, status_id, updated_at, note, initial FROM photo WHERE user_id = ? AND deleted = 0", user_id)
	return result, err
}

// Photo_info contains information necessary to verify access to the photo
type Photo_info struct {
	Owner_id   uint32 `db:"owner_id"`
	Role_level uint8  `db:"role_level"`
	Status_id  uint8  `db:"status_id"`
	Initial    uint8  `db:"initial"`
}

func PhotoInfoByPath(user_id uint64, path string) (Photo_info, error) {
	result := Photo_info{}
	err := database.DB.Get(&result, "SELECT photo.user_id as 'owner_id', photo.status_id as 'status_id', role.level_id as 'role_level', photo.initial as 'initial' FROM photo JOIN user ON user.id = photo.user_id JOIN role ON role.user_id = photo.user_id WHERE photo.path = ? AND photo.user_id = ? AND photo.deleted = 0 LIMIT 1", path, user_id)
	return result, err
}

func PhotoStatusByPath(user_id uint64, path string) uint8 {
	result := Photo{}
	database.DB.Get(&result, "SELECT status_id FROM photo WHERE path = ? AND user_id = ? AND deleted = 0 LIMIT 1", path, user_id)
	return result.Status_id
}

func PhotoApprove(path string, user_id uint64) error {
	_, err := database.DB.Exec("UPDATE photo SET status_id = 1 WHERE path = ? AND user_id = ? AND deleted = 0 LIMIT 1", path, user_id)
	if err != nil {
		return err
	}

	return nil
}

func PhotoReject(path string, user_id uint64, note string) error {
	_, err := database.DB.Exec("UPDATE photo SET status_id = 3, note = ? WHERE path = ? AND user_id = ? AND deleted = 0 LIMIT 1", note, path, user_id)
	if err != nil {
		return err
	}

	return nil
}

func PhotoUnverify(path string, user_id uint64) error {
	_, err := database.DB.Exec("UPDATE photo SET status_id = 2 WHERE path = ? AND user_id = ? AND deleted = 0 LIMIT 1", path, user_id)
	if err != nil {
		return err
	}

	return nil
}
