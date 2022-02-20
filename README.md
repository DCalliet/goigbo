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
    // Create instance of http client
    client = &http.Client{}
    // Create instance of goigbo
    goigboInstance := goigbo.New(os.GetEnv("IGBO_API_KEY"), client)

    // Execute GetWords function
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
    // Create instance of http client
    client = &http.Client{}
    // Create instance of goigbo
    goigboInstance := goigbo.New(os.GetEnv("IGBO_API_KEY"), client)

    // Execute GetExamples function
    examples, err := goigboInstance.GetExamples(value)
    fmt.Print(examples)
}


```