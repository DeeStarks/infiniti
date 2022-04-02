package db

import (
	"fmt"
	"time"

	"github.com/deestarks/infiniti/utils"
	"github.com/lib/pq"
)

type (
	UserModel struct {
		Id 					int 					`json:"id"`
		FirstName 			string 					`json:"first_name"`
		LastName 			string 					`json:"last_name"`
		Email 				string 					`json:"email"`
		Password 			string 					`json:"password"`
		CreatedAt 			time.Time 				`json:"created_at"`
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
func (mAdapt *UserAdapter) Get(col string, value interface{}) (UserModel, error) {
	var user UserModel
	query := fmt.Sprintf(`
		SELECT id, first_name, last_name, email, password, created_at 
		FROM %s
		WHERE %s = $1
	`, mAdapt.tableName, col)

	err := mAdapt.adapter.db.QueryRow(query, value).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return UserModel{}, fmt.Errorf("%s", err.Detail)
    }
	return user, nil
}

// TODO: Populate the user model with the user's permissions and user groups
func (mAdapt *UserAdapter) Filter(col string, value interface{}) ([]UserModel, error) {
	var users []UserModel
	query := fmt.Sprintf(`
		SELECT id, first_name, last_name, email, password, created_at FROM %s
		WHERE %s = $1 ORDER BY id ASC
	`, mAdapt.tableName, col)
	rows, err := mAdapt.adapter.db.Query(query, value)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var user UserModel
		err := rows.Scan(
			&user.Id, &user.FirstName, &user.LastName,
			&user.Email, &user.Password, &user.CreatedAt,
		)
		if err, ok := err.(*pq.Error); ok {
			return nil, fmt.Errorf("%s", err.Detail)
		}
		users = append(users, user)
	}
	return users, nil
}

// TODO: Populate the user model with the user's permissions and user groups
func (mAdapt *UserAdapter) List() ([]UserModel, error) {
	var users []UserModel
	query := fmt.Sprintf(`
		SELECT id, first_name, last_name, email, password, created_at FROM %s
		ORDER BY id ASC
	`, mAdapt.tableName)
	rows, err := mAdapt.adapter.db.Query(query)
    if err, ok := err.(*pq.Error); ok {
		return nil, fmt.Errorf("%s", err.Detail)
    }
	defer rows.Close()

	for rows.Next() {
		var user UserModel
		err := rows.Scan(
			&user.Id, &user.FirstName, &user.LastName,
			&user.Email, &user.Password, &user.CreatedAt,
		)
		if err, ok := err.(*pq.Error); ok {
			return nil, fmt.Errorf("%s", err.Detail)
		}
		users = append(users, user)
	}
	return users, nil
}

func (mAdapt *UserAdapter) Create(data map[string]interface{}) (UserModel, error) {
	var user UserModel

	mToS := utils.MapToStructSlice(data)
	var (
		colStr		string
		valArr		[]interface{}
	)
	for i, s := range mToS {
		colStr += s.Key + ", "
		valArr = append(valArr, s.Value)
		if i == len(mToS)-1 {
			colStr = colStr[:len(colStr)-2] // remove the last ", "
		}
	}

	query := fmt.Sprintf(`
		INSERT INTO %s ( %s ) VALUES ( %s )
		RETURNING id, first_name, last_name, email, password, created_at
	`, mAdapt.tableName, colStr, utils.CreatePlaceholder(len(valArr)))

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return UserModel{}, fmt.Errorf("%s", err.Detail)
    }
	return user, nil
}

func (mAdapt *UserAdapter) Update(col string, colValue interface{}, data map[string]interface{}) (UserModel, error) {
	var (
		user 	UserModel
		valArr		[]interface{}
	)
	
	mToS := utils.MapToStructSlice(data)
	for _, s := range mToS {
		valArr = append(valArr, s.Value)
	}
	query := fmt.Sprintf(`
		UPDATE %s SET %s
		WHERE %s = $%d
		RETURNING id, first_name, last_name, email, password, created_at
	`, mAdapt.tableName, utils.CreateSetConditions(mToS), col, len(data)+1)

	valArr = append(valArr, colValue)

	err := mAdapt.adapter.db.QueryRow(query, valArr...).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return UserModel{}, fmt.Errorf("%s", err.Detail)
    }
	return user, nil
}

func (mAdapt *UserAdapter) Delete(colName string, value interface{}) (UserModel, error) {
	var (
		user	UserModel
		err			error
	)
	query := fmt.Sprintf(`
		DELETE FROM %s
		WHERE %s = $1
		RETURNING id, first_name, last_name, email, password, created_at
	`, mAdapt.tableName, colName)
	err = mAdapt.adapter.db.QueryRow(query, value).Scan(
		&user.Id, &user.FirstName, &user.LastName,
		&user.Email, &user.Password, &user.CreatedAt,
	)
    if err, ok := err.(*pq.Error); ok {
		return UserModel{}, fmt.Errorf("%s", err.Detail)
    }
	return user, nil
}