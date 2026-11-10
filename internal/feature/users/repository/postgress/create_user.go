package users_postgres_repository

import (
	"context"

	"github.com/FITLPI/TODOGOLANG/internal/core/domain"
)

func (r *UserRepository) CreateUser(
	ctx context.Context,
	user domain.User,
) (domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	/

	query := `

	`
}
