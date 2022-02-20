package goigbo

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"testing"
)

type TestCase_New struct {
	description string
	apikey      string
	keyword     string
	client      http_Do
	expected    error
}

// Mock a fake request->response interaction in our fake client
type mockClient_New struct {
	c TestCase_New
}

func (m *mockClient_New) Do(req *http.Request) (*http.Response, error) {
	var boolpointer bool = true
	if req.Header.Get("X-API-Key") != m.c.apikey {
		// Return error if no api key
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(`{"error":"X-API-Key Header doesn't exist"}`)),
		}, nil
	}
	// Return example payload
	return &http.Response{
		Body: &GetWordsReader{
			{
				WordClass:      "NNC",
				Definitions:    []string{"health"},
				Variations:     []string{},
				Stems:          []string{"ezi", "ndu"},
				Word:           "ezi ndu",
				IsStandardIgbo: &boolpointer,
			},
		},
	}, nil
}

// TestMain will control setup and teardown of needed testing environments
func TestMain(m *testing.M) {
	if os.Getenv("IGBO_API_KEY") == "" {
		log.Printf("expected env variable IGBO_API_KEY to be set")
		os.Exit(1)
	}
	code := m.Run()
	os.Exit(code)
}

// Test_New will accept api key and a client adhering to http.Do interface and return an interface that will have a GetWords function & a GetExample functions
func Test_New(t *testing.T) {
	cases := []TestCase_New{
		{
			description: "new should create a goigbdo instance when passed an apikey",
			apikey:      os.Getenv("IGBO_API_KEY"),
			keyword:     "health",
			client:      &mockClient_New{},
			expected:    nil,
		},
		{
			description: "new should fail to create a goigbdo instance when passed an apikey, and throw an error.",
			apikey:      "",
			keyword:     "health",
			client:      &mockClient_New{},
			expected:    &ErrApiKeyRequired{},
		},
	}
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			instance, err := New(c.apikey, c.client)
			if err != c.expected {
				t.Fatalf("Expected returned error to be '%v' and received '%v'", c.expected, err)
			}
			if c.expected == nil && instance.apikey != c.apikey {
				t.Fatalf("Expected apikey to be '%s' and receieved '%s'", c.apikey, instance.apikey)
			}
			if c.expected == nil && instance.client == nil {
				t.Fatal("Expected client to not be nil")
			}

		})
	}
}

type TestCase_GetWords struct {
	description string
	keyword     string
	err         error
	expected    *GetWordsReader
}

type GetWordsReader []GetWordsOutput

// Naive implementation of read will always read from beginning of json
// array and will always return io.EOF
func (g *GetWordsReader) Read(p []byte) (int, error) {
	bytes, err := json.Marshal(g)
	if err != nil {
		return 0, err
	}
	copy(p, bytes)
	return 0, io.EOF
}

func (g *GetWordsReader) Close() error {
	return nil
}

type mockClient_GetWords struct {
	t *testing.T
	c *TestCase_GetWords
}

func (m *mockClient_GetWords) Do(req *http.Request) (*http.Response, error) {
	if req.Header.Get("X-API-Key") != os.Getenv("IGBO_API_KEY") {
		// Return error if no api key
		return &http.Response{
			StatusCode: http.StatusBadRequest,
			Body:       io.NopCloser(strings.NewReader(`{"error":"X-API-Key Header doesn't exist"}`)),
		}, nil
	}
	for key, value := range req.URL.Query() {
		if key == "keyword" && value[0] != m.c.keyword {
			m.t.Errorf("expected keyword '%s' received '%s'", m.c.keyword, value)
		}
	}
	// Return example payload
	return &http.Response{
		Body: m.c.expected,
	}, nil
}

// Test_GetWords will accept a keyword and return an array of GetWordsOutput
func Test_GetWords(t *testing.T) {
	var boolpointer bool = true
	cases := []TestCase_GetWords{
		{
			description: "getWords should return GetWordsOutput when provided a keyword",
			keyword:     "health",
			err:         nil,
			// {NNC [health] [] [ezi ndù] ezi ndù 0xc0001aa6b0 [] [] [] [] }
			expected: &GetWordsReader{
				{
					WordClass:      "NNC",
					Definitions:    []string{"health"},
					Variations:     []string{},
					Stems:          []string{"ezi", "ndu"},
					Word:           "ezi ndu",
					IsStandardIgbo: &boolpointer,
				},
			},
		},
	}
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			// Set client in test so we can self reference the test case
			client := GoIgboClient{
				client: &mockClient_GetWords{
					t: t,
					c: &c,
				},
				apikey: os.Getenv("IGBO_API_KEY"),
			}
			result, err := client.GetWords(c.keyword)
			if err != c.err {
				t.Fatalf("expected error to be '%v' received '%v'", c.err, err)
			}
			for i, output := range result {
				eoutput := *c.expected
				if output.WordClass != eoutput[i].WordClass {
					t.Fatalf("expected output.WordClass '%s' received '%s'", eoutput[i].WordClass, output.WordClass)
				}
				if output.Word != eoutput[i].Word {
					t.Fatalf("expected output.English '%s' received '%s'", eoutput[i].Word, output.Word)
				}
				if output.IsStandardIgbo == nil {
					t.Fatal("expected output.IstandardIgbo to not be nil")
				}
			}
		})
	}

}
