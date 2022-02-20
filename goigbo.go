package goigbo

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// New will accept an api key and return an interface that will have a
// GetWords function & a GetExample function
func New(apikey string, client http_Do) (*GoIgboClient, error) {
	if apikey != "" {
		instance := GoIgboClient{
			apikey: apikey,
			client: client,
		}
		return &instance, nil
	}
	return nil, &ErrApiKeyRequired{}
}

type GoIgboClient struct {
	apikey string
	client http_Do
}

type http_Do interface {
	Do(req *http.Request) (*http.Response, error)
}

// GetWords will retrieve a keyword and return an array of revelant GetWordsOutput
func (g *GoIgboClient) GetWords(keyword string) ([]GetWordsOutput, error) {
	// Create an http request
	request, err := http.NewRequest("GET", "https://www.igboapi.com/api/v1/words", nil)
	if err != nil {
		return []GetWordsOutput{}, err
	}
	// Set the Request Header
	request.Header.Add("X-API-Key", g.apikey)

	// Apply keyword to url.Values
	q := request.URL.Query()
	q.Add("keyword", keyword)
	request.URL.RawQuery = q.Encode()

	// Execute Request
	response, err := g.client.Do(request)
	if err != nil {
		return []GetWordsOutput{}, err
	}
	// http module recommends closing the body after a request
	defer response.Body.Close()

	var n int = -1
	var outputBytes []byte
	var output []GetWordsOutput
	for n != 0 {
		b := make([]byte, 1024)
		n, err = response.Body.Read(b)
		if err != io.EOF && err != nil {
			return []GetWordsOutput{}, err
		}
		outputBytes = append(outputBytes, b...)
	}

	// migrate our byte array into a structure we can return
	err = json.Unmarshal(bytes.Trim(outputBytes, "\x00"), &output)
	if err != nil {
		return []GetWordsOutput{}, &ErrJsonUnrecognized{
			n:     n,
			bytes: outputBytes,
			err:   err,
		}
	}
	return output, err
}

type GetWordsOutput struct {
	WordClass      string   `json:"wordClass"`
	Definitions    []string `json:"definitions"`
	Variations     []string `json:"variations"`
	Stems          []string `json:"stems"`
	Word           string   `json:"word"`
	IsStandardIgbo *bool    `json:"isStandardIgbo"`
	Antonyms       []string `json:"antonyms"`
	Hypernyms      []string `json:"hypernyms"`
	Hyponyms       []string `json:"hyponyms"`
	Synonyms       []string `json:"synonyms"`
	Nsibidi        string   `json:"nsibidi"`
}

type GetExampleOutput struct {
	Igbo            string   `json:"igbo"`
	English         string   `json:"english"`
	AssociatedWords []string `json:"associatedWords"`
	Pronunciation   string   `json:"pronunciation"`
	UpdatedOn       string   `json:"-"`
	Id              string   `json:"id"`
}

type ErrJsonUnrecognized struct {
	n     int
	bytes []byte
	err   error
}

func (e *ErrJsonUnrecognized) Error() string {
	return fmt.Sprintf("failed to recognized %d bytes of json: %v (%s)", e.n, e.err, string(e.bytes))
}

type ErrApiKeyRequired struct{}

func (e *ErrApiKeyRequired) Error() string {
	return "api key is required to create a new instance of goigbo"
}
