package main

import (
	"antibruteforce/internal/app"
	"antibruteforce/internal/config"
	"antibruteforce/internal/logger"
	"antibruteforce/internal/transport/grpc"
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
}

func main() {
	cfg, err := config.InitConfig(configFile)
	if err != nil {
		log.Fatalf("parsing config file %s: %s", configFile, err)
	}

	file, err := logger.GetLogFile(cfg.Logger.OutputPath)
	if err != nil {
		log.Fatalf("error opening file to log %s", err)
	}

	logg, err := logger.New(cfg.Logger, file)
	if err != nil {
		log.Fatalf("init logger: %s", err)
	}

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)

	antibruteforceApp := app.NewApp(ctx, logg)

	grpcServer := grpc.NewServer(ctx, cfg.GRPC, logg, &antibruteforceApp)

	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			log.Fatalf("error closing log file %s", err)
		}
	}(file)

	wg := sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()

		if err = grpcServer.Start(); err != nil {
			logg.Error("failed to start grpc server: " + err.Error())
			cancel()
		}
	}()

	logg.Info("app is running...")

	go func() {
		defer wg.Done()
		<-ctx.Done()

		if err = grpcServer.Stop(); err != nil {
			logg.Error("failed to stop grpc server: " + err.Error())
		}

		logg.Info("app has stopped...")
	}()

	wg.Wait()
}
