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
	"github.com/urfave/cli/v2"
)

var address string
var timeout time.Duration

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	client, err := grpc.NewClient(ctx, address)
	if err != nil {
		log.Fatalf("unable to establish connection: %s", err)
	}

	app := &cli.App{
		Name: "CLI admin panel",
		Usage: `интерфейс для ручного администрирования сервиса. 
		Через CLI есть возможность вызвать сброс бакета и управлять whitelist/blacklist-ом. 
		CLI работает через GRPC интерфейс`,
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:        "address",
				Value:       "localhost:9095",
				Usage:       "Anti-bruteforce server address? e.g.: 127.0.0.1:9091",
				Destination: &address,
			},
			&cli.DurationFlag{
				Name:        "connect-timeout",
				Value:       300 * time.Second,
				Usage:       "Anti-bruteforce client connection timeout",
				Destination: &timeout,
			},
		},
		Commands: []*cli.Command{
			{
				Name:    "Reset",
				Aliases: []string{"reset"},
				Usage:   "Reset bucket by login and ip",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "login",
						Usage:    "Login, e.g.: sidorov",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					r := api.ResetRequest{Login: c.String("login"), Ip: c.String("ip")}
					if response, err := client.Reset(ctx, &r); err != nil {
						log.Println("Error: ", err)
						return err
					} else {
						log.Printf("Response: %v\n", response)
					}

					return nil
				},
			},
			{
				Name:    "Add to whitelist",
				Aliases: []string{"aw"},
				Usage:   "add subnet(ip + mask) to the whitelist (e.g.: 255.0.0.0/12)",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					r := api.WhitelistAddRequest{
						SubNet: &api.SubNet{
							Ip:   c.String("ip"),
							Mask: c.String("mask"),
						},
					}
					if response, err := client.WhitelistAdd(ctx, &r); err != nil {
						log.Println("Error: ", err)
						return err
					} else {
						log.Printf("Response: %v\n", response)
					}

					return nil
				},
			},
			{
				Name:    "Remove from whitelist",
				Aliases: []string{"rw"},
				Usage:   "remove subnet(ip + mask) from the whitelist. For example: 192.168.130.0/24",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					r := api.WhitelistRemoveRequest{
						SubNet: &api.SubNet{
							Ip:   c.String("ip"),
							Mask: c.String("mask"),
						},
					}
					if response, err := client.WhitelistRemove(ctx, &r); err != nil {
						log.Println("Error: ", err)
						return err
					} else {
						log.Printf("Response: %v\n", response)
					}

					return nil
				},
			},
			{
				Name:    "Add to blacklist",
				Aliases: []string{"ab"},
				Usage:   "add subnet(ip + mask) in the blacklist. For example: 192.168.130.0/24",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					r := api.BlacklistAddRequest{
						SubNet: &api.SubNet{
							Ip:   c.String("ip"),
							Mask: c.String("mask"),
						},
					}
					if response, err := client.BlacklistAdd(ctx, &r); err != nil {
						log.Println("Error: ", err)
						return err
					} else {
						log.Printf("Response: %v\n", response)
					}

					return nil
				},
			},
			{
				Name:    "Remove from blacklist",
				Aliases: []string{"rb"},
				Usage:   "remove subnet(ip + mask) from the blacklist. For example: 192.168.130.0/24",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
					&cli.StringFlag{
						Name:     "ip",
						Usage:    "Ip, e.g.: 192.168.0.1",
						Required: true,
					},
				},
				Action: func(c *cli.Context) error {
					r := api.BlacklistRemoveRequest{
						SubNet: &api.SubNet{
							Ip:   c.String("ip"),
							Mask: c.String("mask"),
						},
					}
					if response, err := client.BlacklistRemove(ctx, &r); err != nil {
						log.Println("Error: ", err)
						return err
					} else {
						log.Printf("Response: %v\n", response)
					}

					return nil
				},
			},
		},
	}

	err = app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}

	terminate := make(chan os.Signal, 1)
	signal.Notify(terminate, os.Interrupt, syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	go func() {
		<-terminate

		log.Println("Received system interrupt...")
		cancel()
	}()
}
