package nordic

type generator struct{}

func (g generator) Name() string { return "Nordic" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
	return "NordicFirst", "NordicLast"
}

var GeneratorInstance generator
