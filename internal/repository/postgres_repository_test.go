package repository

import (
	"context"
	"testing"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/go-playground/assert/v2"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func TestRepository(t *testing.T) {
	ctx := context.Background()

	dbname := "test"
	dbuser := "test"
	dbpassword := "test"

	// Start Postgres container
	// 1. Start the postgres ctr and run any migrations on it
	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(dbuser),
		postgres.WithPassword(dbpassword),
		postgres.BasicWaitStrategies(),
		postgres.WithSQLDriver("pgx"),
	)
	require.NoError(t, err)

	dbURI, err := ctr.ConnectionString(ctx)
	require.NoError(t, err)

	repo, err := NewPostgresRepository(ctx, dbURI)
	require.NoError(t, err)

	err = repo.RunMigrations(ctx)
	require.NoError(t, err)

	// cleanup after tests
	testcontainers.CleanupContainer(t, ctr)
	require.NoError(t, err)

	u := &models.User{
		Login:    "user1",
		Password: "password1",
	}

	user1, err := repo.AddUser(ctx, u)
	require.NoError(t, err)
	require.NotZero(t, user1.ID)

	user1Order1 := &models.Order{
		Number: "4561261212345467", UserID: user1.ID,
	}

	_, err = repo.AddOrder(ctx, user1Order1)
	require.NoError(t, err)

	u = &models.User{
		Login:    "user2",
		Password: "password2",
	}

	user2, err := repo.AddUser(ctx, u)
	require.NoError(t, err)
	require.NotZero(t, user2.ID)

	user2Order1 := &models.Order{
		Number: "374245455400126", UserID: user2.ID,
	}
	user2Order2 := &models.Order{
		Number: "5425233430109903", UserID: user2.ID,
	}

	_, err = repo.AddOrder(ctx, user2Order1)
	require.NoError(t, err)

	_, err = repo.AddOrder(ctx, user2Order2)
	require.NoError(t, err)

	// find non-existing user, should be an error
	_, err = repo.FindUserByLogin(ctx, "userunknown")
	require.ErrorIs(t, err, common.ErrorNotFound)

	//find user1
	t.Run("FinUser1ByLogin", func(t *testing.T) {
		user, err := repo.FindUserByLogin(ctx, "user1")
		require.NoError(t, err)
		assert.Equal(t, user.Login, user1.Login)
	})

	t.Run("FindNonExistingOrder", func(t *testing.T) {
		_, err = repo.FindOrderByNumber(ctx, "123123")
		require.ErrorIs(t, err, common.ErrorNotFound)
	})

	t.Run("AddWithdrawalToUser1", func(t *testing.T) {
		w := &models.Withdrawal{UserID: user1.ID, Amount: 1}
		err = repo.AddWithdrawal(ctx, w)
		require.NoError(t, err)
	})

	t.Run("AddAnotherWithdrawalToUser1", func(t *testing.T) {
		w := &models.Withdrawal{UserID: user1.ID, Amount: 2.50}
		err = repo.AddWithdrawal(ctx, w)
		require.NoError(t, err)
	})

	t.Run("GetWithdrawalsTotalAmountByUserID", func(t *testing.T) {
		res, err := repo.GetWithdrawalsTotalAmountByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, res, float32(3.50))
	})

}
