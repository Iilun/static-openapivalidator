package json

import (
	_ "embed"
	"encoding/json"
	"os"
	"static-openapivalidator/reports"
)

func NewReporter(filePath string) Reporter {
	return Reporter{filePath: filePath}
}

type Reporter struct {
	filePath string
}

func (r Reporter) Generate(report reports.Report) error {
	bytes, err := json.Marshal(report)
	if err != nil {
		return err
	}
	return os.WriteFile(r.filePath, bytes, 0644)
}
