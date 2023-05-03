package tests

import (
	"log"
	"reflect"
	"testing"

	"github.com/burxondv/new-services/post-service/storage/repo"
)

func TestPost_Create(t *testing.T) {
	tests := []struct {
		name    string
		input   repo.Post
		want    repo.Post
		wantErr bool
	}{
		{
			name: "success create post",
			input: repo.Post{
				Title:       "coding",
				Description: "I was start 18",
				UserId:      "7",
			},
			want: repo.Post{
				Title:       "coding",
				Description: "I was start 18",
				UserId:      "7",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.CreatePost(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, tCase.wantErr, err)
			}
			tCase.want.Id = got.Id
			tCase.want.Likes = got.Likes
			tCase.want.Comments = got.Comments
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}
}

func TestPost_Get(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    repo.Post
		wantErr bool
	}{
		{
			name:  "success get user",
			input: "30",
			want: repo.Post{
				Id:          "30",
				Title:       "how",
				Description: "i don't know",
				UserId:      "9",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.GetPostById(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, tCase.wantErr, err)
			}
			tCase.want.Likes = got.Likes
			tCase.want.Comments = got.Comments
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}
}

func TestPost_Update(t *testing.T) {
	tests := []struct {
		name    string
		input   repo.Post
		want    repo.Post
		wantErr bool
	}{
		{
			name: "success update post",
			input: repo.Post{
				Title:       "new_updated_title",
				Description: "new_updated_description",
				Id:          "31",
			},
			want: repo.Post{
				Title:       "new_updated_title",
				Description: "new_updated_description",
				Id:          "31",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.UpdatePost(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, err: %v", tCase.name, tCase.wantErr, err)
			}

			tCase.want.Likes = got.Likes
			tCase.want.Comments = got.Comments
			tCase.want.UserId = got.UserId
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s: expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}
}

func TestPost_Delete(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    repo.Post
		wantErr bool
	}{
		{
			name:  "success delete post",
			input: "31",
			want: repo.Post{
				Title:       "new_updated_title",
				Description: "new_updated_description",
				Id:          "31",
				UserId:      "9",
			},
			wantErr: false,
		},
	}

	for _, tCase := range tests {
		t.Run(tCase.name, func(t *testing.T) {
			got, err := pgRepo.DeletePost(tCase.input)
			if err != nil {
				log.Printf("%s: expected: %v, got: %v", tCase.name, err, got)
			}

			tCase.want.Likes = got.Likes
			tCase.want.Comments = got.Comments
			tCase.want.CreatedAt = got.CreatedAt
			tCase.want.UpdatedAt = got.UpdatedAt

			if !reflect.DeepEqual(tCase.want, got) {
				t.Fatalf("%s, expected: %v, got: %v", tCase.name, tCase.want, got)
			}
		})
	}

}
