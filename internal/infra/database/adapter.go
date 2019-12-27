package database

import (
	"context"
	"database/sql"
	"errors"

	"github.com/L11R/wotbot/internal/domain"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"go.uber.org/zap"
)

type adapter struct {
	logger *zap.Logger
	config *Config
	db     *sqlx.DB
}

func NewAdapter(logger *zap.Logger, config *Config) (domain.Database, error) {
	a := &adapter{
		logger: logger,
		config: config,
	}

	db, err := sqlx.Open("postgres", config.ConnectionString())
	if err != nil {
		return nil, err
	}
	a.db = db

	db.SetMaxOpenConns(config.MaxOpenConns)
	db.SetMaxIdleConns(config.MaxIdleConns)
	db.SetConnMaxLifetime(config.ConnMaxLifeTime)

	// Migrations block
	driver, err := postgres.WithInstance(db.DB, &postgres.Config{})
	if err != nil {
		return nil, err
	}

	m, err := migrate.NewWithDatabaseInstance(config.MigrationsSourceURL, config.Name, driver)
	if err != nil {
		return nil, err
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return nil, err
	}

	return a, nil
}

func (a *adapter) GetUserByTelegramID(telegramID int) (*domain.User, error) {
	row := a.db.QueryRowx(`SELECT id, telegram_id, nickname, wargaming_id, created_at, updated_at FROM users WHERE telegram_id = $1`, telegramID)
	if row.Err() != nil {
		if errors.Is(row.Err(), sql.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		a.logger.Error("Error getting user!", zap.Error(row.Err()))
		return nil, domain.ErrInternalDatabase
	}

	var res domain.User
	if err := row.StructScan(&res); err != nil {
		a.logger.Error("Error scanning result!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}

	return &res, nil
}

func (a *adapter) UpsertUser(user *domain.User) (*domain.User, error) {
	rows, err := a.db.NamedQuery(`INSERT INTO users (telegram_id, nickname, wargaming_id)
VALUES (:telegram_id, :nickname, :wargaming_id)
ON CONFLICT (telegram_id) DO UPDATE SET nickname = COALESCE(EXCLUDED.nickname, users.nickname), wargaming_id = COALESCE(EXCLUDED.wargaming_id, users.wargaming_id) RETURNING id, telegram_id, nickname, wargaming_id, created_at, updated_at;`, user)
	if err != nil {
		a.logger.Error("Error upserting user!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	for rows.Next() {
		var res domain.User
		if err := rows.StructScan(&res); err != nil {
			a.logger.Error("Error scanning result!", zap.Error(err))
			return nil, domain.ErrInternalDatabase
		}

		return &res, nil
	}

	a.logger.Error("There is no results to scan!")
	return nil, domain.ErrInternalDatabase
}

func (a *adapter) GetStatsByUserID(userID int) ([]*domain.Stat, error) {
	rows, err := a.db.Queryx(`SELECT * FROM stats WHERE user_id = $1`, userID)
	if err != nil {
		a.logger.Error("Error selecting stats!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	results := make([]*domain.Stat, 0)
	for rows.Next() {
		var res domain.Stat
		if err := rows.StructScan(&res); err != nil {
			a.logger.Error("Error scanning result!", zap.Error(err))
			return nil, domain.ErrInternalDatabase
		}

		results = append(results, &res)
	}

	return results, nil
}

func (a *adapter) UpdateStatsByUserID(userID int, stats []*domain.Stat) ([]*domain.Stat, error) {
	tx, err := a.db.BeginTxx(context.Background(), nil)
	if err != nil {
		a.logger.Error("Error beginning database transaction!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}

	defer func(err *error) {
		if err != nil && *err != nil {
			if err := tx.Rollback(); err != nil {
				a.logger.Error("Error while rollback transaction!", zap.Error(err))
			}
		}
	}(&err)

	_, err = tx.Exec(`DELETE FROM stats WHERE user_id = $1`, userID)
	if err != nil {
		a.logger.Error("Error deleting old stats!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}

	// Set user_id
	for i := range stats {
		stats[i].UserID = userID
	}

	_, err = tx.NamedExec(
		`INSERT INTO stats (user_id, type, name, value, html_id, img) VALUES (:user_id, :type, :name, :value, :html_id, :img)`,
		stats,
	)
	if err != nil {
		a.logger.Error("Error inserting new stats!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}

	err = tx.Commit()
	if err != nil {
		a.logger.Error("Error committing transaction!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}

	rows, err := a.db.Queryx(`SELECT * FROM stats WHERE user_id = $1`, userID)
	if err != nil {
		a.logger.Error("Error selecting stats!", zap.Error(err))
		return nil, domain.ErrInternalDatabase
	}
	//noinspection GoUnhandledErrorResult
	defer rows.Close()

	results := make([]*domain.Stat, 0)
	for rows.Next() {
		var res domain.Stat
		if err := rows.StructScan(&res); err != nil {
			a.logger.Error("Error scanning result!", zap.Error(err))
			return nil, domain.ErrInternalDatabase
		}

		results = append(results, &res)
	}

	return results, nil
}
