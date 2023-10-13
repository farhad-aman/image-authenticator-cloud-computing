package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

var db *sql.DB

type User struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	National string `json:"national"`
	IP       string `json:"ip"`
	State    string `json:"state"`
}

func InitPG() error {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return fmt.Errorf("failed to connect to the database: %v", err)
	}

	err = db.Ping()
	if err != nil {
		return fmt.Errorf("failed to ping the database: %v", err)
	}

	return nil
}

func GetUserData(nationalID string) (*User, error) {
	query := "SELECT name, email, national, ip, state FROM users WHERE national = $1"

	var user User
	err := db.QueryRow(query, EncodeNationalID(nationalID)).Scan(&user.Name, &user.Email, &user.National, &user.IP, &user.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user not found")
		}
		return nil, err
	}
	user.National, err = DecodeNationalID(user.National)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func SaveUserData(name, email, nationalID, ip string) error {
	hashedNationalID := EncodeNationalID(nationalID)
	_, err := db.Exec("INSERT INTO users (name, email, national, ip, state) VALUES ($1, $2, $3, $4, $5)", name, email, hashedNationalID, ip, "pending")
	return err
}
