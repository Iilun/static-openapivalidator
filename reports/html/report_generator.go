package html

import (
	_ "embed"
	"encoding/json"
	"os"
	"static-openapivalidator/logger"
	"static-openapivalidator/reports"
	"strings"
)

//go:embed report-template.html
var htmlTemplate string

func NewReporter(filePath string) Reporter {
	return Reporter{filePath: filePath}
}

type Reporter struct {
	filePath string
}

func (r Reporter) Generate(report reports.Report) error {
	logger.Log("Reporting: generating HTML report")
	bytes, err := json.Marshal(report)
	if err != nil {
		return err
	}

	fileContent := strings.Replace(htmlTemplate, "__RESULTS_JSON__", string(bytes), 1)
	return os.WriteFile(r.filePath, []byte(fileContent), 0644)
}
