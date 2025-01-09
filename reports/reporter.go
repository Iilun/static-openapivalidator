package reports

type Reporter interface {
	Generate(report Report) error
}
