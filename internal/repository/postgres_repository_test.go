package repository

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

func runMigrations(t *testing.T, db *sql.DB, migrationsDir string) {
	t.Helper()

	err := goose.SetDialect("postgres")
	require.NoError(t, err)

	err = goose.Up(db, migrationsDir)
	require.NoError(t, err)
}

func findMigrationsDir() (string, error) {
	wd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// Сначала попробуем "migrations" в корне проекта
	path := filepath.Join(wd, "migrations")
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	// Или попробуем рядом с исполняемым файлом
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	exeDir := filepath.Dir(exePath)
	path = filepath.Join(exeDir, "migrations")
	if _, err := os.Stat(path); err == nil {
		return path, nil
	}

	return "", fmt.Errorf("migrations folder not found")
}

func TestRepository(t *testing.T) {
	ctx := context.Background()

	dbname := "test"
	user := "test"
	password := "test"

	// Start Postgres container
	// 1. Start the postgres ctr and run any migrations on it
	ctr, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithDatabase(dbname),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
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

}

// 	require.NoError(t, err)
// 	defer pgContainer.Terminate(ctx)

// 	dbURI, err := pgContainer.ConnectionString(ctx, "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable")
// 	require.NoError(t, err)

// 	//dbURI := "postgresql://postgres:postgres@postgres/praktikum?sslmode=disable"

// 	db, err := sql.Open("pgx", dbURI)
// 	require.NoError(t, err)

// 	err = db.PingContext(ctx)
// 	require.NoError(t, err)

// 	// Create test table
// 	_, err = db.ExecContext(ctx, `
//         CREATE TABLE users (
//             id SERIAL PRIMARY KEY,
//             name TEXT NOT NULL,
//             email TEXT NOT NULL UNIQUE
//         );
//     `)
// 	require.NoError(t, err)

// 	// Now test the repository
// 	_, err = NewPostgresRepository(ctx, dbURI)
// 	require.NoError(t, err)

// 	// user := &repository.User{
// 	// 	Name:  "Alice",
// 	// 	Email: "alice@example.com",
// 	// }

// 	// err = userRepo.Create(ctx, user)
// 	// require.NoError(t, err)
// 	// require.NotZero(t, user.ID)

// 	// foundUser, err := userRepo.GetByID(ctx, user.ID)
// 	// require.NoError(t, err)
// 	// require.Equal(t, user.Name, foundUser.Name)
// 	// require.Equal(t, user.Email, foundUser.Email)
// }
