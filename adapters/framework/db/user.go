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
func (userAdpt *UserAdapter) Get(col string, value interface{}) (*UserModel, error) {
	var user UserModel
	query := fmt.Sprintf(`
		SELECT id, username, email, password FROM %s
		WHERE %s = $1
	`, userAdpt.tableName, col)
	err := userAdpt.adapter.db.QueryRow(query, value).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	return &user, nil
}

func (userAdpt *UserAdapter) Filter(col string, value interface{}) (*[]UserModel, error) {
	var users []UserModel
	query := fmt.Sprintf(`
		SELECT id, username, email, password FROM %s
		WHERE %s = $1
	`, userAdpt.tableName, col)
	rows, err := userAdpt.adapter.db.Query(query, value)
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

func (userAdpt *UserAdapter) All() (*[]UserModel, error) {
	var users []UserModel
	query := "SELECT id, username, email, password FROM " + userAdpt.tableName
	rows, err := userAdpt.adapter.db.Query(query)
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