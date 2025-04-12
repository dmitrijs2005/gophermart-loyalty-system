package repository

import (
	"context"
	"database/sql"
	"errors"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/common"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/models"
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

	goose.SetBaseFS(nil) // default is os.DirFS(".")

	if err := goose.UpContext(ctx, db, "./migrations"); err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (r *PostgresRepository) BeginTransaction(ctx context.Context) error {
	//r.mu.Lock()
	//fmt.Println("BEGIN TRANSACTIOB")
	return nil
}

func (r *PostgresRepository) CommitTransaction(ctx context.Context) error {
	//r.mu.Unlock()
	//fmt.Println("COMMIT TRANSACTIOB")
	return nil
}

func (r *PostgresRepository) RollbackTransaction(ctx context.Context) error {
	//r.mu.Unlock()
	//fmt.Println("ROLLBACK TRANSACTIOB")
	return nil
}

// func (r *PostgresRepository) findUserIDByLogin(_ context.Context, login string) string {
// 	id, exists := r.userLookupByLogin[login]
// 	if !exists {
// 		return ""
// 	}
// 	return id
// }

func (r *PostgresRepository) FindUserByLogin(ctx context.Context, login string) (models.User, error) {

	s := "select id, login, password from users where login=$1"

	exec := r.db

	var user models.User

	_, err := common.RetryWithResult(ctx, func() (*sql.Row, error) {
		r := exec.QueryRowContext(ctx, s, login)
		err := r.Scan(&user.ID, &user.Login, &user.Password)
		return r, err
	})

	return user, err
}

func (r *PostgresRepository) AddUser(ctx context.Context, user *models.User) (models.User, error) {

	s := "insert into users (login, password) values ($1, $2) RETURNING id"

	exec := r.db

	_, err := common.RetryWithResult(ctx, func() (interface{}, error) {
		err := exec.QueryRowContext(ctx, s, user.Login, user.Password).Scan(&user.ID)
		return nil, err
	})

	return *user, err
}

func (r *PostgresRepository) FindOrderByID(ctx context.Context, id string) (models.Order, error) {
	// o, exists := r.orders[id]
	// if !exists {
	// 	return models.Order{}, common.ErrorOrderDoesNotExist
	// }
	// return o, nil

	return models.Order{}, nil

}

func (r *PostgresRepository) FindOrderByNumber(ctx context.Context, number string) (models.Order, error) {

	exec := r.db
	var order models.Order

	s := "select id, user_id, number, uploaded_at, accrual from orders where number = $1 order by uploaded_at desc"

	_, err := common.RetryWithResult(ctx, func() (*sql.Row, error) {
		r := exec.QueryRowContext(ctx, s, number)
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

	exec := r.db

	_, err := common.RetryWithResult(ctx, func() (interface{}, error) {
		err := exec.QueryRowContext(ctx, s, order.UserID, order.Number, order.Status).Scan(&order.ID)
		return nil, err
	})

	return *order, err

}

func (r *PostgresRepository) GetOrdersByUserID(ctx context.Context, userID string) ([]models.Order, error) {

	s := "select id, user_id, number, uploaded_at, accrual, status from orders where user_id = $1 order by uploaded_at desc"

	exec := r.db

	rows, err := common.RetryWithResult(ctx, func() (*sql.Rows, error) {
		rows, err := exec.QueryContext(ctx, s, userID)
		return rows, err
	})

	if err != nil {
		return nil, err
	}

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

	return orders, err
}

func (r *PostgresRepository) GetUnprocessedOrders(ctx context.Context) ([]models.Order, error) {

	s := "select id, user_id, number, uploaded_at, accrual from orders where status in ($1,  $2)"

	exec := r.db

	rows, err := common.RetryWithResult(ctx, func() (*sql.Rows, error) {
		rows, err := exec.QueryContext(ctx, s, models.OrderStatusNew, models.OrderStatusProcessing)
		return rows, err
	})

	var orders = []models.Order{}

	defer rows.Close()
	for rows.Next() {
		var order = models.Order{}
		err := rows.Scan(&order.ID, &order.UserID, &order.Number, &order.UploadedAt, &order.Accrual)
		if err != nil {
			return nil, err
		}
		orders = append(orders, order)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return orders, err

}

func (r *PostgresRepository) UpdateOrderAccrualStatus(ctx context.Context, orderID string, status models.OrderStatus, accrual float32) error {

	s := "update orders set status = $1, accrual = $2 where id = $3"

	exec := r.db

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := exec.ExecContext(ctx, s, status, accrual, orderID)
		return res, err
	})

	return err

}

func (r *PostgresRepository) UpdateUserAccruedTotel(ctx context.Context, userID string, amount float32) error {

	s := "update users set accrued_total = $1 where id = $2"

	exec := r.db

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := exec.ExecContext(ctx, s, amount, userID)
		return res, err
	})

	return err

}

func (r *PostgresRepository) UpdateUserWithdrawnTotel(ctx context.Context, userID string, amount float32) error {

	s := "update users set withdrawn_total = $1 where id = $2"

	exec := r.db

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := exec.ExecContext(ctx, s, amount, userID)
		return res, err
	})

	return err

}

func (r *PostgresRepository) FindUserByID(ctx context.Context, userID string) (models.User, error) {
	s := "select id, login, password, accrued_total, withdrawn_total from users where id=$1"

	exec := r.db

	var user models.User

	_, err := common.RetryWithResult(ctx, func() (*sql.Row, error) {
		r := exec.QueryRowContext(ctx, s, userID)
		err := r.Scan(&user.ID, &user.Login, &user.Password, &user.AccruedTotal, &user.WithdrawnTotal)
		return r, err
	})

	return user, err

}

func (r *PostgresRepository) AddWithdrawal(ctx context.Context, item *models.Withdrawal) error {

	s := "insert into withdrawals (user_id, \"order\", amount) values ($1, $2, $3)"

	exec := r.db

	_, err := common.RetryWithResult(ctx, func() (sql.Result, error) {
		res, err := exec.ExecContext(ctx, s, item.UserID, item.Order, item.Amount)
		return res, err
	})

	return err

}

func (r *PostgresRepository) GetWithdrawalsByUserID(ctx context.Context, userID string) ([]models.Withdrawal, error) {

	s := "select id, user_id, \"order\", uploaded_at, amount from withdrawals where user_id = $1 order by uploaded_at desc"

	exec := r.db

	rows, err := common.RetryWithResult(ctx, func() (*sql.Rows, error) {
		rows, err := exec.QueryContext(ctx, s, userID)
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
