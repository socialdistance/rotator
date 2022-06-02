package app

import (
	"context"
	"go.uber.org/zap"
	bandit "rotator/internal/alghoritms"
	sqlstorage "rotator/internal/storage/sql"
	"time"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	Debug(message string, fields ...zap.Field)
	Info(message string, fields ...zap.Field)
	Error(message string, fields ...zap.Field)
	Fatal(message string, fields ...zap.Field)
}

type Storage interface {
	GetBannerId(ctx context.Context, bannerID int64) (*sqlstorage.Banner, error)
	GetSlotByID(ctx context.Context, slotID int64) (*sqlstorage.Slot, error)
	GetSocialGroupByID(ctx context.Context, socialGroupID int64) (*sqlstorage.SocialGroup, error)
	AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error
	RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error
	CountTransition(ctx context.Context, bannerID, slotID, socialGroupID int64) error
	CountDisplay(ctx context.Context, bannerID, slotID, socialGroupID int64) error
	GetBannersStat(ctx context.Context, slotID, socialGroupID int64) ([]sqlstorage.BannerStats, int, error)
}

func New(logger Logger, storage Storage) *App {
	return &App{
		Logger:  logger,
		Storage: storage,
	}
}

func (a *App) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	opCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	return a.Storage.AddBannerToSlot(opCtx, bannerID, slotID)
}

func (a *App) RemoveBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	opCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	return a.Storage.RemoveBannerFromSlot(opCtx, bannerID, slotID)
}

func (a *App) CountTransition(ctx context.Context, bannerID, slotID, socialGroupID int64) error {
	opCtx, cancel := context.WithTimeout(ctx, time.Second*3)
	defer cancel()

	return a.Storage.CountTransition(opCtx, bannerID, slotID, socialGroupID)
}

func (a *App) ChooseBanner(ctx context.Context, slotID, socialGroupID int64) (bannerID int64, err error) {
	opCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	bannerStat, totalDisplay, err := a.Storage.GetBannersStat(opCtx, slotID, socialGroupID)
	if err != nil {
		return 0, err
	}

	stat := make([]bandit.Bandit, len(bannerStat))
	for i, v := range bannerStat {
		stat[i] = bandit.Bandit{
			ID:     int(v.ID),
			Trials: int(v.Display),
			Reward: int(v.Click),
		}
	}

	banner, err := bandit.ChooseAlgorithm(stat, totalDisplay)
	if err != nil {
		return 0, err
	}

	err = a.Storage.CountTransition(ctx, int64(banner), slotID, socialGroupID)
	if err != nil {
		return 0, err
	}

	return int64(banner), nil
}
