run:
	go run cmd/main.go

proto-gen:
	./scripts/gen-proto.sh

migrate_up:
	migrate -path migrations -database postgres://postgres:bnnfav@localhost:5432/suserdb -verbose up

migrate_down:
	migrate -path migrations -database postgres://postgres:bnnfav@localhost:5432/suserdb -verbose down

migrate_force:
	migrate -path migrations -database postgres://postgres:bnnfav@localhost:5432/suserdb -verbose force 0

.PHONY: start migrateup migratedown