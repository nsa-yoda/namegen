
package main

type generator struct{}

func (g generator) Name() string { return "Chinese" }

func (g generator) GenerateName(gender string, realism int) (string, string) {
    return "ChineseFirst", "ChineseLast"
}

var GeneratorInstance generator
