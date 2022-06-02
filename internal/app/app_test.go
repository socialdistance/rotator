package app

import (
	"context"
	"github.com/stretchr/testify/require"
	"log"
	"rotator/internal/config"
	"testing"

	internallogger "rotator/internal/logger"
	internalstorage "rotator/internal/storage/sql"
)

func TestApp(t *testing.T) {
	ctx := context.Background()

	storage := internalstorage.New(ctx, "postgres://postgres:postgres@localhost:54321/rotator?sslmode=disable")
	if err := storage.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB server", err)
	}

	logger, err := internallogger.NewLogger(config.LoggerConf{
		Level:    "debug",
		Encoding: "json",
	})
	if err != nil {
		log.Fatalf("Failed logger %s", err)
	}

	testApp := New(logger, storage)

	t.Run("Test choose banner", func(t *testing.T) {
		_, err := testApp.ChooseBanner(ctx, 1, 1)
		require.NoError(t, err)
	})

	t.Run("Test out of data", func(t *testing.T) {
		_, err := testApp.ChooseBanner(ctx, 1, 5)
		require.Error(t, err)
	})

	t.Run("Add banner to slot", func(t *testing.T) {
		err = testApp.AddBannerToSlot(ctx, 2, 2)
		require.NoError(t, err)
	})

	t.Run("Add banner to slot duplicate", func(t *testing.T) {
		err = testApp.AddBannerToSlot(ctx, 1, 1)
		require.Error(t, err)
	})

	t.Run("Remove banner from slot", func(t *testing.T) {
		err = testApp.RemoveBannerToSlot(ctx, 2, 2)
		require.NoError(t, err)
	})

	t.Run("Count transition", func(t *testing.T) {
		err = testApp.CountTransition(ctx, 2, 2, 2)
		require.NoError(t, err)
	})

}
