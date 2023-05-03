package tests

import (
	"log"
	"os"
	"testing"

	"github.com/burxondv/new-services/post-service/config"
	"github.com/burxondv/new-services/post-service/pkg/db"
	"github.com/burxondv/new-services/post-service/storage/postgres"
)

var pgRepo *postgres.PostRepo

func TestMain(m *testing.M) {
	conf := config.Load()
	connDB, err := db.ConnectToDB(conf)
	if err != nil {
		log.Println("failed to connect to database: ", err)
	}

	pgRepo = postgres.NewPostRepo(connDB)

	os.Exit(m.Run())
}
