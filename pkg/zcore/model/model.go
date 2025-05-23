package model

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/pgx/v5"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/pkg/errors"
	"github.com/risingwavelabs/wavekit"
	"github.com/risingwavelabs/wavekit/pkg/config"
	"github.com/risingwavelabs/wavekit/pkg/logger"
	"github.com/risingwavelabs/wavekit/pkg/zgen/querier"
)

var log = logger.NewLogAgent("model")

var (
	ErrAlreadyInTransaction = errors.New("already in transaction")
)

type ModelInterface interface {
	querier.Querier
	RunTransaction(ctx context.Context, f func(model ModelInterface) error) error
	RunTransactionWithTx(ctx context.Context, f func(tx pgx.Tx, model ModelInterface) error) error
	InTransaction() bool
	SpawnWithTx(tx pgx.Tx) ModelInterface
}

type Model struct {
	querier.Querier
	beginTx       func(ctx context.Context) (pgx.Tx, error)
	p             *pgxpool.Pool
	inTransaction bool
	now           func() time.Time
}

func (m *Model) InTransaction() bool {
	return m.inTransaction
}

func (m *Model) SpawnWithTx(tx pgx.Tx) ModelInterface {
	return &Model{
		Querier: querier.New(tx),
		beginTx: func(ctx context.Context) (pgx.Tx, error) {
			return nil, ErrAlreadyInTransaction
		},
		inTransaction: true,
	}
}

func (m *Model) RunTransactionWithTx(ctx context.Context, f func(tx pgx.Tx, model ModelInterface) error) error {
	tx, err := m.beginTx(ctx)
	if err != nil {
		return err
	}
	defer tx.Rollback(ctx)

	txm := m.SpawnWithTx(tx)

	if err := f(tx, txm); err != nil {
		return err
	}

	return tx.Commit(ctx)
}

func (m *Model) RunTransaction(ctx context.Context, f func(model ModelInterface) error) error {
	return m.RunTransactionWithTx(ctx, func(_ pgx.Tx, model ModelInterface) error {
		return f(model)
	})
}

func NewModel(cfg *config.Config) (ModelInterface, error) {
	if cfg.Pg.DSN == nil {
		return nil, errors.New("dsn is not set")
	}
	dsn := *cfg.Pg.DSN

	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse pgxpool config: %s", dsn)
	}
	config.MaxConns = 30
	config.MinConns = 5

	var (
		retryLimit = 10
		retry      = 0
	)

	var p *pgxpool.Pool

	for {
		ctx, cancel := context.WithTimeout(context.TODO(), 10*time.Second)
		defer cancel()

		pool, err := pgxpool.NewWithConfig(ctx, config)
		if err != nil {
			log.Warnf("failed to init pgxpool: %s", err.Error())
			if retry >= retryLimit {
				return nil, errors.Wrapf(err, "failed to init pgxpool: %s", dsn)
			}
			continue
		}

		p = pool

		if err := pool.Ping(ctx); err != nil {
			log.Warnf("failed to ping database: %s", err.Error())
			if retry >= retryLimit {
				return nil, errors.Wrap(err, "failed to ping db")
			}
		} else {
			break
		}
		retry++
		time.Sleep(3 * time.Second)
	}

	d, err := iofs.New(wavekit.Migrations, "sql/migrations")
	if err != nil {
		return nil, errors.Wrap(err, "failed to create migration source driver")
	}

	url := fmt.Sprintf("pgx5://%s:%s@%s:%d/%s",
		config.ConnConfig.User,
		config.ConnConfig.Password,
		config.ConnConfig.Host,
		config.ConnConfig.Port,
		config.ConnConfig.Database,
	)
	m, err := migrate.NewWithSourceInstance("iofs", d, url)
	if err != nil {
		return nil, errors.Wrap(err, "failed to init migrate")
	}
	if err := m.Up(); err != nil {
		if !errors.Is(err, migrate.ErrNoChange) {
			return nil, errors.Wrap(err, "failed to migrate up")
		}
	}

	return &Model{Querier: querier.New(p), beginTx: p.Begin, p: p, now: time.Now}, nil
}
