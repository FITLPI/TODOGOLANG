// Package core_postgres_pool — антикоррупционный слой над пулом pgx:
// repo фич импортируют только этот пакет и нейтральные ошибки (ErrNoRows,
// ErrViolatesForeignKey); pgx-специфичные ошибки транслируются здесь,
// в mapErrors.
package core_postgres_pool

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	// ErrNoRows — запрос не вернул ни одной строки (аналог sql.ErrNoRows).
	ErrNoRows = errors.New("no rows in result set")
	// ErrViolatesForeignKey — нарушение внешнего ключа (SQLSTATE 23503).
	ErrViolatesForeignKey = errors.New("violates foreign key constraint")
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, args ...any) (CommandTag, error)
	Close()

	OperationTimeout() time.Duration
}

// Rows — нейтральное отображение pgx.Rows.
type Rows interface {
	Next() bool
	Scan(dest ...any) error
	Values() ([]any, error)
	Err() error
	Close()
	RawRows() any
}

// Row — нейтральное отображение pgx.Row.
type Row interface {
	Scan(dest ...any) error
	RawRow() any
}

// CommandTag — нейтральное отображение pgconn.CommandTag.
type CommandTag interface {
	RowsAffected() int64
}

// mapErrors — перевод pgx-ошибок в нейтральные ошибки пула.
func mapErrors(err error) error {
	if err == nil {
		return nil
	}

	if errors.Is(err, pgx.ErrNoRows) {
		return fmt.Errorf("%w: %s", ErrNoRows, err)
	}

	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23503" {
		return fmt.Errorf("%w: %s", ErrViolatesForeignKey, pgErr.Message)
	}

	return err
}

// pgxRows — адаптер pgx.Rows -> Rows.
type pgxRows struct {
	pgx.Rows
}

func (r pgxRows) Scan(dest ...any) error {
	return mapErrors(r.Rows.Scan(dest...))
}

func (r pgxRows) Values() ([]any, error) {
	values, err := r.Rows.Values()
	return values, mapErrors(err)
}

func (r pgxRows) Err() error {
	return mapErrors(r.Rows.Err())
}

func (r pgxRows) RawRows() any {
	return r.Rows
}

// pgxRow — адаптер pgx.Row -> Row.
type pgxRow struct {
	pgx.Row
}

func (row pgxRow) Scan(dest ...any) error {
	return mapErrors(row.Row.Scan(dest...))
}

func (row pgxRow) RawRow() any {
	return row.Row
}

// pgxCommandTag — адаптер pgconn.CommandTag -> CommandTag.
type pgxCommandTag struct {
	pgconn.CommandTag
}

type ConnectionPool struct {
	*pgxpool.Pool
	opTimeout time.Duration
}

func (p *ConnectionPool) Query(ctx context.Context, sql string, args ...any) (Rows, error) {
	rows, err := p.Pool.Query(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return pgxRows{rows}, nil
}

func (p *ConnectionPool) QueryRow(ctx context.Context, sql string, args ...any) Row {
	return pgxRow{p.Pool.QueryRow(ctx, sql, args...)}
}

func (p *ConnectionPool) Exec(ctx context.Context, sql string, args ...any) (CommandTag, error) {
	cmdTag, err := p.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return nil, mapErrors(err)
	}

	return pgxCommandTag{cmdTag}, nil
}

func NewConnectionPool(ctx context.Context, config Config) (*ConnectionPool, error) {
	conectionString := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	pgxConfig, err := pgxpool.ParseConfig(conectionString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse connection string: %w", err)
	}

	pool, err := pgxpool.NewWithConfig(ctx, pgxConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create connection pool: %w", err)
	}

	if err := pool.Ping(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return &ConnectionPool{
		Pool:      pool,
		opTimeout: config.Timeout,
	}, nil
}

func (p *ConnectionPool) OperationTimeout() time.Duration {
	return p.opTimeout
}
