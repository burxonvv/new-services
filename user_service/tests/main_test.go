package postgres

import (
	"log"
	"os"
	"testing"

	"github.com/burxondv/new-services/user-service/config"
	"github.com/burxondv/new-services/user-service/pkg/db"
	"github.com/burxondv/new-services/user-service/storage/postgres"
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
