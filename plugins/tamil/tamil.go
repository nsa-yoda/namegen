package tamil

type generator struct{}

func (g generator) Name() string { return "Tamil" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
	return "TamilFirst", "TamilLast"
}

var GeneratorInstance generator
