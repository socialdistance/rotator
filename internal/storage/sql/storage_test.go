package sql

import (
	"context"
	pgx4 "github.com/jackc/pgx/v4"
	"github.com/stretchr/testify/require"

	"testing"
)

func TestStorage(t *testing.T) {
	ctx := context.Background()
	storage := New(ctx, "postgres://postgres:postgres@localhost:54321/rotator?sslmode=disable")
	if err := storage.Connect(ctx); err != nil {
		t.Fatal("Failed to connect to DB server", err)
	}

	t.Run("test SQL", func(t *testing.T) {
		tx, err := storage.conn.BeginTx(ctx, pgx4.TxOptions{
			IsoLevel:       pgx4.Serializable,
			AccessMode:     pgx4.ReadWrite,
			DeferrableMode: pgx4.NotDeferrable,
		})
		if err != nil {
			t.Fatal("Failed to connect to DB server", err)
		}

		_, err = storage.GetBannerId(ctx, 1)
		require.NoError(t, err)

		_, err = storage.GetSlotByID(ctx, 1)
		require.NoError(t, err)

		_, err = storage.GetSocialGroupByID(ctx, 2)
		require.NoError(t, err)

		//addBannerToSlot := storage.AddBannerToSlot(ctx, 2, 1)
		//require.NoError(t, addBannerToSlot)

		//removeBannerToSlot := storage.RemoveBannerFromSlot(ctx, 1, 1)
		//require.NoError(t, removeBannerToSlot)

		err = storage.CountTransition(ctx, 1, 1, 2)
		require.NoError(t, err)

		err = storage.CountDisplay(ctx, 1, 1, 2)
		require.NoError(t, err)

		_, _, err = storage.GetBannersStat(ctx, 1, 1)
		require.NoError(t, err)

		err = tx.Rollback(ctx)
		if err != nil {
			t.Fatal("Failed to rollback tx", err)
		}
	})
}
