package internal

import (
	"context"
	"static-openapivalidator/validator"
)

type Params struct {
	Ctx             context.Context
	ApiFilePath     string   `validate:"required,file"`
	ReportFilePaths []string `validate:"gt=0,dive,file"`
	Format          string   `validate:"required"`
	JunitFilePath   string   `validate:"omitempty,filepath"`
	HtmlFilePath    string   `validate:"omitempty,filepath"`
	JsonFilePath    string   `validate:"omitempty,filepath"`
	ConfigFilePath  string   `validate:"omitempty,file"`
	Debug           bool
	config          validator.Config
}

type Config struct {
	Ignore Ignore `yaml:"ignore"`
}

type Ignore struct {
	Requests  []string `yaml:"requests"`
	Responses []string `yaml:"responses"`
	Routes    []string `yaml:"routes"`
	Servers   bool     `yaml:"servers"`
}
