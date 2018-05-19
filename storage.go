package main

import (
	cloudStorage "cloud.google.com/go/storage"
	"context"
	"firebase.google.com/go/storage"
	"fmt"
	"github.com/olekukonko/tablewriter"
	"github.com/urfave/cli"
	"os"
	"strings"
)

type BucketAttrs cloudStorage.BucketAttrs

func (b *BucketAttrs) String() string {
	return fmt.Sprintf(
		"Name: "+b.Name+"\n"+
			"Location: "+b.Location+"\n"+
			"Class: "+b.StorageClass+"\n"+
			"Created: %v"+
			"", b.Created)
}

var st *storage.Client

func setupStorage(ctx *cli.Context) (err error) {
	st, err = fb.Storage(context.Background())
	if err != nil {
		return err
	}

	return nil
}

func storagePing(ctx *cli.Context) error {
	bucketname := ctx.String("bucket")
	if bucketname == "" {
		bucketname = ctx.GlobalString("bucket")
	}

	bucket, err := st.Bucket(bucketname)
	if err != nil {
		return err
	}

	attrs, err := bucket.Attrs(context.Background())
	if err != nil {
		return err
	}

	printableAttrs := BucketAttrs(*attrs)

	fmt.Println(printableAttrs.String())
	return nil
}

func storageList(ctx *cli.Context) error {
	const timeFormat = "2006/01/02 15:04"

	bucket, err := st.Bucket(ctx.GlobalString("bucket"))
	if err != nil {
		return err
	}

	objectIt := bucket.Objects(context.Background(), nil)

	table := tablewriter.NewWriter(os.Stdout)
	table.SetBorder(false)
	table.SetHeader([]string{"Name", "CacheControl", "Created", "Updated"})

	for object, err := objectIt.Next(); object != nil && err == nil; object, err = objectIt.Next() {
		table.Append([]string{
			object.Prefix + object.Name,
			object.CacheControl,
			object.Created.Format(timeFormat),
			object.Updated.Format(timeFormat),
		})
	}

	table.Render()

	return nil
}

func storageSetCache(ctx *cli.Context) error {
	bucket, err := st.Bucket(ctx.GlobalString("bucket"))
	if err != nil {
		return err
	}

	_, err = bucket.Attrs(context.Background())
	if err != nil {
		return err
	}

	cache := ctx.String("cache")
	if cache == "" || !strings.Contains(cache, "max-age") {
		return fmt.Errorf(`invalid cache format, should be "[public|private, ]max-age=<seconds>"`)
	}

	suffixes := ctx.StringSlice("suffix")

	if len(suffixes) <= 0 {
		suffixes = []string{"mp4", "webm", "png", "jpg", "jpeg", "bin"}
	}

	objectIt := bucket.Objects(context.Background(), nil)

	for object, err := objectIt.Next(); object != nil && err == nil; object, err = objectIt.Next() {
		for _, s := range suffixes {
			if strings.HasSuffix(object.Name, s) {
				if object.CacheControl == cache {
					fmt.Println("Cache-Control already set for " + object.Prefix + object.Name)
				} else {
					fmt.Println("Setting Cache-Control to " + cache + " for " + object.Prefix + object.Name)
					_, err = bucket.Object(object.Name).Update(context.Background(), cloudStorage.ObjectAttrsToUpdate{CacheControl: cache})
				}
				break
			}
		}
	}

	return err
}
