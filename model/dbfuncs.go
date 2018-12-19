package model

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func SetDB(database *sql.DB) {
	db = database
}

type HostType struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Port1    string `json:"port1"`
	Port2    string `json:"port2"`
	Image1   string `json:"image1"`
	Image2   string `json:"image2"`
	Channel  string `json:"channel"`
	Hosting  bool   `json:"hosting"`
}

// create sessions. Only 1 session can be created per host at one time

func CreateSessions(h HostType) (bool, string) {
	var temp string
	roww := db.QueryRow(`
		SELECT EMAIL
		FROM HOSTS
		WHERE PORT1=$1 OR PORT2=$2
	`, h.Port1, h.Port2)

	err := roww.Scan(&temp)

	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return false, "Some error occurrened"
	default:
		if len(temp) > 0 {
			return false, "Some session is already going on with these ports, please choose different ports"
		}
	}

	var mail string
	row := db.QueryRow(`
		SELECT EMAIL FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

	err = row.Scan(&mail)

	switch {
	case err == sql.ErrNoRows:
		return false, "User not found"
	case err != nil:
		return false, "Some error occurred while scanning"
	default:
		break

	}

	_, err = db.Exec(`
		UPDATE HOSTS
		SET HOSTING=$1, PORT1=$2, PORT2=$3, IMAGE1=$4, IMAGE2=$5, CHANNEL=$6
		WHERE EMAIL=$7
	`, 1, h.Port1, h.Port2, h.Image1, h.Image2, h.Channel, h.Email)

	if err != nil {
		return false, "Some error occurred while executing"
	}

	return true, ""

}

// Show all active sessions

func ReadSessions() ([]HostType, error) {
	var (
		arr  []HostType
		data HostType
	)

	rows, err := db.Query(`
		SELECT EMAIL, PORT1, PORT2, IMAGE1, IMAGE2, HOSTING, CHANNEL 
		FROM HOSTS WHERE HOSTING=$1
	`, 1)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&data.Email, &data.Port1, &data.Port2, &data.Image1, &data.Image2, &data.Hosting, &data.Channel)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			arr = append(arr, data)

		}

	}
	return arr, nil
}

func DeleteSessions(e string) (bool, error) {
	_, err := db.Exec(`
		UPDATE HOSTS 
		SET PORT1='', PORT2='', IMAGE1='', IMAGE2='', CHANNEL='', HOSTING=$1 
		WHERE EMAIL=$2
	`, 0, e)

	if err != nil {
		return false, err
	}

	return true, nil
}

func CreateHost(h HostType) (bool, error) {
	var checker string
	row := db.QueryRow(`
		SELECT EMAIL FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

	err := row.Scan(&checker)

	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		return false, err
	case err == nil:
		if len(checker) > 0 {
			return false, fmt.Errorf("User already exists")
		}
	}

	_, err = db.Exec(`
		INSERT INTO HOSTS(EMAIL, PASSWORD)
		VALUES($1, $2)
	`, h.Email, h.Password)

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetHost(h HostType) (HostType, error) {

	var data HostType

	row := db.QueryRow(`
	SELECT EMAIL,PASSWORD,PORT1,PORT2,IMAGE1,IMAGE2,CHANNEL,HOSTING
	FROM HOSTS WHERE EMAIL=$1
	`, h.Email)

	err := row.Scan(&data.Email, &data.Password, &data.Port1, &data.Port2, &data.Image1,
		&data.Image2, &data.Channel, &data.Hosting)

	if err != nil {
		return data, nil
	}

	return data, nil
}
