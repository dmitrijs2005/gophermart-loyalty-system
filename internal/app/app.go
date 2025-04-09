package app

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/config"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/repository"
	"github.com/dmitrijs2005/gophermart-loyalty-system/internal/server"
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

func (app *App) startHTTPServer(ctx context.Context, cancelFunc context.CancelFunc, wg *sync.WaitGroup, repository repository.Repository) {

	wg.Add(1)

	go func() {
		defer wg.Done()
		server, err := server.NewHTTPServer(app.config, repository)
		if err != nil {
			cancelFunc()
		}
		if err := server.Run(ctx); err != nil {
			cancelFunc()
		}
	}()
}

func (app *App) startCheckingTask(ctx context.Context, wg *sync.WaitGroup, repository repository.Repository) {

	wg.Add(1)

	go func() {
		defer wg.Done()
		task := task.NewAccrualCheckerTask(app.config, repository)
		task.Start(ctx)
	}()
}

func (app *App) Run() error {

	ctx, cancelFunc := context.WithCancel(context.Background())

	app.initSignalHandler(cancelFunc)

	repository, err := repository.NewInMemoryRepository()
	if err != nil {
		return err
	}

	var wg sync.WaitGroup

	app.startHTTPServer(ctx, cancelFunc, &wg, repository)
	app.startCheckingTask(ctx, &wg, repository)

	wg.Wait()

	return nil

}
