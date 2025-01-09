package main

import (
	"context"
	"github.com/urfave/cli/v3"
	"log"
	"os"
	"static-openapivalidator/internal"
)

const (
	specFlagName        = "spec"
	reportFlagName      = "report"
	formatFlagName      = "format"
	reportHTMLFlagName  = "report-html"
	reportJUNITFlagName = "report-junit"
	reportJSONFlagName  = "report-json"
)

func main() {
	cmd := &cli.Command{
		Name:  "static-openapivalidator",
		Usage: "Check openapi against static results",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     specFlagName,
				Aliases:  []string{"s"},
				Usage:    "Load openapi spec from `FILE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:     reportFlagName,
				Aliases:  []string{"r"},
				Usage:    "Load report from `FILE`",
				Required: true,
			},
			&cli.StringFlag{
				Name:    formatFlagName,
				Aliases: []string{"f"},
				Value:   "bruno",
				Usage:   "Use report format `FORMAT`",
			},
			&cli.StringFlag{
				Name:  reportHTMLFlagName,
				Usage: "Export HTML report to `FILE`",
			},
			&cli.StringFlag{
				Name:  reportJUNITFlagName,
				Usage: "Export JSON report to `FILE`",
			},
			&cli.StringFlag{
				Name:  reportJSONFlagName,
				Usage: "Export JUNIT report to `FILE`",
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			return internal.Params{
				Ctx:            ctx,
				ApiFilePath:    cmd.String(specFlagName),
				ReportFilePath: cmd.String(reportFlagName),
				Format:         cmd.String(formatFlagName),
				JunitFilePath:  cmd.String(reportJUNITFlagName),
				HtmlFilePath:   cmd.String(reportHTMLFlagName),
				JsonFilePath:   cmd.String(reportJSONFlagName),
			}.Execute()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
