# go-igbo_api
Golang wrapper for utilizing https://github.com/ijemmao/igbo_api


# GetWords

Get Words accepts a string value and queries the go-igbo_api for words and pronounciations related to the input.

```golang

package main

import (
    goigbo "github.com/DCalliet/go-igbo_api"
)


func main() {
    // Accepts input
    value := "health"
    goigboInstance := goigbo.GetNew(os.GetEnv("IGBO_API_KEY"))

    // GetWords
    words, err := goigboInstance.GetWords(value)
    fmt.Print(words)
}


```


# GetExamples

Get Examples accets a string values and queries the go-igbo_api for phrases in igbo and in english.

```golang

package main

import (
    goigbo "github.com/DCalliet/go-igbo_api"
)


func main() {
    // Accepts input
    value := "health"
    goigboInstance := goigbo.GetNew(os.GetEnv("IGBO_API_KEY"))

    // GetWords
    examples, err := goigboInstance.GetWords(value)
    fmt.Print(examples)
}


```