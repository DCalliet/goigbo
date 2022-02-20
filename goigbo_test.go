package goigbo

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
type mockClient_New struct {
	c TestCase_New
}

func (m *mockClient_New) Do(req *http.Request) (*http.Response, error) {
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

	cases := []TestCase_GetWords{
		{
			description: "getWords should return GetWordsOutput when provided a keyword",
			keyword:     "health",
			err:         nil,
			expected: &GetWordsReader{
				{
					Igbo:            "Igbo_1",
					English:         "English_1",
					AssociatedWords: []string{"AssociatedWords_1", "AssociatedWords_1"},
					Pronunciation:   "Pronunciation_1",
					UpdatedOn:       "UpdatedOn_1",
					Id:              "Id_1",
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
				if output.Igbo != eoutput[i].Igbo {
					t.Fatalf("expected output.Igbo '%s' received '%s'", eoutput[i].Igbo, output.Igbo)
				}
				if output.English != eoutput[i].English {
					t.Fatalf("expected output.English '%s' received '%s'", eoutput[i].English, output.English)
				}
				if output.Id != eoutput[i].Id {
					t.Fatalf("expected output.Id '%s' received '%s'", eoutput[i].Id, output.Id)
				}
			}
		})
	}

}
