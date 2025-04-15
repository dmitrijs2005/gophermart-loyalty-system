package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
	"github.com/dmitrijs2005/gophermart-loyalty-system/migrations"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(ctx context.Context, dsn string) (*PostgresRepository, error) {
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) RunMigrations(ctx context.Context) error {

	goose.SetBaseFS(migrations.Migrations) // Вот здесь передаём embed FS!

	if err := goose.UpContext(ctx, r.db, "."); err != nil {
		return err
	}

	return nil

}

func (r *PostgresRepository) UnitOfWork() UnitOfWork {
	return &PgUnitOfWork{r.db}
}

func (r *PostgresRepository) FindUserByLogin(ctx context.Context, login string) (models.User, error) {

	s := "select id, login, password, salt from users where login=$1"

	var user models.User

	_, err := common.RetryWithResult(ctx, func() (*sql.Row, error) {
		r := r.db.QueryRowContext(ctx, s, login)
		err := r.Scan(&user.ID, &user.Login, &user.Password, &user.Salt)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, common.ErrorNotFound
			}
			return nil, err
		}

		return r, nil
	})

	return user, err
}

func (r *PostgresRepository) AddUser(ctx context.Context, user *models.User) (models.User, error) {

	s := "insert into users (login, password, salt) values ($1, $2, $3) RETURNING id"

	_, err := common.RetryWithResult(ctx, func() (interface{}, error) {
		err := r.db.QueryRowContext(ctx, s, user.Login, user.Password, user.Salt).Scan(&user.ID)
		return nil, err
	})

	return *user, err
}

func (r *PostgresRepository) FindOrderByNumber(ctx context.Context, number string) (models.Order, error) {

	var order models.Order

	s := "select id, user_id, number, uploaded_at, accrual from orders where number = $1 order by uploaded_at desc"

	_, err := common.RetryWithResult(ctx, func() (*sql.Row, error) {
		r := r.db.QueryRowContext(ctx, s, number)
		err := r.Scan(&order.ID, &order.UserID, &order.Number, &order.UploadedAt, &order.Accrual)

		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, common.ErrorNotFound
			}
			return nil, err
		}

		return r, err
	})
	return order, err
}

func (r *PostgresRepository) AddOrder(ctx context.Context, order *models.Order) (models.Order, error) {

	s := "insert into orders (user_id, number, status) values ($1, $2, $3) RETURNING id"

	_, err := common.RetryWithResult(ctx, func() (interface{}, error) {
		err := r.db.QueryRowContext(ctx, s, order.UserID, order.Number, order.Status).Scan(&order.ID)
		return nil, err
	})

	return *order, err

}

func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {

	s := "select id, user_id, number, uploaded_at, accrual, status from orders where user_id = $1 order by uploaded_at desc"

	rows, err := common.RetryWithResult(ctx, func() (*sql.Rows, error) {
		rows, err := r.db.QueryContext(ctx, s, userID)
		return rows, err
	})

	if err != nil {
		return nil, err
	}

	return r.getOrdersFromRows(rows)
}

func (r *PostgresRepository) getOrdersFromRows(rows *sql.Rows) ([]models.Order, error) {
	var orders = []models.Order{}

	defer rows.Close()
	for rows.Next() {
		var order = models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.UploadedAt, &order.Accrual, &order.Status)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, nil
}

func (r *PostgresRepository) GetUnprocessedOrders(ctx context.Context) ([]models.Order, error) {

	s := "select id, user_id, number, uploaded_at, accrual, status from orders where status in ($1,  $2)"

	rows, err := common.RetryWithResult(ctx, func() (*sql.Rows, error) {
		rows, err := r.db.QueryContext(ctx, s, models.OrderStatusNew, models.OrderStatusProcessing)
		return rows, err
	})

	if err != nil {
		return nil, err
	}

	return r.getOrdersFromRows(rows)

}

func (r *PostgresRepository) UpdateOrderAccrualStatus(ctx context.Context, orderID string, status models.OrderStatus, accrual float32) error {

	s := "update orders set status = $1, accrual = $2 where id = $3"

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := r.db.ExecContext(ctx, s, status, accrual, orderID)
		return res, err
	})

	return err

}

func (r *PostgresRepository) UpdateUserAccruedTotel(ctx context.Context, userID string, amount float32) error {

	s := "update users set accrued_total = $1 where id = $2"

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := r.db.ExecContext(ctx, s, amount, userID)
		return res, err
	})

	return err

}

func (r *PostgresRepository) UpdateUserWithdrawnTotel(ctx context.Context, userID string, amount float32) error {

	s := "update users set withdrawn_total = $1 where id = $2"

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := r.db.ExecContext(ctx, s, amount, userID)
		return res, err
	})

	return err

}

func (r *PostgresRepository) FindUserByID(ctx context.Context, userID string) (models.User, error) {
	s := "select id, login, password, accrued_total, withdrawn_total from users where id=$1"

	var user models.User

	_, err := common.RetryWithResult(ctx, func() (*sql.Row, error) {
		r := r.db.QueryRowContext(ctx, s, userID)
		err := r.Scan(&user.ID, &user.Login, &user.Password, &user.AccruedTotal, &user.WithdrawnTotal)
		return r, err
	})

	return user, err

}

func (r *PostgresRepository) AddWithdrawal(ctx context.Context, item *models.Withdrawal) error {

	s := "insert into withdrawals (user_id, \"order\", amount) values ($1, $2, $3)"

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := r.db.ExecContext(ctx, s, item.UserID, item.Order, item.Amount)
		return res, err
	})

	return err

}

func (r *PostgresRepository) GetWithdrawalsTotalAmountByUserID(ctx context.Context, userID string) (float32, error) {
	s := "select coalesce(sum(amount), 0) from withdrawals where user_id = $1"

	var res float32

	_, err := common.RetryWithResult(ctx, func() (any, error) {
		err := r.db.QueryRowContext(ctx, s, userID).Scan(&res)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, common.ErrorNotFound
			}
			return nil, err
		}
		return nil, err
	})

	return res, err

}

func (r *PostgresRepository) GetAccrualsTotalAmountByUserID(ctx context.Context, userID string) (float32, error) {
	s := "select coalesce(sum(accrual),0) from orders where user_id = $1 and status = $2"

	var res float32

	_, err := common.RetryWithResult(ctx, func() (any, error) {
		err := r.db.QueryRowContext(ctx, s, userID, models.OrderStatusProcessed).Scan(&res)
		if err != nil {
			if errors.Is(err, sql.ErrNoRows) {
				return nil, common.ErrorNotFound
			}
			return nil, err
		}
		return nil, err
	})

	return res, err
}

func (r *PostgresRepository) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]models.Withdrawal, error) {

	s := "select id, user_id, \"order\", uploaded_at, amount from withdrawals where user_id = $1 order by uploaded_at desc"

	rows, err := common.RetryWithResult(ctx, func() (*sql.Rows, error) {
		rows, err := r.db.QueryContext(ctx, s, userID)
		return rows, err
	})

	var withdrawals = []models.Withdrawal{}

	defer rows.Close()
	for rows.Next() {
		var withdrawal = models.Withdrawal{}
		err := rows.Scan(&withdrawal.ID, &withdrawal.UserID, &withdrawal.Order, &withdrawal.UploadedAt, &withdrawal.Amount)
		if err != nil {
			return nil, err
		}
		withdrawals = append(withdrawals, withdrawal)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return withdrawals, err

}
