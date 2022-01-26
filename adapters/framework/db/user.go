package db

import (
	"fmt"
	"log"
)

type (
	UserModel struct {
		Id 			int 	`json:"id"`
		Username 	string 	`json:"username"`
		Email 		string 	`json:"email"`
		Password 	string 	`json:"password"`
	}
	
	UserAdapter struct {
		adapter		*DBAdapter
		tableName	string
	}
)

func (adpt *DBAdapter) NewUserAdapter() *UserAdapter {
	return &UserAdapter{
		adapter: adpt,
		tableName: "users",
	}
}

// Define the methods of the UserAdapter
func (uAdpt *UserAdapter) Get(col string, value interface{}) (*UserModel, error) {
	var user UserModel
	query := fmt.Sprintf(`
		SELECT id, username, email, password FROM %s
		WHERE %s = $1
	`, uAdpt.tableName, col)
	err := uAdpt.adapter.db.QueryRow(query, value).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &user, nil
}

func (uAdpt *UserAdapter) Filter(col string, value interface{}) (*[]UserModel, error) {
	var users []UserModel
	query := fmt.Sprintf(`
		SELECT id, username, email, password FROM %s
		WHERE %s = $1
	`, uAdpt.tableName, col)
	rows, err := uAdpt.adapter.db.Query(query, value)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user UserModel
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}

func (uAdpt *UserAdapter) All() (*[]UserModel, error) {
	var users []UserModel
	query := "SELECT id, username, email, password FROM " + uAdpt.tableName
	rows, err := uAdpt.adapter.db.Query(query)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user UserModel
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			log.Fatal(err)
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}