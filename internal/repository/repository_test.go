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

func TestRepositoryImplementations(t *testing.T) {

	ctx := context.Background()

	t.Run("InMemory", func(t *testing.T) {
		repo, err := NewInMemoryRepository()
		require.NoError(t, err)
		RunRepositoryTests(t, ctx, "InMemory", repo)
	})

	t.Run("Postgres", func(t *testing.T) {

		dbname := "test"
		dbuser := "test"
		dbpassword := "test"

		// Start the postgres ctr and run any migrations on it
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

		RunRepositoryTests(t, ctx, "Postgres", repo)

		testcontainers.CleanupContainer(t, ctr)

	})
}

func RunRepositoryTests(t *testing.T, ctx context.Context, name string, repo Repository) {

	var err error

	var user1, user2 models.User
	var user1order1 models.Order

	t.Run(name+"AddUser1", func(t *testing.T) {

		u := &models.User{
			Login:    "user1",
			Password: "password1",
		}

		user1, err = repo.AddUser(ctx, u)
		require.NoError(t, err)
		require.NotZero(t, user1.ID)

	})

	t.Run(name+"AddOrder1", func(t *testing.T) {
		user1order1, err = repo.AddOrder(ctx, &models.Order{
			Number: "4561261212345467", UserID: user1.ID, Status: models.OrderStatusNew,
		})
		require.NoError(t, err)
		require.NotZero(t, user1order1.ID)
	})

	t.Run(name+"AddUser2", func(t *testing.T) {
		u := &models.User{
			Login:    "user2",
			Password: "password2",
		}

		user2, err = repo.AddUser(ctx, u)
		require.NoError(t, err)
		require.NotZero(t, user2.ID)
	})

	t.Run(name+"AdOrders", func(t *testing.T) {
		user2Order1 := &models.Order{
			Number: "374245455400126", UserID: user2.ID, Status: models.OrderStatusNew,
		}
		user2Order2 := &models.Order{
			Number: "5425233430109903", UserID: user2.ID, Status: models.OrderStatusProcessing,
		}

		_, err = repo.AddOrder(ctx, user2Order1)
		require.NoError(t, err)

		_, err = repo.AddOrder(ctx, user2Order2)
		require.NoError(t, err)
	})

	t.Run(name+"FindNonExistingUser", func(t *testing.T) {
		orders, err := repo.GetOrdersByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, len(orders), 1)
		assert.Equal(t, orders[0].Number, "4561261212345467")
	})

	t.Run(name+"FindNonExistingUser", func(t *testing.T) {
		// find non-existing user, should be an error
		_, err = repo.FindUserByLogin(ctx, "userunknown")
		require.ErrorIs(t, err, common.ErrorNotFound)
	})

	//find user1
	t.Run(name+"FinUser1ByLogin", func(t *testing.T) {
		user, err := repo.FindUserByLogin(ctx, "user1")
		require.NoError(t, err)
		assert.Equal(t, user.Login, user1.Login)
	})

	t.Run(name+"FindNonExistingOrder", func(t *testing.T) {
		_, err = repo.FindOrderByNumber(ctx, "123123")
		require.ErrorIs(t, err, common.ErrorNotFound)
	})

	t.Run(name+"FindOrderByNumber", func(t *testing.T) {
		order, err := repo.FindOrderByNumber(ctx, "5425233430109903")
		require.NoError(t, err)
		assert.Equal(t, order.UserID, user2.ID)
	})

	t.Run(name+"AddWithdrawalsToUser1", func(t *testing.T) {
		w := &models.Withdrawal{UserID: user1.ID, Amount: 1}
		err = repo.AddWithdrawal(ctx, w)
		require.NoError(t, err)

		w = &models.Withdrawal{UserID: user1.ID, Amount: 2.50}
		err = repo.AddWithdrawal(ctx, w)
		require.NoError(t, err)
	})

	t.Run(name+"AddWithdrawalsToUser2", func(t *testing.T) {
		w := &models.Withdrawal{UserID: user2.ID, Amount: 4.50}
		err = repo.AddWithdrawal(ctx, w)
		require.NoError(t, err)
	})

	t.Run(name+"GetWithdrawalsByUserID1", func(t *testing.T) {
		res, err := repo.GetWithdrawalsByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, len(res), 2)
	})

	t.Run(name+"GetWithdrawalsByUserID2", func(t *testing.T) {
		res, err := repo.GetWithdrawalsByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Equal(t, len(res), 1)
	})

	t.Run(name+"GetUnprocessedOrdersTry1", func(t *testing.T) {
		res, err := repo.GetUnprocessedOrders(ctx)
		require.NoError(t, err)
		assert.Equal(t, len(res), 3)
	})

	t.Run(name+"GetUnprocessedOrdersAfterAccrual", func(t *testing.T) {
		err := repo.UpdateOrderAccrualStatus(ctx, user1order1.ID, models.OrderStatusProcessed, 5.0)
		require.NoError(t, err)

		res, err := repo.GetUnprocessedOrders(ctx)
		require.NoError(t, err)
		assert.Equal(t, len(res), 2)

	})

	t.Run(name+"GetWithdrawalsTotalAmountByUserID2", func(t *testing.T) {
		res, err := repo.GetWithdrawalsTotalAmountByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Equal(t, res, float32(4.50))
	})

	t.Run(name+"GetWithdrawalsTotalAmountByUserID1", func(t *testing.T) {
		res, err := repo.GetWithdrawalsTotalAmountByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, res, float32(3.50))
	})

	t.Run(name+"GetAccrualsTotalAmountByUserID2", func(t *testing.T) {
		res, err := repo.GetAccrualsTotalAmountByUserID(ctx, user2.ID)
		require.NoError(t, err)
		assert.Equal(t, res, float32(0))
	})

	t.Run(name+"GetAccrualsTotalAmountByUserID1", func(t *testing.T) {
		res, err := repo.GetAccrualsTotalAmountByUserID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, res, float32(5))
	})

	t.Run(name+"CheckUserBalanceAfterRecalculation", func(t *testing.T) {
		err := repo.UpdateUserAccruedTotal(ctx, user1.ID, 5)
		require.NoError(t, err)

		err = repo.UpdateUserWithdrawnTotal(ctx, user1.ID, 3.5)
		require.NoError(t, err)

		user, err := repo.FindUserByID(ctx, user1.ID)
		require.NoError(t, err)
		assert.Equal(t, user.ID, user1.ID)

		assert.Equal(t, user.AccruedTotal, float32(5))
		assert.Equal(t, user.WithdrawnTotal, float32(3.5))

	})

}
