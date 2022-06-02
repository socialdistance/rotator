package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	internalapp "rotator/internal/app"
	internalconfig "rotator/internal/config"
	internallogger "rotator/internal/logger"
	"rotator/internal/rq"
	internalhttp "rotator/internal/server/http"
	internalstore "rotator/internal/storage/store"
	"syscall"
	"time"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.json", "Path to configuration file")
}

func main() {
	config, err := internalconfig.LoadConfig(configFile)
	if err != nil {
		log.Fatalf("Failed load config %s", err)
	}

	logger, err := internallogger.NewLogger(config.Logger)
	if err != nil {
		log.Fatalf("Failed load config %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	store := internalstore.CreateStorage(ctx, *config)
	logger.Info("[+] Connected to database")

	application := internalapp.New(logger, store)

	server := internalhttp.NewServer(config.HTTP.Host, config.HTTP.Port, application, logger)

	// reciver
	_, err = rq.NewRabbit(ctx, config.Rabbit.Url, config.Rabbit.Exchange, config.Rabbit.Queue, logger)
	if err != nil {
		logger.Error("[-] Failed to start RabbitMQ")
	}
	logger.Info("[+] Connected to RabbitMQ")

	go func() {
		<-ctx.Done()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logger.Error("failed to stop http server: " + err.Error())
		}
	}()

	logger.Info("[+] Application starting...")

	if err := server.Start(ctx); err != nil {
		logger.Error("failed to start http server: " + err.Error())
		cancel()
		os.Exit(1) //nolint:gocritic
	}
}
