package sql

import (
	"context"
	"errors"
	"fmt"
	pgx4 "github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"os"
)

type Storage struct {
	ctx  context.Context
	conn *pgxpool.Pool
	dsn  string
}

type Banner struct {
	ID           int64  `db:"banner_id"`
	Description  string `db:"banner_description"`
	TotalDisplay int64  `db:"total_display"`
}

type Slot struct {
	ID           int64  `db:"slot_id"`
	Description  string `db:"slot_description"`
	TotalDisplay int64  `db:"total_display"`
}

type SocialGroup struct {
	ID          int64  `db:"social_group_id"`
	Description string `db:"description"`
}

type BannerStats struct {
	ID           int64 `db:"banner_id"`
	Display      int64 `db:"display"`
	Click        int64 `db:"click"`
	TotalDisplay int64 `db:"total_display"`
}

func New(ctx context.Context, dsn string) *Storage {
	return &Storage{
		ctx: ctx,
		dsn: dsn,
	}
}

func (s *Storage) Connect(ctx context.Context) error {
	conn, err := pgxpool.Connect(ctx, s.dsn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect database %s", err)
	}

	s.conn = conn

	return nil
}

func (s *Storage) Close() {
	s.conn.Close()
}

func (s *Storage) GetBannerId(ctx context.Context, bannerID int64) (*Banner, error) {
	var b Banner

	sql := `
		SELECT banner_id, banner_description, total_display FROM banner WHERE banner_id = $1
	`

	err := s.conn.QueryRow(ctx, sql, bannerID).Scan(
		&b.ID, &b.Description, &b.TotalDisplay)

	if err == nil {
		return &b, nil
	}

	if errors.Is(err, pgx4.ErrNoRows) {
		return nil, nil
	}

	return nil, fmt.Errorf("cant scan SQL result to struct %w", err)
}

func (s *Storage) GetSlotByID(ctx context.Context, slotID int64) (*Slot, error) {
	var slot Slot

	sql := `
		SELECT slot_id, slot_description, total_display FROM slot WHERE slot_id = $1
	`

	err := s.conn.QueryRow(ctx, sql, slotID).Scan(
		&slot.ID, &slot.Description, &slot.TotalDisplay)

	if err == nil {
		return &slot, nil
	}

	if errors.Is(err, pgx4.ErrNoRows) {
		return nil, nil
	}

	return nil, fmt.Errorf("cant scan SQL result to struct %w", err)
}

func (s *Storage) GetSocialGroupByID(ctx context.Context, socialGroupID int64) (*SocialGroup, error) {
	var socialGroup SocialGroup

	sql := `
		SELECT social_group_id, social_description FROM social_group WHERE social_group_id = $1
	`

	err := s.conn.QueryRow(ctx, sql, socialGroupID).Scan(
		&socialGroup.ID, &socialGroup.Description)

	if err == nil {
		return &socialGroup, nil
	}

	if errors.Is(err, pgx4.ErrNoRows) {
		return nil, nil
	}

	return nil, fmt.Errorf("cant scan SQL result to struct %w", err)
}

// AddBannerToSlot relation banner <-> slot
func (s *Storage) AddBannerToSlot(ctx context.Context, bannerID, slotID int64) error {
	tx, err := s.conn.BeginTx(ctx, pgx4.TxOptions{
		IsoLevel:       pgx4.Serializable,
		AccessMode:     pgx4.ReadWrite,
		DeferrableMode: pgx4.NotDeferrable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		INSERT INTO banner_to_slot (banner_id, slot_id) VALUES ($1, $2)
	`
	_, err = tx.Exec(ctx, query, bannerID, slotID)
	if err != nil {
		return err
	}

	query = `
		INSERT INTO statistics (banner_id, social_group_id, slot_id) SELECT banner_id, social_group_id, slot_id 
		FROM banner_to_slot CROSS JOIN social_group WHERE banner_id = $1 AND slot_id = $2
	`

	_, err = tx.Exec(ctx, query, bannerID, slotID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

func (s *Storage) RemoveBannerFromSlot(ctx context.Context, bannerID, slotID int64) error {
	tx, err := s.conn.BeginTx(ctx, pgx4.TxOptions{
		IsoLevel:       pgx4.Serializable,
		AccessMode:     pgx4.ReadWrite,
		DeferrableMode: pgx4.NotDeferrable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		DELETE FROM banner_to_slot WHERE banner_id = $1 AND slot_id = $2
	`
	_, err = tx.Exec(ctx, query, bannerID, slotID)
	if err != nil {
		return err
	}

	query = `
		DELETE FROM statistics WHERE banner_id = $1 AND slot_id = $2
	`

	_, err = tx.Exec(ctx, query, bannerID, slotID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// CountTransition Регистрирует переход
func (s *Storage) CountTransition(ctx context.Context, bannerID, slotID, socialGroupID int64) error {
	query := `UPDATE statistics SET click = click + 1
		WHERE slot_id = $1 AND banner_id = $2 AND social_group_id = $3`

	result, err := s.conn.Exec(ctx, query, bannerID, slotID, socialGroupID)
	if err != nil {
		return fmt.Errorf("can't count transition slot %d banner = %d social group %d: %w", slotID, bannerID, socialGroupID, err)
	}

	count := result.RowsAffected()
	if count == 0 {
		return err
	}

	return nil
}

// CountDisplay Регистрирует показ баннера
func (s *Storage) CountDisplay(ctx context.Context, bannerID, slotID, socialGroupID int64) error {
	tx, err := s.conn.BeginTx(ctx, pgx4.TxOptions{
		IsoLevel:       pgx4.Serializable,
		AccessMode:     pgx4.ReadWrite,
		DeferrableMode: pgx4.NotDeferrable,
	})
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	query := `
		UPDATE statistics SET display = display + 1 WHERE slot_id = $1 AND banner_id = $2 AND social_group_id = $3
	`
	_, err = tx.Exec(ctx, query, bannerID, slotID, socialGroupID)
	if err != nil {
		return err
	}

	query = `
		UPDATE slot SET total_display = total_display + 1 WHERE slot_id = $1
	`

	_, err = tx.Exec(ctx, query, slotID)
	if err != nil {
		return err
	}

	if err = tx.Commit(ctx); err != nil {
		return err
	}

	return nil
}

// GetBannersStat Выбирает баннеры с их статистиками
// которые могут быть показаны в указанном слоте и для указанной соц.группы.
func (s *Storage) GetBannersStat(ctx context.Context, slotID, socialGroupID int64) ([]BannerStats, int, error) {
	result := make([]BannerStats, 0)

	query := `
		SELECT banner_id, display, click FROM statistics WHERE slot_id = $1 AND social_group_id = $2
	`

	rows, err := s.conn.Query(ctx, query, slotID, socialGroupID)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	for rows.Next() {
		var b BannerStats
		if err := rows.Scan(&b.ID, &b.Display, &b.Click); err != nil {
			return nil, 0, fmt.Errorf("cant convert result: %w", err)
		}

		result = append(result, b)
	}

	var totalDisplay int
	query = `SELECT total_display FROM slot WHERE slot_id = $1`
	err = s.conn.QueryRow(ctx, query, slotID).Scan(&totalDisplay)

	if errors.Is(err, pgx4.ErrNoRows) {
		return nil, 0, nil
	}

	return result, totalDisplay, err
}
