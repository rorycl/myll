/*
Package model provides the interaction between the database and the rest
of the go middleware.

There is a useful pgxpool connection setting blog article at
https://hexacluster.ai/postgresql/postgresql-client-side-connection-pooling-in-golang-using-pgxpool/
*/

package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

var model Model

type Model struct {
	ctx    context.Context
	cancel context.CancelFunc
	inited bool
	*pgxpool.Pool
}

// default database pool settings
var (
	defaultMinConns              int32         = 2                // default 0
	defaultMaxConns              int32         = 8                // default 4
	defaultMaxConnLifetime       time.Duration = 12 * time.Minute // default 1 hour
	defaultMaxConnIdleTime                     = 10 * time.Minute // default 30 minutes
	defaultHealthCheckPeriod                   = 1 * time.Minute  // default 1 minute
	defaultMaxConnLifetimeJitter               = 10 * time.Second // default 0
)

// NewModel creates a new database model with database context and
// pgxpool connection settings as supplied or as set out by the package
// defaults.
func NewModel(dsn string) (*Model, error) {
	if model.inited {
		return &model, nil
	}
	ctx, cancel := context.WithCancel(context.Background())
	model = Model{
		ctx:    ctx,
		cancel: cancel,
		inited: true,
	}
	// parse configuration
	config, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("database parseconfig error: %w", err)
	}

	// override defaults if not changed by dsn
	if config.MinConns == 4 {
		config.MinConns = defaultMinConns
	}
	if config.MaxConns == 8 {
		config.MaxConns = defaultMaxConns
	}
	if config.MaxConnLifetime == 1*time.Hour {
		config.MaxConnLifetime = defaultMaxConnLifetime
	}
	if config.MaxConnIdleTime == 30*time.Minute {
		config.MaxConnIdleTime = defaultMaxConnIdleTime
	}
	if config.HealthCheckPeriod == 1*time.Minute {
		config.HealthCheckPeriod = 1 * time.Minute
	}
	if config.MaxConnLifetimeJitter == 0 {
		config.MaxConnLifetimeJitter = defaultMaxConnLifetimeJitter
	}

	// make new database pool
	model.Pool, err = pgxpool.NewWithConfig(model.ctx, config)
	if err != nil {
		return nil, fmt.Errorf("database pool connection error: %w", err)
	}
	// test the connection
	err = model.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database ping error: %w", err)
	}
	return &model, nil
}

// Shutdown cancels the connections through the ctx.CancelFunc held in
// Model.
func (m *Model) Shutdown() {
	m.Close()
	m.cancel() // probably not needed
}
