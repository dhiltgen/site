package main

import (
	"fmt"
	"os"

	"github.com/dhiltgen/site/static"

	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func main() {
	app := cli.NewApp()
	app.Name = "site"
	app.Usage = "Run the web site"
	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:   "debug",
			EnvVar: "DEBUG",
			Usage:  "Turn on debug logging to troubleshoot",
		},
		cli.IntFlag{
			Name:   "port",
			Value:  8080,
			EnvVar: "PORT",
			Usage:  "specify the port to bind the server too",
		},
		cli.StringFlag{
			Name:   "site-prefix",
			Value:  "/",
			EnvVar: "SITE_PREFIX",
			Usage:  "prefix all routes with this path",
		},
		cli.StringFlag{
			Name:   "templates-dir",
			Value:  "/templates",
			EnvVar: "TEMPLATES_DIR",
			Usage:  "Location to find template files",
		},
		// TODO - Is this useful or should blobs just be served directly?
		cli.StringFlag{
			Name:   "blob-dir",
			Value:  "/blobs",
			EnvVar: "BLOB_DIR",
			Usage:  "Local file path for blob content",
		},
		cli.StringFlag{
			Name:   "blob-prefix",
			Value:  "/blobs",
			EnvVar: "BLOB_PREFIX",
			Usage:  "Site path for blob access",
		},
	}
	app.Action = func(c *cli.Context) error {
		if c.Bool("debug") {
			log.SetLevel(log.DebugLevel)
		}
		srv := &static.Server{
			Port:         c.Int("port"),
			SitePrefix:   c.String("site-prefix"),
			TemplatesDir: c.String("templates-dir"),
			BlobDir:      c.String("blob-dir"),
			BlobPrefix:   c.String("blob-prefix"),
		}
		return srv.Serve()
	}
	err := app.Run(os.Args)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

}
