package main

import (
	"context"
	"errors"
	"github.com/quangnc2k/do-an-gang/internal/app"
	"github.com/quangnc2k/do-an-gang/pkg/env"
	"github.com/urfave/cli/v2"
	"log"
	"os"
	"sort"
	"time"
)

func Migrate() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		ctx, ok := c.App.Metadata["context"].(context.Context)
		if !ok {
			return errors.New("invalid root context")
		}

		return app.Migrate(ctx)
	}
}

func Run() func(c *cli.Context) error {
	return func(c *cli.Context) error {
		//ctx, ok := c.App.Metadata["context"].(context.Context)
		//if !ok {
		//	return errors.New("invalid root context")
		//}

		ctx := context.Background()

		return app.Process(ctx)
	}
}

func main() {
	a := &cli.App{}
	a.Name = "Do An"
	a.Usage = "Post process"
	a.Version = "1.0.0"
	a.Compiled = time.Now()
	a.Commands = []*cli.Command{
		{
			Name:   "run",
			Usage:  "fetching events from queue, serve backend",
			Action: Run(),
		},

		{
			Name:   "migrate",
			Usage:  "create database",
			Action: Migrate(),
		},
	}
	a.Flags = []cli.Flag{
		&cli.StringFlag{
			Name:    "env",
			Aliases: []string{"e"},
			Value:   "./.env",
			Usage:   "set path to environment file",
		},
	}

	a.Before = func(c *cli.Context) error {
		err := env.LoadEnv(c.String("env"))
		if err != nil {
			return cli.Exit(err.Error(), 1)
		}
		ctx := context.Background()

		a.Metadata["context"] = ctx

		return nil
	}

	sort.Sort(cli.FlagsByName(a.Flags))
	sort.Sort(cli.CommandsByName(a.Commands))

	err := a.Run(os.Args)
	if err == nil {
		return
	}

	log.Println(err)
	os.Exit(1)
}
