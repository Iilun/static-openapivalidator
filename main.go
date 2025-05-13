package main

import (
	"context"
	"fmt"
	altsrc "github.com/urfave/cli-altsrc/v3"
	"github.com/urfave/cli-altsrc/v3/yaml"
	"github.com/urfave/cli/v3"
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
	debugFlagName       = "debug"
)

func main() {
	// TODO: improve this
	configFilePath := os.Getenv("CONFIG_FILE")

	cmd := &cli.Command{
		Name:  "static-openapivalidator",
		Usage: "Check openapi against static results",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:     specFlagName,
				Aliases:  []string{"s"},
				Usage:    "Load openapi spec from `FILE`",
				Required: true,
				Sources:  cli.NewValueSourceChain(yaml.YAML(specFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.StringSliceFlag{
				Name:     reportFlagName,
				Aliases:  []string{"r"},
				Usage:    "Load report from `FILES`",
				Required: true,
				Sources:  cli.NewValueSourceChain(yaml.YAML(reportFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.StringFlag{
				Name:    formatFlagName,
				Aliases: []string{"f"},
				Value:   "bruno",
				Usage:   "Use report format `FORMAT`",
				Sources: cli.NewValueSourceChain(yaml.YAML(formatFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.StringFlag{
				Name:    reportHTMLFlagName,
				Usage:   "Export HTML report to `FILE`",
				Sources: cli.NewValueSourceChain(yaml.YAML(reportHTMLFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.StringFlag{
				Name:    reportJUNITFlagName,
				Usage:   "Export JUNIT report to `FILE`",
				Sources: cli.NewValueSourceChain(yaml.YAML(reportJUNITFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.StringFlag{
				Name:    reportJSONFlagName,
				Usage:   "Export JSON report to `FILE`",
				Sources: cli.NewValueSourceChain(yaml.YAML(reportJSONFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.StringFlag{
				Name:    configFileFlagName,
				Usage:   "Export JSON report to `FILE`",
				Sources: cli.NewValueSourceChain(yaml.YAML(configFileFlagName, altsrc.StringSourcer(configFilePath))),
			},
			&cli.BoolFlag{
				Name:    debugFlagName,
				Usage:   "Enable debug logging",
				Sources: cli.NewValueSourceChain(yaml.YAML(debugFlagName, altsrc.StringSourcer(configFilePath))),
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
				Debug:           cmd.Bool(debugFlagName),
				ConfigFilePath:  configFilePath,
			}
			return params.Execute()
		},
	}

	if err := cmd.Run(context.Background(), os.Args); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
