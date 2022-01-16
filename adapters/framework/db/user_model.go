package db

import (
	"fmt"
)

type (
	User struct {
		Id 			int 	`json:"id"`
		Username 	string 	`json:"username"`
		Email 		string 	`json:"email"`
		Password 	string 	`json:"password"`
	}
	
	UserModel struct {
		adapter		*DBAdapter
	}
)

func (adpt *DBAdapter) NewUserModel() *UserModel {
	return &UserModel{adapter: adpt}
}

// Define the methods of the UserModel
func (model *UserModel) Get(col string, value interface{}) (*User, error) {
	var user User
	query := fmt.Sprintf(`
		SELECT id, username, email, password FROM users
		WHERE %s = ?
	`, col)
	err := model.adapter.db.QueryRow(query, value).Scan(&user.Id, &user.Username, &user.Email, &user.Password)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (model *UserModel) Filter(col string, value interface{}) (*[]User, error) {
	var users []User
	query := fmt.Sprintf(`
		SELECT id, username, email, password FROM users
		WHERE %s = ?
	`, col)
	rows, err := model.adapter.db.Query(query, value)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			fmt.Println(err)
		}
		users = append(users, user)
	}
	return &users, nil
}

func (model *UserModel) All() (*[]User, error) {
	var users []User
	query := "SELECT id, username, email, password FROM users"
	rows, err := model.adapter.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var user User
		err := rows.Scan(&user.Id, &user.Username, &user.Email, &user.Password)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return &users, nil
}