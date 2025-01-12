package test_report

type BrunoReport struct {
	Results []BrunoResult `json:"results"`
}

type BrunoResult struct {
	FileOrigin string
	Test       BrunoTest     `json:"test"`
	Request    BrunoRequest  `json:"request"`
	Response   BrunoResponse `json:"response"`
}

type BrunoTest struct {
	Filename string `json:"filename"`
}

type BrunoRequest struct {
	Method  string         `json:"method"`
	Url     string         `json:"url"`
	Headers map[string]any `json:"headers"`
	Body    CustomString   `json:"data"`
}

type BrunoResponse struct {
	Status  int            `json:"status"`
	Headers map[string]any `json:"headers"`
	Body    CustomString   `json:"data"`
}

// CustomString unmarshals all the value of the key as string
// In the case of Bruno, bodies are stored as json
type CustomString string

func (s *CustomString) UnmarshalJSON(data []byte) error {
	*s = CustomString(data)
	return nil
}
