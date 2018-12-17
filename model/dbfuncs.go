package model

import (
	"database/sql"
	"fmt"
)

var db *sql.DB

func SetDB(db *sql.DB) {
	db = db
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
			break
		}
	}

	_, err = db.Exec(`
		INSERT INTO HOSTS(EMAIL,PASSWORD,PORT1,PORT2,IMAGE1,IMAGE2,CHANNEL,HOSTING)
		VALUES($1, $2, $3, $4, $5, $6, $7, $8)
	`, h.Email, h.Password, h.Port1, h.Port2, h.Image1, h.Image2, h.Channel, h.Hosting)

	if err != nil {
		return false, err
	}

	return true, nil

}

// Show all active sessions

func ReadSessions(h *HostType) ([]HostType, error) {
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

func DeleteSessions(h *HostType) (bool, error) {
	_, err := db.Exec(`
		DELETE * FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

	if err != nil {
		return false, err
	}

	return true, nil
}

func StopSession(h *HostType) (bool, error) {
	_, err := db.Exec(`
		DELETE PORT1, PORT2, IMAGE1, IMAGE2, HOSTING
		FROM HOSTS
		WHERE EMAIL=$1
	`, h.Email)

	if err != nil {
		return false, err
	}

	return true, nil
}
