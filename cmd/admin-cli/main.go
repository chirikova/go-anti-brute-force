package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"antibruteforce/internal/transport/grpc"
	"antibruteforce/internal/transport/grpc/api"
	"github.com/urfave/cli"
)

func main() {
	app := &cli.App{
		Name: "CLI admin panel",
		Usage: `интерфейс для ручного администрирования сервиса. 
		Через CLI есть возможность вызвать сброс бакета и управлять whitelist/blacklist-ом. 
		CLI рабоатет через GRPC интерфейс`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "address",
				Value: "localhost:9095",
				Usage: "Anti-bruteforce server address? e.g.: 127.0.0.1:9091",
			},
			&cli.DurationFlag{
				Name:  "connect-timeout",
				Value: 3 * time.Second,
				Usage: "Anti-bruteforce client connection timeout",
			},
		},
		Commands: []cli.Command{
			{
				Name:    "Reset",
				Aliases: []string{"reset"},
				Usage:   "Reset bucket by login and ip",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "login",
						Usage: "Login",
					},
					&cli.StringFlag{
						Name:  "ip",
						Usage: "Ip",
					},
				},
				Action: func(c *cli.Context) error {
					ctx, cancel := context.WithTimeout(context.Background(), c.Duration("connect-timeout"))
					defer cancel()

					client, err := grpc.NewClient(ctx, c.String("address"))
					if err != nil {
						log.Fatalf("unable to establish connection: %s", err)
					}

					if response, err := client.Reset(ctx, &api.ResetRequest{}); err != nil {
						log.Println("Error: ", err)
					} else {
						log.Printf("Response: %v\n", response)
					}

					go func() {
						terminate := make(chan os.Signal, 1)
						signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
						<-terminate
						log.Println("Received system interrupt...")
						cancel()
					}()

					return nil
				},
			},
			{
				Name:    "Auth",
				Aliases: []string{"auth"},
				Usage:   "Authorization attempt",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "Add to whitelist",
				Aliases: []string{"add2whitelist"},
				Usage:   "add subnet(ip + mask) to the whitelist (e.g.: 255.0.0.0/12)",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "remove_from_whitelist",
				Aliases: []string{"rw"},
				Usage:   "remove subnet(ip + mask) from the whitelist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "add_in_blacklist",
				Aliases: []string{"ab"},
				Usage:   "add subnet(ip + mask) in the blacklist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
			{
				Name:    "remove_from_blacklist",
				Aliases: []string{"rb"},
				Usage:   "remove subnet(ip + mask) from the blacklist. For example: 192.168.130.0/24",
				Action: func(c *cli.Context) error {
					return nil
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

//type clientResponse func(client api.ApiServiceClient) (string, error)
//
//func runCommand(clientResponse clientResponse, cfg config.Config) error {
//	ctx, cancel := context.WithTimeout(context.Background(), cfg.GRPC.Client.ConnectionTimeout)
//	defer cancel()
//
//	client, err := grpc.NewClient(ctx, cfg.GRPC)
//	if err != nil {
//		log.Fatalf("unable to establish connection: %s", err)
//	}
//
//	if response, err := clientResponse(client); err != nil {
//		log.Println("Error: ", err)
//	} else {
//		log.Printf("Response: %v\n", response)
//	}
//
//	go func() {
//		terminate := make(chan os.Signal, 1)
//		signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
//		<-terminate
//		log.Println("Received system interrupt...")
//		cancel()
//	}()
//
//	return nil
//}
