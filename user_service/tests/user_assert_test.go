package postgres

import (
	"errors"
	"reflect"
	"testing"

	"github.com/burxondv/new-services/user-service/storage/repo"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

var (
	InvalidArgumentError = errors.New("invalid argument")
)

func TestAssert_User_Create(t *testing.T) {
	ast := assert.New(t)
	tests := []struct {
		name string
		req  repo.User
		res  repo.User
		err  error
	}{
		{
			name: "error case invalid argument",
			req: repo.User{
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			res: repo.User{
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			err: InvalidArgumentError,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			id := uuid.New()
			tCase.req.Id = id.String()
			got, err := pgRepo.CreateUser(tCase.req)
			if err != nil && tCase.err != nil {
				ast.True(tCase.err.Error() == err.Error())
			}

			tCase.res.Id = got.Id
			tCase.res.CreatedAt = got.CreatedAt
			tCase.res.UpdatedAt = got.UpdatedAt
			tCase.res.Password = got.Password

			if !reflect.DeepEqual(tCase.res, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.res, got)
			}

		})
	}
}

func TestAssert_User_Get(t *testing.T) {
	ast := assert.New(t)
	tests := []struct {
		name string
		req  string
		res  repo.User
		err  error
	}{
		{
			name: "successfully get user",
			req:  "f40b017b-eee8-4fdb-b4d9-19ad211f5f53",
			res: repo.User{
				Id:        "f40b017b-eee8-4fdb-b4d9-19ad211f5f53",
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			err: nil,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			res, err := pgRepo.GetUserById(tCase.req)
			if tCase.err != nil && err != nil {
				ast.True(tCase.err.Error() == err.Error())
			}
			tCase.res.Password = res.Password
			tCase.res.Posts = res.Posts
			tCase.res.CreatedAt = res.CreatedAt
			tCase.res.UpdatedAt = res.UpdatedAt

			if !reflect.DeepEqual(tCase.res, res) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.res, res)
			}
		})
	}

}

func TestAssert_User_Update(t *testing.T) {
	ast := assert.New(t)
	tests := []struct {
		name string
		req  repo.User
		res  repo.User
		err  error
	}{
		{
			name: "successfully update user",
			req: repo.User{
				FirstName: "updated_first_name",
				LastName:  "updated_last_name",
				Email:     "updated_email",
				Id:        "3c69e7ca-0740-4d59-b8f1-351d03018830",
			},
			res: repo.User{
				Id:        "3c69e7ca-0740-4d59-b8f1-351d03018830",
				FirstName: "updated_first_name",
				LastName:  "updated_last_name",
				Email:     "updated_email",
			},
			err: nil,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.UpdateUser(tCase.req)
			if tCase.err != nil && err != nil {
				ast.True(tCase.err.Error() == err.Error())
			}

			tCase.res.Password = got.Password
			tCase.res.Posts = got.Posts
			tCase.res.CreatedAt = got.CreatedAt
			tCase.res.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.res, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.res, got)
			}
		})
	}
}

func TestUser_Delete(t *testing.T) {
	ast := assert.New(t)
	tests := []struct {
		name string
		req  string
		res  repo.User
		err  error
	}{
		{
			name: "successfully delete user",
			req:  "e7f11044-5fc4-4e22-bb2a-d5369edf2d8c",
			res: repo.User{
				Id:        "e7f11044-5fc4-4e22-bb2a-d5369edf2d8c",
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			err: nil,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			res, err := pgRepo.DeleteUser(tCase.req)
			if tCase.err != nil && err != nil {
				ast.True(tCase.err.Error() == err.Error())
			}

			tCase.res.Password = res.Password
			tCase.res.Posts = res.Posts
			tCase.res.CreatedAt = res.CreatedAt
			tCase.res.UpdatedAt = res.UpdatedAt

			if !reflect.DeepEqual(tCase.res, res) {
				t.Fatalf("%s: expected: %v, res: %v", tCase.name, tCase.res, res)
			}
		})
	}
}
