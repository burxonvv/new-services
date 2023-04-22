package postgres

import (
	"log"
	"reflect"
	"testing"

	"github.com/new-york-services/user_service/storage/repo"
)

func TestUnit_User_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   repo.User
		want    repo.User
		wantErr bool
	}{
		{
			name: "success create user",
			input: repo.User{
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			want: repo.User{
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.CreateUser(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, tCase.wantErr, err)
			}
			tCase.want.Id = got.Id
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt
			tCase.want.Password = got.Password

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}

		})
	}

}

func TestUnit_User_Get(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    repo.User
		wantErr bool
	}{
		{
			name:  "success get user",
			input: "17",
			want: repo.User{
				Id:        "17",
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.GetUserById(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, tCase.wantErr, err)
			}
			tCase.want.Password = got.Password
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}

}

func TestUnit_User_Update(t *testing.T) {
	tests := []struct {
		name    string
		input   repo.User
		want    repo.User
		wantErr bool
	}{
		{
			name: "success update user",
			input: repo.User{
				FirstName: "New_name",
				LastName:  "new_last_name",
				Email:     "new_email",
				Id:        "9",
			},
			want: repo.User{
				Id:        "9",
				FirstName: "New_name",
				LastName:  "new_last_name",
				Email:     "new_email",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.UpdateUser(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, tCase.wantErr, err)
			}

			tCase.want.Password = got.Password
			tCase.want.Posts = got.Posts
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}
}

func TestUnit_User_Delete(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    repo.User
		wantErr bool
	}{
		{
			name:  "success delete user",
			input: "13",
			want: repo.User{
				Id:        "13",
				FirstName: "Justin",
				LastName:  "Bieber",
				Email:     "drew@gmail.com",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.DeleteUser(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, tCase.wantErr, got)
			}

			tCase.want.Password = got.Password
			tCase.want.Posts = got.Posts
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}

}
