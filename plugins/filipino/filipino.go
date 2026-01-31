package filipino

type generator struct{}

func (g generator) Name() string { return "Filipino" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
	return "FilipinoFirst", "FilipinoLast"
}

var GeneratorInstance generator
