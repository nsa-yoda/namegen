
package main

import (
    "flag"
    "fmt"
    "log"
    "plugin"
)

type Generator interface {
    Name() string
    GenerateName(gender string, realism int) (string, string)
}

func main() {
    lang := flag.String("lang", "english", "Language plugin to load")
    gender := flag.String("gender", "neutral", "male, female, neutral")
    realism := flag.Int("realism", 50, "0-100 realism blending")
    includeLast := flag.Bool("l", false, "include last name")
    reverse := flag.Bool("r", false, "reverse name order")
    flag.Parse()

    plug, err := plugin.Open(fmt.Sprintf("plugins/%s.so", *lang))
    if err != nil {
        log.Fatalf("Could not load plugin: %v", err)
    }

    sym, err := plug.Lookup("GeneratorInstance")
    if err != nil {
        log.Fatalf("Invalid plugin: %v", err)
    }

    gen := sym.(Generator)

    first, last := gen.GenerateName(*gender, *realism)

    if *includeLast {
        if *reverse {
            fmt.Println(last, first)
        } else {
            fmt.Println(first, last)
        }
    } else {
        fmt.Println(first)
    }
}
