package main

import (
	"github.com/urfave/cli"
	"fmt"
	"firebase.google.com/go"
	"google.golang.org/api/option"
	"context"
	"os"
	"time"
)

func main() {
	app := cli.NewApp()
	app.Name = "pscli"
	app.Description = "Manage PhysiqSquad storage, database, users, and settings"
	app.Version = "0.1.0"

	app.Before = setup

	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:  "config, c",
			Usage: "Path to firebase json config file",
			Value: "firebase.json",
		},
	}

	app.Commands = []cli.Command{
		{
			Name:        "ping",
			Description: "Check connection",
			Action:      ping,
		},
		{
			Name:        "storage",
			Description: "Manipulate firebase storage",
			Before:      setupStorage,
			Action:      storagePing,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:  "bucket, b",
					Usage: "Name of the bucket to use",
				},
			},
			Subcommands: []cli.Command{
				{
					Name:        "show",
					Description: "Show bucket attributes",
					Action:      storagePing,
				},
				{
					Name:        "list",
					Description: "List bucket contents",
					Action:      storageList,
				},
				{
					Name:        "set-cache",
					Description: "Set storage cache by wildcard",
					Action:      storageSetCache,
					Flags: []cli.Flag {
						cli.StringSliceFlag {
							Name:  "suffix, s",
							Usage: "Suffix to use as filter. Can be specified multiple times",
						},
						cli.StringFlag {
							Name:  "cache, c",
							Usage: "String to use as cache",
							Value: fmt.Sprintf("public, max-age=%d", 15 * 24 * time.Hour / time.Second),
						},
					},
				},
			},
		},
		{
			Name:        "auth",
			Description: "Manipulate firebase users",
			Before:      setupUsers,
			Action:      usersList,
			Subcommands: []cli.Command{
				{
					Name:        "list",
					Description: "Show bucket attributes",
					Action:      usersList,
				},
				{
					Name:        "update",
					Description: "Update a user's properties",
					Action:      userUpdate,
					Flags: []cli.Flag {
						cli.StringFlag {
							Name:  "uid, u",
							Usage: "ID of the user to update",
						},
						cli.StringFlag {
							Name:  "name",
						},
						cli.StringFlag {
							Name:  "email",
						},
						cli.BoolFlag {
							Name:  "disabled",
						},
						cli.BoolFlag {
							Name:  "verified",
						},
						cli.StringSliceFlag {
							Name:  "claim, c",
							Usage: "Custom user claims, key=value. Can be specified multiple times",
						},
					},
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
	}
}

var fb *firebase.App

func setup(ctx *cli.Context) (err error) {
	opt := option.WithCredentialsFile(ctx.GlobalString("config"))
	fb, err = firebase.NewApp(context.Background(), nil, opt)
	return err
}

func ping(ctx *cli.Context) error {
	fmt.Println("Connection succeeded")
	return nil
}
