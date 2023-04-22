package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/new-york-services/user_service/config"
	"github.com/new-york-services/user_service/pkg/db"
	"github.com/new-york-services/user_service/storage/postgres"
)

var pgRepo *postgres.UserRepo

func TestMain(m *testing.M) {
	conf := config.Load()
	connDb, err := db.ConnectToDB(conf)
	if err != nil {
		log.Println("failed to connect to database: ", err)
	}

	pgRepo = postgres.NewUserRepo(connDb)

	os.Exit(m.Run())
}
