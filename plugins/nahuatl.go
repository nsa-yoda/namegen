
package main

type generator struct{}

func (g generator) Name() string { return "Nahuatl" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
    return "NahuatlFirst", "NahuatlLast"
}

var GeneratorInstance generator
