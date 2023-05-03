package postgres

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	u "github.com/burxondv/new-services/user-service/genproto/user"
	"github.com/burxondv/new-services/user-service/storage/repo"
)

func (r *UserRepo) CreateUser(user repo.User) (repo.User, error) {
	var res repo.User
	err := r.db.QueryRow(`
		insert into 
			users(id, first_name, last_name, email, password, refresh_token)
		values
			($1, $2, $3, $4, $5, $6)
		returning 
			id, first_name, last_name, user_type, email, refresh_token, created_at, updated_at`, user.Id, user.FirstName, user.LastName, user.Email, user.Password, user.RefreshToken).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to create user in sql: ", err)
		return repo.User{}, err
	}

	return res, nil
}

func (r *UserRepo) GetUserById(id string) (repo.User, error) {
	var res repo.User
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, user_type, email, refresh_token, created_at, updated_at
		from 
			users 
		where id = $1 and deleted_at is null`, id).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user by id in sql: ", err)
		return repo.User{}, err
	}

	return res, nil
}

func (r *UserRepo) GetUserForClient(user_id string) (repo.User, error) {
	var res repo.User
	err := r.db.QueryRow(`
		select 
			id, first_name, last_name, user_type, email, created_at, updated_at 
		from 
			users 
		where id = $1 and deleted_at is null`, user_id).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user for client in sql: ", err)
		return repo.User{}, err
	}

	return res, nil
}

func (r *UserRepo) GetUserByEmail(email string) (repo.User, error) {
	var res repo.User
	query := fmt.Sprint("select id, first_name, last_name, user_type, email, password, refresh_token, created_at, updated_at from users where email ilike '%" + email + "%' and deleted_at is null")

	err := r.db.QueryRow(query).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.Password, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to get user by email in sql: ", err)
		return repo.User{}, err
	}

	return res, nil
}

func (r *UserRepo) GetAllUsers(page, limit int64) ([]repo.User, error) {
	var res []repo.User
	offset := (page - 1) * limit
	rows, err := r.db.Query(`
		select 
			id, first_name, last_name, user_type, email, created_at, updated_at 
		from 
			users 
		where 
			deleted_at is null 
		limit $1 offset $2`, limit, offset)

	if err != nil {
		log.Println("failed to get all users in sql: ", err)
		return []repo.User{}, err
	}

	for rows.Next() {
		temp := repo.User{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.UserType,
			&temp.Email,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning all users in sql: ", err)
			return []repo.User{}, err
		}

		res = append(res, temp)
	}

	return res, nil
}

func (r *UserRepo) SearchUsers(req string) ([]repo.User, error) {
	var res []repo.User
	query := fmt.Sprint("select id, first_name, last_name, user_type, email, created_at, updated_at from users where first_name ilike '%" + req + "%' or last_name ilike '%" + req + "%' and deleted_at is null")

	rows, err := r.db.Query(query)
	if err != nil {
		log.Println("failed to searching user in sql: ", err)
		return []repo.User{}, err
	}

	for rows.Next() {
		temp := repo.User{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.UserType,
			&temp.Email,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scanning search user in sql: ", err)
			return []repo.User{}, err
		}

		res = append(res, temp)
	}

	return res, nil
}

func (r *UserRepo) UpdateUser(user repo.User) (repo.User, error) {
	res := repo.User{}
	err := r.db.QueryRow(`
		update
			users
		set
			first_name = $1, last_name = $2, email = $3, updated_at = $4
		where 
			id = $5 and deleted_at is null
		returning id, first_name, last_name, user_type, email, created_at, updated_at`, user.FirstName, user.LastName, user.Email, time.Now(), user.Id).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to update user in sql: ", err)
		return repo.User{}, err
	}

	fmt.Println("res -:>>>", res)
	return res, nil
}

func (r *UserRepo) DeleteUser(user_id string) (repo.User, error) {
	temp := repo.User{}
	err := r.db.QueryRow(`
		update 
			users
		set 
			deleted_at = $1 
		where 
			id = $2 and deleted_at is null
		returning 
			id, first_name, last_name, user_type, email, created_at, updated_at`, time.Now(), user_id).Scan(&temp.Id, &temp.FirstName, &temp.LastName, &temp.UserType, &temp.Email, &temp.CreatedAt, &temp.UpdatedAt)

	if err != nil {
		log.Println("failed to delete user in sql", err)
		return repo.User{}, err
	}

	return temp, nil
}

func (r *UserRepo) CheckField(req *u.CheckFieldRequest) (*u.CheckFieldResponse, error) {
	query := fmt.Sprintf("select 1 from users where %s=$1", req.Field)
	var temp int
	err := r.db.QueryRow(query, req.Value).Scan(&temp)
	if err == sql.ErrNoRows {
		return &u.CheckFieldResponse{Exists: false}, nil
	}

	if err != nil {
		return &u.CheckFieldResponse{}, err
	}

	if temp == 0 {
		return &u.CheckFieldResponse{Exists: true}, nil
	}

	return &u.CheckFieldResponse{Exists: false}, nil
}

func (r *UserRepo) UpdateUserTokens(req *u.UpdateUserTokensRequest) (*u.UserResponse, error) {
	res := u.UserResponse{}
	err := r.db.QueryRow(`
		update
			users
		set 
			refresh_token = $1
		where
			id = $2 and deleted_at is null
		returning id, first_name, last_name, user_type, email, refresh_token, created_at, updated_at`, req.RefreshToken, req.Id).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to update tokens in sql: ", err)
		return nil, err
	}

	return &res, nil
}

func (r *UserRepo) ChangeRoleUser(req *u.ChangeRoleRequest) (*u.UserResponse, error) {
	res := u.UserResponse{}
	err := r.db.QueryRow(`
		update
			users
		set
			user_type = $1, updated_at = $2
		where
			id = $3
		returning id, first_name, last_name, user_type, email, refresh_token, created_at, updated_at`, req.Role, time.Now(), req.Id).Scan(&res.Id, &res.FirstName, &res.LastName, &res.UserType, &res.Email, &res.RefreshToken, &res.CreatedAt, &res.UpdatedAt)

	if err != nil {
		log.Println("failed to update user type in sql: ", err)
		return nil, err
	}

	return &res, nil
}

func (r *UserRepo) GetSameRoleUsers(str string) (*u.UsersResponse, error) {
	users := u.UsersResponse{}
	rows, err := r.db.Query(`
		select
			id, first_name, last_name, user_type, email, refresh_token, created_at, updated_at
		from
			users
		where
			user_type = $1 and deleted_at is null`, str)

	if err != nil {
		log.Println("failed to get the same role users in sql: ", err)
		return nil, err
	}

	for rows.Next() {
		temp := u.UserResponse{}

		err = rows.Scan(
			&temp.Id,
			&temp.FirstName,
			&temp.LastName,
			&temp.UserType,
			&temp.Email,
			&temp.RefreshToken,
			&temp.CreatedAt,
			&temp.UpdatedAt,
		)
		if err != nil {
			log.Println("failed to scan the same role users in sql: ", err)
			return nil, err
		}

		users.Users = append(users.Users, &temp)
	}

	return &users, nil
}
