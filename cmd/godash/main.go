package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/callmemhz/godash/internal/provider"
	"github.com/callmemhz/godash/model"
	"github.com/callmemhz/godash/model/alfred"
	"github.com/samber/lo"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "godash",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "file",
				Value: "",
			},
			&cli.StringFlag{
				Name:  "sqlite",
				Value: "./godash.db",
			},
		},
		Action: run,
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}

func run(cCtx *cli.Context) error {
	var providerFactory provider.ProviderFactory
	if repo := cCtx.String("file"); repo != "" {
		providerFactory = &provider.TxtProviderFactory{RepoPath: repo}
	} else if path := cCtx.String("sqlite"); path != "" {
		providerFactory = &provider.SqliteProviderFactory{Path: path}
	} else {
		return fmt.Errorf("unknown repo type")
	}

	provider, err := providerFactory.NewProvider()
	if err != nil {
		return err
	}

	var namespace, keyword string
	if cCtx.Args().Len() < 2 {
		keyword = cCtx.Args().First()
	} else {
		namespace, keyword = cCtx.Args().Get(0), cCtx.Args().Get(1)
	}
	entries, err := provider.Search(namespace, keyword)
	if err != nil {
		return err
	}

	items := lo.Map(entries, func(entry model.Entry, _ int) alfred.Item {
		return alfred.Item{
			Title:    entry.Title(),
			Subtitle: entry.Subtitle(),
			Arg:      entry.Content(),
		}
	})

	debugMsg := fmt.Sprintf("args: %v, queryNamespace: %v, keyword: %v", os.Args, namespace, keyword)
	res, _ := json.Marshal(alfred.Items{Items: items, DebugMsg: debugMsg})
	fmt.Print(string(res))
	return nil
}
