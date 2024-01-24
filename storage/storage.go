package storage

import (
	"database/sql"
	"gorillaMux/models"

	_ "github.com/lib/pq"
)

// this function connected database and return sql.DB and error
func connect() (*sql.DB, error) {
	dsn := "user=newuser password=1234 dbname=newdb sslmode=disable"

	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic(err)
	}

	return db, nil
}

// this function creates new user
func CreateUser(user *models.User) (*models.User, error) {
	db, err := connect()
	if err != nil {
		return &models.User{}, err
	}

	defer db.Close()

	query := `INSERT INTO users (id, name, last_name) VALUES($1, $2, $3) 
  	RETURNING id, name, last_name`

	var respUser models.User
	if err = db.QueryRow(query,
		user.Id,
		user.FirstName,
		user.LastName).Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName); err != nil {
		return &models.User{}, err
	}

	return &respUser, nil
}

// this function updates user.FirstName and user.LastName with userId
func UpdateUser(userId string, user *models.User) (*models.User, error) {
	db, err := connect()
	if err != nil {
		return &models.User{}, err
	}
	defer db.Close()

	query := `
  	UPDATE 
    	users 
  	SET 
    	name = $1, 
    	last_name = $2
  	WHERE 
    	id = $3
  	RETURNING 
    	id, 
    	name, 
    	last_name`

	var respUser models.User
	if err := db.QueryRow(query, user.FirstName, user.LastName, userId).Scan(
		&respUser.Id,
		&respUser.FirstName,
		&respUser.LastName,
	); err != nil {
		return &models.User{}, err
	}

	return &respUser, nil
}

// this function deletes user with userId
func DeleteUser(userId string) error {
	db, err := connect()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM users WHERE id = $1`

	_, err = db.Exec(query, userId)
	if err != nil {
		return err
	}
	return nil
}

// this function gets user with userId
func GetUser(userId string) (*models.User, error) {
	db, err := connect()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	query := `SELECT id, name, last_name FROM users WHERE id = $1`

	var respUser models.User
	if err = db.QueryRow(query, userId).Scan(&respUser.Id, &respUser.FirstName, &respUser.LastName); err != nil {
		return nil, err
	}

	return &respUser, nil
}

// this func gets page, limit and return users with offset
func GetAllUsers(page, limit int) ([]*models.User, error) {
	db, err := connect()
	if err != nil {
		return []*models.User{}, err
	}
	defer db.Close()

	var users []*models.User
	offset := limit * (page - 1)
	query := `SELECT id, name, last_name FROM users LIMIT $1 OFFSET $2`

	rows, err := db.Query(query, limit, offset)
	if err != nil {
		return []*models.User{}, err
	}

	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.Id, &user.FirstName, &user.LastName); err != nil {
			return []*models.User{}, err
		}
		users = append(users, &user)
	}
	return users, nil
}
