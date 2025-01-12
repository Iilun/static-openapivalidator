package main

import (
	"context"
	altsrc "github.com/urfave/cli-altsrc/v3"
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
	configFileFlagName  = "config-file"
)

func main() {

	var configFiles []string
	// TODO: improve this
	configFilePath := os.Getenv("CONFIG_FILE")
	if configFilePath != "" {
		configFiles = []string{configFilePath}
	}

	cmd := &cli.Command{
		Name:  "static-openapivalidator",
		Usage: "Check openapi against static results",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     specFlagName,
				Aliases:  []string{"s"},
				Usage:    "Load openapi spec from `FILE`",
				Required: true,
				Sources:  altsrc.YAML(specFlagName, configFiles...),
			},
			&cli.StringSliceFlag{
				Name:     reportFlagName,
				Aliases:  []string{"r"},
				Usage:    "Load report from `FILES`",
				Required: true,
				Sources:  altsrc.YAML(reportFlagName, configFiles...),
			},
			&cli.StringFlag{
				Name:    formatFlagName,
				Aliases: []string{"f"},
				Value:   "bruno",
				Usage:   "Use report format `FORMAT`",
				Sources: altsrc.YAML(formatFlagName, configFiles...),
			},
			&cli.StringFlag{
				Name:    reportHTMLFlagName,
				Usage:   "Export HTML report to `FILE`",
				Sources: altsrc.YAML(reportHTMLFlagName, configFiles...),
			},
			&cli.StringFlag{
				Name:    reportJUNITFlagName,
				Usage:   "Export JUNIT report to `FILE`",
				Sources: altsrc.YAML(reportJUNITFlagName, configFiles...),
			},
			&cli.StringFlag{
				Name:    reportJSONFlagName,
				Usage:   "Export JSON report to `FILE`",
				Sources: altsrc.YAML(reportJSONFlagName, configFiles...),
			},
			&cli.StringFlag{
				Name:    configFileFlagName,
				Usage:   "Export JSON report to `FILE`",
				Sources: altsrc.YAML(reportJSONFlagName, configFiles...),
			},
		},
		Action: func(ctx context.Context, cmd *cli.Command) error {
			params := internal.Params{
				Ctx:             ctx,
				ApiFilePath:     cmd.String(specFlagName),
				ReportFilePaths: cmd.StringSlice(reportFlagName),
				Format:          cmd.String(formatFlagName),
				JunitFilePath:   cmd.String(reportJUNITFlagName),
				HtmlFilePath:    cmd.String(reportHTMLFlagName),
				JsonFilePath:    cmd.String(reportJSONFlagName),
				ConfigFilePath:  configFilePath,
			}
			return params.Execute()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		log.Fatal(err)
	}
}
