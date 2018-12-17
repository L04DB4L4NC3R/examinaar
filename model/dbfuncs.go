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
	Channel  uint16 `json:"channel"`
	Hosting  bool   `json:"hosting"`
}

// create sessions. Only 1 session can be created per host at one time

func CreateSessions(h HostType) (bool, error) {
	var flag uint8
	row := db.QueryRow(`
		SELECT HOSTING FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

	err := row.Scan(&flag)

	switch {
	case err == sql.ErrNoRows:
		break
	case err != nil:
		if flag == 1 {
			return false, fmt.Errorf("Session already in place")
		} else {
			return false, fmt.Errorf("Error in creating session")
		}
	}
	fmt.Println("ducky")
	_, err = db.Exec(`
		INSERT INTO HOSTS(PORT1,PORT2,IMAGE1,IMAGE2,CHANNEL,HOSTING)
		VALUES($2, $3, $4, $5, $6, $7)
		WHERE EMAIL=$1
	`, h.Email, h.Port1, h.Port2, h.Image1, h.Image2, h.Channel, h.Hosting)

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
		SELECT EMAIL,PORT1,PORT2,IMAGE1,IMAGE2,CHANNEL,HOSTING 
		FROM HOSTS WHERE HOSTING=$1
	`, 1)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		err = rows.Scan(&data.Email, &data.Port1, &data.Port2, &data.Image1, &data.Image2,
			&data.Channel, &data.Hosting)
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
	SELECT EMAIL,PORT1,PORT2,IMAGE1,IMAGE2,CHANNEL,HOSTING 
	FROM HOSTS WHERE EMAIL=$1
	`, h.Email)

	err := row.Scan(&data.Email, &data.Port1, &data.Port2, &data.Image1,
		&data.Image2, &data.Channel, &data.Hosting)

	if err != nil {
		return data, nil
	}

	return data, nil
}
