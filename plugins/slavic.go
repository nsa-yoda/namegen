
package main

type generator struct{}

func (g generator) Name() string { return "Slavic" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
    return "SlavicFirst", "SlavicLast"
}

var GeneratorInstance generator
