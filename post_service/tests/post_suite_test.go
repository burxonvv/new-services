package tests

import (
	"testing"

	"github.com/burxondv/new-services/post-service/config"
	"github.com/burxondv/new-services/post-service/pkg/db"
	"github.com/burxondv/new-services/post-service/storage/postgres"
	"github.com/burxondv/new-services/post-service/storage/repo"

	"github.com/stretchr/testify/suite"
)

type PostSuiteTest struct {
	suite.Suite
	CleanUpfunc func()
	repo        repo.PostStorageI
}

func (s *PostSuiteTest) SetupSuite() {
	pgPool, cleanUp := db.ConnectToDBForSuite(config.Load())
	s.repo = postgres.NewPostRepo(pgPool)
	s.CleanUpfunc = cleanUp
}

func (s *PostSuiteTest) TestPostCrud() {
	post := repo.Post{
		Title:       "My fifty sixth post",
		Description: "I love my job",
		UserId:      "5",
	}

	createPostResp, err := s.repo.CreatePost(post)
	s.Nil(err)
	s.NotNil(createPostResp)
	s.Equal(post.Title, createPostResp.Title)
	s.Equal(post.Description, createPostResp.Description)
	s.Equal(post.UserId, createPostResp.UserId)

	getPostResp, err := s.repo.GetPostById(createPostResp.Id)
	s.Nil(err)
	s.NotNil(getPostResp)
	s.Equal(getPostResp.Title, post.Title)
	s.Equal(getPostResp.Description, post.Description)
	s.Equal(getPostResp.UserId, post.UserId)

	allPostsResp, err := s.repo.GetPostByUserId(getPostResp.UserId)
	s.Nil(err)
	s.NotNil(allPostsResp)

	updateBody := repo.Post{
		Title:       "Seventeenth post",
		Description: "It's jog, my seventeenth post",
		Id:          getPostResp.Id,
	}

	updateRes, err := s.repo.UpdatePost(updateBody)
	s.Nil(err)
	s.NotNil(updateRes)

	deleteResp, err := s.repo.DeletePost(getPostResp.Id)
	s.Nil(err)
	s.NotNil(deleteResp)
	s.Equal(deleteResp.Title, updateBody.Title)
	s.Equal(deleteResp.Description, updateBody.Description)
}

func (suite *PostSuiteTest) TearDownSuite() {
	suite.CleanUpfunc()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(PostSuiteTest))
}
