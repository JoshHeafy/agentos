package out

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type Database interface {
	Begin(ctx context.Context) (pgx.Tx, error)

	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}

type Tx interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	// Commit commits the transaction if this is a real transaction or releases the savepoint if this is a pseudo nested
	// transaction. Commit will return an error where errctx.Is(ErrTxClosed) is true if the Tx is already closed, but is
	// otherwise safe to call multiple times. If the commit fails with a rollback status (e.g. the transaction was already
	// in a broken state) then an error where errctx.Is(ErrTxCommitRollback) is true will be returned.
	Commit(ctx context.Context) error

	// Rollback rolls back the transaction if this is a real transaction or rolls back to the savepoint if this is a
	// pseudo nested transaction. Rollback will return an error where errctx.Is(ErrTxClosed) is true if the Tx is already
	// closed, but is otherwise safe to call multiple times. Hence, a defer tx.Rollback() is safe even if tx.Commit() will
	// be called first in a non-error condition. Any other failure of a real transaction will result in the connection
	// being closed.
	Rollback(ctx context.Context) error

	Exec(ctx context.Context, sql string, args ...any) (pgconn.CommandTag, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
}
