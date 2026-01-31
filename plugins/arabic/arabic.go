package arabic

type generator struct{}

func (g generator) Name() string { return "Arabic" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
	return "ArabicFirst", "ArabicLast"
}

var GeneratorInstance generator
