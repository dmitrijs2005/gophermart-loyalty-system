package common

import (
	"context"
	"errors"
	"net"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func CheckLuhn(s string) (bool, error) {
	ctr := 0

	double := false
	for i := len(s) - 1; i >= 0; i-- {
		n, err := strconv.Atoi(string(s[i]))
		if err != nil {
			return false, err
		}

		add := n
		if double {
			add = add * 2
			if add >= 10 {
				add -= 9
			}
		}
		ctr += add
		double = !double
	}
	return ctr%10 == 0, nil
}

func CheckForAllDigits(s string) bool {
	pattern := `^[0-9]*$`

	matched, err := regexp.MatchString(pattern, s)
	if err != nil {
		return false
	}

	return matched
}

func CheckOrderNumberFormat(number string) (bool, error) {
	if number == "" {
		return false, nil
	}

	if !CheckForAllDigits(number) {
		return false, nil
	}

	valid, err := CheckLuhn(number)

	if err != nil {
		return false, err
	}

	if !valid {
		return false, err
	}

	return true, nil
}

var (
	MaxRetries int           = 3
	ExpBackoff time.Duration = 2 * time.Second
)

// tries to detect if there is a point to retry
func isErrorRetriable(err error) bool {

	if err == nil {
		return false
	}

	// checking if syscall error (connect refused)
	var sysErr *os.SyscallError
	if errors.As(err, &sysErr) {
		return true
	}

	// connection timeout
	if errors.Is(err, os.ErrDeadlineExceeded) {
		return true
	}

	// Проверка PostgreSQL спец. ошибок
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) {
		// Class 08 - Connection Exception
		if pgerrcode.IsConnectionException(pgErr.Code) {
			return true
		}
	}

	// checking if error is network-related
	var netErr *net.OpError
	if errors.As(err, &netErr) && netErr.Temporary() {
		return true
	}

	return false

}

// function that processes retries
func RetryWithResult[T any](ctx context.Context, request func() (T, error)) (T, error) {
	var result T
	var err error

	for i := 0; i < 1+MaxRetries; i++ {

		result, err = request()

		if err == nil {
			return result, nil
		}

		if !isErrorRetriable(err) {
			return result, err
		}

		backoff := 1*time.Second + time.Duration(i)*ExpBackoff

		select {
		case <-time.After(backoff):
		case <-ctx.Done():
			return result, ctx.Err() // Context timeout/cancellation
		}

	}

	return result, err
}

func FilterMap[T any](m map[string]T, predicate func(T) bool) []T {
	var result []T
	for _, item := range m {
		if predicate(item) {
			result = append(result, item)
		}
	}
	return result
}
