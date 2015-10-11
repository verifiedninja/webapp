package model

import (
	//"database/sql"
	"time"

	"github.com/verifiedninja/webapp/shared/database"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type Demographic struct {
	Id            uint32    `db:"id"`
	User_id       uint32    `db:"user_id"`
	Birth_month   uint8     `db:"birth_month"`
	Birth_day     uint8     `db:"birth_day"`
	Birth_year    uint16    `db:"birth_year"`
	Gender        string    `db:"gender"`
	Height_feet   uint8     `db:"height_feet"`
	Height_inches uint8     `db:"height_inches"`
	Weight        uint16    `db:"weight"`
	Created_at    time.Time `db:"created_at"`
	Updated_at    time.Time `db:"updated_at"`
	Deleted       uint8     `db:"deleted"`
}

//
type Ethnicity_type struct {
	Id         uint8     `db:"id"`
	Name       string    `db:"name"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

//
type Ethnicity struct {
	Id         uint32    `db:"id"`
	User_id    uint32    `db:"user_id"`
	Type_id    uint8     `db:"type_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    uint8     `db:"deleted"`
}

func DemographicByUserId(user_id uint64) (Demographic, error) {
	result := Demographic{}
	err := database.DB.Get(&result, "SELECT * FROM demographic WHERE user_id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}

func DemographicAdd(user_id uint64, d Demographic) error {
	_, err := database.DB.Exec("INSERT INTO demographic (user_id, birth_month, birth_day,"+
		" birth_year, gender, height_feet, height_inches, weight) VALUES (?,?,?,?,?,?,?,?) "+
		"ON DUPLICATE KEY UPDATE birth_month = ?, birth_day = ?,"+
		" birth_year = ?, gender = ?, height_feet = ?, height_inches = ?, weight = ?", user_id,
		d.Birth_month, d.Birth_day, d.Birth_year,
		d.Gender, d.Height_feet, d.Height_inches, d.Weight, d.Birth_month, d.Birth_day, d.Birth_year,
		d.Gender, d.Height_feet, d.Height_inches, d.Weight)

	if err != nil {
		return err
	}

	return nil
}

func EthnicityAdd(user_id uint64, e []string) error {
	_, err := database.DB.Exec("UPDATE ethnicity SET deleted = 1 WHERE user_id = ? AND deleted = 0", user_id)
	if err != nil {
		return err
	}

	for _, v := range e {
		_, err := database.DB.Exec("INSERT INTO ethnicity (user_id, type_id) VALUES (?,?)", user_id, v)
		if err != nil {
			return err
		}
	}

	return nil
}

func EthnicityByUserId(user_id uint64) ([]Ethnicity, error) {
	result := []Ethnicity{}
	err := database.DB.Select(&result, "SELECT * FROM ethnicity WHERE user_id = ? AND deleted = 0 LIMIT 1", user_id)
	return result, err
}
