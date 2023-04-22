package postgres

import (
	"fmt"

	"github.com/new-york-services/user_service/config"
	u "github.com/new-york-services/user_service/genproto/user"
	"github.com/new-york-services/user_service/pkg/db"
	"github.com/new-york-services/user_service/storage/postgres"
	"github.com/new-york-services/user_service/storage/repo"

	"testing"

	"github.com/stretchr/testify/suite"
)

type UserSuiteTest struct {
	suite.Suite
	CleanUpfunc func()
	repo        repo.UserStoreI
}

func (s *UserSuiteTest) SetupSuite() {
	pgPool, cleanUp := db.ConnectToDBForSuite(config.Load())
	s.repo = postgres.NewUserRepo(pgPool)
	s.CleanUpfunc = cleanUp
}

func (s *UserSuiteTest) TestUserCrud() {
	user := &u.UserResponse{
		FirstName: "Justin",
		LastName:  "Bieber",
		Email:     "drew@gmail.com",
	}

	createUserResp, err := s.repo.CreateUser(repo.User{FirstName: user.FirstName, LastName: user.LastName, Email: user.Email})
	fmt.Println(createUserResp)
	s.Nil(err)
	s.NotNil(createUserResp)
	s.Equal(user.FirstName, createUserResp.FirstName)
	s.Equal(user.LastName, createUserResp.LastName)

	getUserResp, err := s.repo.GetUserById(createUserResp.Id)
	s.Nil(err)
	s.NotNil(getUserResp)
	s.Equal(getUserResp.FirstName, user.FirstName)
	s.Equal(getUserResp.LastName, user.LastName)

	updateBody := repo.User{
		Id:        createUserResp.Id,
		FirstName: "Eminem",
		LastName:  "Habit",
		Email:     "eminem@gmail.com",
	}

	updateResp, err := s.repo.UpdateUser(updateBody)
	s.Nil(err)
	s.NotNil(updateResp)

	usersResp, err := s.repo.GetAllUsers(1, 100)
	s.Nil(err)
	s.NotNil(usersResp)

	deleteResp, err := s.repo.DeleteUser(createUserResp.Id)
	s.Nil(err)
	s.NotNil(deleteResp)

}

func (suite *UserSuiteTest) TearDownSuite() {
	suite.CleanUpfunc()
}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestUserRepositoryTestSuite(t *testing.T) {
	suite.Run(t, new(UserSuiteTest))
}
