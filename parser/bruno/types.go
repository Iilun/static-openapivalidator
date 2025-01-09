package bruno

type Report struct {
	Results []Result `json:"results"`
}

type Result struct {
	FileOrigin string
	Test       Test     `json:"test"`
	Request    Request  `json:"request"`
	Response   Response `json:"response"`
}

type Test struct {
	Filename string `json:"filename"`
}

type Request struct {
	Method  string            `json:"method"`
	Url     string            `json:"url"`
	Headers map[string]string `json:"headers"`
	Body    CustomString      `json:"data"`
}

type Response struct {
	Status  int               `json:"status"`
	Headers map[string]string `json:"headers"`
	Body    CustomString      `json:"data"`
}

type CustomString string

func (s *CustomString) UnmarshalJSON(data []byte) error {
	*s = CustomString(data)
	return nil
}
