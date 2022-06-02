package store

import (
	"context"
	"log"
	internalapp "rotator/internal/app"
	internalconfig "rotator/internal/config"
	"rotator/internal/storage/sql"
)

func CreateStorage(ctx context.Context, config internalconfig.Config) internalapp.Storage {
	var store internalapp.Storage

	switch config.Storage.Type {
	case internalconfig.SQL:
		sqlStore := sql.New(ctx, config.Storage.Dsn)
		if err := sqlStore.Connect(ctx); err != nil {
			log.Fatalf("Unable to connect database %s", err)
		}
		store = sqlStore
	default:
		log.Fatalf("Dont know type storage: %s", config.Storage.Type)
	}

	return store
}
