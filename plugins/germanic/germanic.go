package germanic

type generator struct{}

func (g generator) Name() string { return "Germanic" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
	return "GermanicFirst", "GermanicLast"
}

var GeneratorInstance generator
