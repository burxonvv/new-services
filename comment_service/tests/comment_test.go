package postgres

import (
	"testing"

	"github.com/burxondv/new-services/comment-service/config"
	"github.com/burxondv/new-services/comment-service/pkg/db"
	"github.com/burxondv/new-services/comment-service/storage/postgres"
	"github.com/burxondv/new-services/comment-service/storage/repo"

	"github.com/stretchr/testify/suite"
)

type CommentSuiteTest struct {
	suite.Suite
	CleanUpfunc func()
	repo        repo.CommentStorageI
}

func (s *CommentSuiteTest) SetupSuite() {
	pgPool, cleanUpfunc := db.ConnectToDBForSuite(config.Load())
	s.repo = postgres.NewCommentRepo(pgPool)
	s.CleanUpfunc = cleanUpfunc
}

func (s *CommentSuiteTest) TestCommentCrud() {
	// comment := &c.CommentRequest{
	// 	PostId: 33,
	// 	UserId: 5,
	// 	Text:   "Woow, it's since the last year",
	// }

	// writeCommentResp, err := s.repo.WriteComment(comment)
	// s.Nil(err)
	// s.NotNil(writeCommentResp)
	// s.Equal(writeCommentResp.Text, comment.Text)
	// s.Equal(writeCommentResp.PostId, comment.PostId)
	// s.Equal(writeCommentResp.UserId, comment.UserId)

	// getCommentResp, err := s.repo.GetComments(&c.GetAllCommentsRequest{PostId: 33})
	// s.Nil(err)
	// s.NotNil(getCommentResp)

	// deleteCommentResp, err := s.repo.DeleteComment(&c.IdRequest{Id: writeCommentResp.Id})
	// s.Nil(err)
	// s.NotNil(deleteCommentResp)
	// s.Equal(deleteCommentResp.Text, writeCommentResp.Text)
	// s.Equal(deleteCommentResp.PostId, writeCommentResp.PostId)
	// s.Equal(deleteCommentResp.UserId, writeCommentResp.UserId)
}

func (suite *CommentSuiteTest) TearDownSuite() {
	suite.CleanUpfunc()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(CommentSuiteTest))
}
