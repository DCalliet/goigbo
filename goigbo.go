package goigbo

import (
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
func (g *GoIgboClient) GetWords(keyword string) ([]*GetWordsOutput, int, error) {
	// Bytes downloaded
	var n int
	// Create an http request
	request, err := http.NewRequest("GET", "https://www.igboapi.com/api/v1/words", nil)
	if err != nil {
		return []*GetWordsOutput{}, n, err
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
		return []*GetWordsOutput{}, n, err
	}
	// http module recommends closing the body after a request
	defer response.Body.Close()

	outputBytes := make([]byte, 512)
	var output []*GetWordsOutput

	n, err = response.Body.Read(outputBytes)
	if err != nil {
		return []*GetWordsOutput{}, n, err
	}
	// migrate our byte array into a structure we can return
	err = json.Unmarshal(outputBytes, &output)
	if err != nil {
		return []*GetWordsOutput{}, n, &ErrJsonUnrecognized{
			n:     n,
			bytes: outputBytes,
			err:   err,
		}
	}
	return output, n, err
}

/**
{
    "igbo": "Igwē nà-èji nji",
    "english": "The sky looks black",
    "associatedWords": [
      "5f90c35f49f7e863e92b825b"
    ],
    "pronunciation": "",
    "updatedOn": "2020-11-22T03:29:58.053Z",
    "id": "5f90c36949f7e863e92b916a"
  }
**/

type GetWordsOutput struct {
	Igbo            string   `json:"igbo"`
	English         string   `json:"english"`
	AssociatedWords []string `json:"associatedWords"`
	Pronunciation   string   `json:"pronunciation"`
	UpdatedOn       string   `json:"updatedOn"`
	Id              string   `json:"id"`
}

type GetWordsReader []GetWordsOutput

// Naive implementation of read will always read from beginning of json
// array and will always return io.EOF
func (g *GetWordsReader) Read(p []byte) (int, error) {
	bytes, err := json.Marshal(g)
	if err != nil {
		return 0, err
	}
	size := copy(p, bytes)
	return size, io.EOF
}

func (g *GetWordsReader) Close() error {
	return nil
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
