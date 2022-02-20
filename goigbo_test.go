package main

import (
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
type mockClient struct {
	c TestCase_New
}

func (m *mockClient) Do(req *http.Request) (*http.Response, error) {
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
				Igbo:            "Igbo_1",
				English:         "English_1",
				AssociatedWords: []string{"AssociatedWords_1", "AssociatedWords_1"},
				Pronunciation:   "Pronunciation_1",
				UpdatedOn:       "UpdatedOn_1",
				Id:              "Id_1",
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
			client:      &mockClient{},
			expected:    nil,
		},
		{
			description: "new should create a goigbdo instance when passed an apikey",
			apikey:      "",
			keyword:     "health",
			client:      &mockClient{},
			expected:    &ErrApiKeyRequired{},
		},
	}
	for _, c := range cases {
		t.Run(c.description, func(t *testing.T) {
			_, err := New(c.apikey, c.client)
			if err != c.expected {
				t.Fatalf("Expected returned error to be '%v' and received '%v'", c.expected, err)
			}

		})
	}
}
