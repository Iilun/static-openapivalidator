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
	config          validator.Config
}

type Config struct {
	Banned Banned `yaml:"banned"`
	Ignore Ignore `yaml:"ignore"`
}

type Banned struct {
	Requests  []string `yaml:"requests"`
	Responses []string `yaml:"responses"`
	Routes    []string `yaml:"routes"`
}

type Ignore struct {
	Servers bool `yaml:"servers"`
}
