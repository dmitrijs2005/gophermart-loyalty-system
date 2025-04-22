package app

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/logging"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/server"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/service"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/task"
)

type App struct {
	config *config.Config
}

func NewApp() *App {

	config, err := config.ParseConfig()
	if err != nil {
		panic(err)
	}

	return &App{config: config}
}

func (app *App) initSignalHandler(cancelFunc context.CancelFunc) {

	// Channel to catch OS signals.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		cancelFunc()
	}()
}

func (app *App) initRepository(ctx context.Context) (repository.Repository, error) {

	var repo repository.Repository
	var err error

	if app.config.DatabaseURI == "" {

		repo, err = repository.NewInMemoryRepository()

	} else {

		repo, err = repository.NewPostgresRepository(ctx, app.config.DatabaseURI)
		if err != nil {
			return nil, err
		}

	}

	if err != nil {
		return nil, err
	}

	dbRepo, ok := repo.(repository.DBStorage)
	if ok {
		err := dbRepo.RunMigrations(ctx)
		if err != nil {
			return nil, err
		}
	}

	return repo, nil

}

func (app *App) startHTTPServer(ctx context.Context, cancelFunc context.CancelFunc,
	wg *sync.WaitGroup, serviceProvider *service.ServiceProvider, logger *slog.Logger) {

	wg.Add(1)

	go func() {
		defer wg.Done()
		server, err := server.NewHTTPServer(app.config, serviceProvider, logger)
		if err != nil {
			cancelFunc()
		}
		if err := server.Run(ctx); err != nil {
			cancelFunc()
		}
	}()
}

func (app *App) startCheckingTask(ctx context.Context, wg *sync.WaitGroup,
	serviceProvider *service.ServiceProvider, logger *slog.Logger) {

	wg.Add(1)

	go func() {
		defer wg.Done()
		task := task.NewAccrualCheckerTask(app.config, serviceProvider.BalanceService, logger)
		task.Start(ctx)
	}()
}

func (app *App) Run() error {

	ctx, cancelFunc := context.WithCancel(context.Background())

	app.initSignalHandler(cancelFunc)

	repository, err := app.initRepository(ctx)
	if err != nil {
		return err
	}

	logger := logging.NewLogger()

	serviceProvider := service.NewServiceProvider(repository, app.config, logger)

	var wg sync.WaitGroup

	app.startHTTPServer(ctx, cancelFunc, &wg, serviceProvider, logger)
	app.startCheckingTask(ctx, &wg, serviceProvider, logger)

	wg.Wait()

	return nil

}
