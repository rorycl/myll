/*
Package model provides the interaction between the database and the rest
of the go middleware.

There is a useful pgxpool connection setting blog article at
https://hexacluster.ai/postgresql/postgresql-client-side-connection-pooling-in-golang-using-pgxpool/
*/

package main

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var model Model

// Model shares a database pool for all database interactions in an app.
// Model should only be initalised through NewModel which ensures only
// one instance of a database pool exists.
type Model struct {
	ctx    context.Context
	cancel context.CancelFunc
	mu     sync.RWMutex
	inited bool // mutex protected field
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
// defaults. NewModel operates on the package global `model` and may
// only be initialised once.
func NewModel(dsn string) (*Model, error) {

	model.mu.Lock()
	defer model.mu.Unlock()
	if model.inited {
		return &model, nil
	}

	// start setting up database connection pool
	model.inited = true
	model.ctx, model.cancel = context.WithCancel(context.Background())

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

// searchPathMaker constructs a valid search_path string in the format
//
//	abc,def,'hij kl'
func searchPathMaker(schemas []string) string {
	s := pgx.Identifier(schemas)
	s.Sanitize()
	return `set search_path to ` + strings.Join(s, ",")
}

// Row provides a database pool query affecting a single row, typically
// to be used by plpgsql row-wise create/update/delete/view functions.
// Row should be called with the type of the intrinsic go struct type
// to be mapped to the returned type of the query.
//
// qCtx is the query-specific context, m is the model type providing the
// database connection pool, schemas is the (possibly empty or nil) list
// of database schemas to use from most specific to least, the query
// itself using $1 type arguments (if any) and the variadic arguments to
// fill the $1 placeholders. pgx.NamedArgs can also be used for the
// args, in which case '@' type placeholders should be used. For
// example:
//
//	type x struct{ A int }
//	n, _ := Row[x](ctx, m, nil, "select $1::int as a", 1)
//	fmt.Printf("%#v\n", n)
//	// returns 1
//
//	type y struct{ B string }
//	o, _ := Row[y](ctx, m, nil, "select @b as b", pgx.NamedArgs{"b": "hi"})
//	fmt.Printf("%#v\n", n)
//	// returns hi
func Row[T any](qCtx context.Context, m *Model, schemas []string, query string, args ...any) (T, error) {
	var value T
	if m == nil || m.inited == false {
		return value, errors.New("could not get row, model not inited")
	}
	if len(schemas) > 0 {
		search_path := searchPathMaker(schemas)
		if _, err := m.Exec(qCtx, search_path); err != nil {
			return value, fmt.Errorf("row exec error: %w", err)
		}
	}
	rows, err := m.Query(qCtx, query, args...)
	if err != nil {
		return value, fmt.Errorf("row query error: %w", err)
	}
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByName[T])
}

// Rows provides a database pool query affecting 0 to many rows,
// typically to be used by plpgsql table functions. Rows should be
// called with the type of the intrinsic go struct type to be mapped to
// a slice of returned type of the query.
//
// qCtx is the query-specific context, m is the model type providing the
// database connection pool, schemas is the (possibly empty or nil) list
// of database schemas to use from most specific to least, the query
// itself using $1 type arguments (if any) and the variadic arguments to
// fill the $1 placeholders. pgx.NamedArgs can also be used for the
// args, in which case '@' type placeholders should be used. For
// example:
//
//	users, err := Rows[User](ctx, m, search_path, "select * from users")
func Rows[T any](qCtx context.Context, m *Model, schemas []string, query string, args ...any) ([]T, error) {
	var values []T
	if m == nil || m.inited == false {
		return values, errors.New("could not get row, model not inited")
	}
	if len(schemas) > 0 {
		search_path := searchPathMaker(schemas)
		if _, err := m.Exec(qCtx, search_path); err != nil {
			return values, fmt.Errorf("row exec error: %w", err)
		}
	}
	rows, err := m.Query(qCtx, query, args...)
	if err != nil {
		return values, fmt.Errorf("row query error: %w", err)
	}
	return pgx.CollectRows(rows, pgx.RowToStructByName[T])
}

// Shutdown cancels the connections through the ctx.CancelFunc held in
// Model.
func (m *Model) Shutdown() {
	m.Close()
	m.cancel() // probably not needed
}
