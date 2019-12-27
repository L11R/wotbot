package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/L11R/wotbot/internal/infra/database"

	"github.com/L11R/wotbot/internal/configs"
	"github.com/L11R/wotbot/internal/domain"
	"github.com/L11R/wotbot/internal/infra/telegram"
	"github.com/L11R/wotbot/internal/infra/wargaming"
	"github.com/L11R/wotbot/internal/infra/xvm"
	"github.com/jessevdk/go-flags"
	"go.uber.org/zap"
)

func main() {
	// Parse command line arguments or environment variables
	config, err := configs.Parse()
	if err != nil {
		if err, ok := err.(*flags.Error); ok {
			fmt.Println(err)
			os.Exit(0)
		}

		fmt.Printf("Invalid args: %v\n", err)
		os.Exit(1)
	}

	// Init logger
	logger, err := zap.NewProduction()
	if len(config.Verbose) != 0 && config.Verbose[0] {
		logger, err = zap.NewDevelopment()
	}
	if err != nil {
		log.Fatalln(err)
	}

	db, err := database.NewAdapter(logger, config.Database)
	if err != nil {
		logger.Fatal("Error creating new database adapter!", zap.Error(err))
	}
	ws := wargaming.NewAdapter(logger, config.Wargaming)
	x := xvm.NewAdapter(logger, config.XVM)

	service := domain.NewService(logger, db, ws, x)

	ts, err := telegram.NewAdapter(logger, config.Telegram, service)
	if err != nil {
		logger.Panic("Error creating new Telegram adapter!", zap.Error(err))
	}

	shutdown := make(chan error, 1)

	go func(shutdown chan<- error) {
		shutdown <- ts.ListenAndServe()
	}(shutdown)

	// Graceful shutdown block
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)

	select {
	case s := <-sig:
		logger.Info("Got the signal!", zap.Any("signal", s))
	case err := <-shutdown:
		logger.Error("Error running the application!", zap.Error(err))
	}

	logger.Info("Stopping bot...")
	ts.Shutdown()
	logger.Info("Bot stopped")
}
