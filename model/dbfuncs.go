package model

import (
	"database/sql"
	"fmt"
	"log"
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

func CreateSessions(h HostType) (bool, error) {
	var mail string
	row := db.QueryRow(`
		SELECT EMAIL FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

	err := row.Scan(&mail)

	switch {
	case err == sql.ErrNoRows:
		return false, fmt.Errorf("User not found")
	case err != nil:
		return false, err
	default:
		break

	}

	_, err = db.Exec(`
		UPDATE HOSTS
		SET HOSTING=$1, PORT1=$2, PORT2=$3, IMAGE1=$4, IMAGE2=$5
		WHERE EMAIL=$6
	`, 1, h.Port1, h.Port2, h.Image1, h.Image2, h.Email)

	if err != nil {
		return false, err
	}

	return true, nil

}

// Show all active sessions

func ReadSessions() ([]HostType, error) {
	var (
		arr  []HostType
		data HostType
	)

	rows, err := db.Query(`
		SELECT EMAIL, PORT1, PORT2, IMAGE1, IMAGE2, HOSTING 
		FROM HOSTS WHERE HOSTING=$1
	`, 1)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&data.Email, &data.Port1, &data.Port2, &data.Image1, &data.Image2, &data.Hosting)
		switch {
		case err == sql.ErrNoRows:
			return nil, nil
		case err != nil:
			return nil, err
		default:
			arr = append(arr, data)

		}

	}
	log.Println(arr)
	return arr, nil
}

func DeleteSessions(h HostType) (bool, error) {
	_, err := db.Exec(`
		DELETE FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

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
