package repository

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/jeheskielSunloy77/go-kickstart/internal/model"
	internaltesting "github.com/jeheskielSunloy77/go-kickstart/internal/testing"
	"gorm.io/gorm"

	"github.com/stretchr/testify/require"
)

// Ensures the user repository supports CRUD operations and pagination against Postgres.
func TestUserRepository_ResourceLifecycle(t *testing.T) {
	testDB, cleanup := internaltesting.SetupTestDB(t)
	defer cleanup()

	ctx := context.Background()

	err := internaltesting.WithRollbackTransaction(ctx, testDB, func(tx *gorm.DB) error {
		repo := NewUserRepository(tx)

		user1 := &model.User{ID: uuid.New(), Email: "user1@example.com", Username: "user1"}
		user2 := &model.User{ID: uuid.New(), Email: "user2@example.com", Username: "user2"}

		require.NoError(t, repo.Store(ctx, user1))
		require.NoError(t, repo.Store(ctx, user2))
		require.NotEqual(t, uuid.Nil, user1.ID)

		fetched, err := repo.GetByID(ctx, user1.ID, nil)
		require.NoError(t, err)
		require.Equal(t, user1.ID, fetched.ID)

		updates := map[string]any{"username": "user1-updated"}
		_, err = repo.Update(ctx, *fetched, updates)
		require.NoError(t, err)

		updated, err := repo.GetByID(ctx, user1.ID, nil)
		require.NoError(t, err)
		require.Equal(t, "user1-updated", updated.Username)

		list, total, err := repo.GetMany(ctx, GetManyOptions{Limit: 1, Offset: 0, OrderBy: "created_at", OrderDirection: "asc"})
		require.NoError(t, err)
		require.Equal(t, int64(2), total)
		require.Len(t, list, 1)

		require.NoError(t, repo.Destroy(ctx, user1.ID))
		_, err = repo.GetByID(ctx, user1.ID, nil)
		require.Error(t, err)

		restored, err := repo.Restore(ctx, user1.ID)
		require.NoError(t, err)
		require.Equal(t, user1.ID, restored.ID)

		return nil
	})
	require.NoError(t, err)
}
