package main

import (
	"errors"
	"os"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestOxfordServiceMapToString(t *testing.T) {
	cases := []struct {
		in   map[string]int
		want string
	}{
		{
			map[string]int{
				"1": 0,
				"2": 0,
				"3": 0,
				"4": 0,
			},
			"1, 2, 3 and 4",
		},
		{
			map[string]int{
				"1": 0,
				"2": 0,
				"3": 0,
				"4": 0,
				"5": 0,
			},
			"1, 2, 3, 4 and 5",
		},
		{
			map[string]int{},
			"",
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:          os.Getenv("OXFORD_API_ID"),
			appKey:         os.Getenv("OXFORD_API_KEY"),
			endpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got := service.MapToString(c.in)
		if got != c.want {
			t.Errorf("OxfordService.MapToString(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestOxfordServiceUnmarshallSynonyms(t *testing.T) {
	cases := []struct {
		in   []byte
		want string
	}{
		{
			[]byte(""),
			"",
		},
		{
			[]byte(`{    "results": [        {            "lexicalEntries": [                {                    "entries": [                        {                            "senses": [                                {                                    "synonyms": [                                        {                                            "text": "1"                                        },                                        {                                            "text": "2"                                        },                                        {                                            "text": "3"                                        }                                    ]                                }                            ]                        }                    ]                }            ]        }    ]}`),
			"1, 2 and 3",
		},
		{
			[]byte(`{    "results": [        {            "lexicalEntries": [                {                    "entries": [                        {                            "senses": [                                {                                    "subsenses": [                                        {                                            "synonyms": [                                                {                                                    "text": "1"                                                },                                                {                                                    "text": "2"                                                },                                                {                                                    "text": "3"                                                }                                            ]                                        }                                    ]                                }                            ]                        }                    ]                }            ]        }    ]}`),
			"1, 2 and 3",
		},
		{
			[]byte(`{    "results": [        {            "lexicalEntries": [                {                    "entries": [                        {                            "senses": [                                {                                    "synonyms": [                                        {                                            "text": "1"                                        },                                        {                                            "text": "2"                                        },                                        {                                            "text": "3"                                        }                                    ],                                    "subsenses": [                                        {                                            "synonyms": [                                                {                                                    "text": "4"                                                },                                                {                                                    "text": "5"                                                },                                                {                                                    "text": "6"                                                }                                            ]                                        }                                    ]                                }                            ]                        }                    ]                }            ]        }    ]}`),
			"1, 2, 3, 4 and 5",
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:          os.Getenv("OXFORD_API_ID"),
			appKey:         os.Getenv("OXFORD_API_KEY"),
			endpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got := service.UnmarshallSynonyms(c.in)
		if got != c.want {
			t.Errorf("OxfordService.UnmarshallSynonyms(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestOxfordServiceUnmarshallDefinitions(t *testing.T) {
	cases := []struct {
		in   []byte
		want string
	}{
		{
			[]byte(""),
			"",
		},
		{
			[]byte(`{    "results": [        {            "lexicalEntries": [                {                    "entries": [                        {                            "senses": [                                {                                    "definitions": [                                        "1",                                        "2",                                        "3"                                    ]                                }                            ]                        }                    ]                }            ]        }    ]}`),
			"1",
		},
		{
			[]byte(`{    "results": [        {            "lexicalEntries": [                {                    "entries": [                        {                            "senses": [                                {                                    "subsenses": [                                        {                                            "definitions": [                                                "1",                                                "2",                                                "3"                                            ]                                        }                                    ]                                }                            ]                        }                    ]                }            ]        }    ]}`),
			"1",
		},
		{
			[]byte(`{    "results": [        {            "lexicalEntries": [                {                    "entries": [                        {                            "senses": [                                {                                    "definitions": [                                        "1",                                        "2",                                        "3"                                    ],                                    "subsenses": [                                        {                                            "definitions": [                                                "1",                                                "2",                                                "3"                                            ]                                        }                                    ]                                }                            ]                        }                    ]                }            ]        }    ]}`),
			"1",
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:          os.Getenv("OXFORD_API_ID"),
			appKey:         os.Getenv("OXFORD_API_KEY"),
			endpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got := service.UnmarshallDefinitions(c.in)
		if got != c.want {
			t.Errorf("OxfordService.UnmarshallDefinitions(%q) == %q, want %q", c.in, got, c.want)
		}
	}
}

func TestOxfordServiceFindDefinitions(t *testing.T) {
	cases := []struct {
		in   string
		want string
		err  error
	}{
		{
			"line",
			"a long, narrow mark or band",
			nil,
		},
		{
			"square",
			"a plane figure with four equal straight sides and four right angles",
			nil,
		},
		{
			"choopong",
			"No definition for 'choopong'.",
			nil,
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:          os.Getenv("OXFORD_API_ID"),
			appKey:         os.Getenv("OXFORD_API_KEY"),
			endpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got, err := service.FindDefinitions(c.in)
		if got != c.want || (c.err == nil && err != nil) || (c.err != nil && c.err.Error() != err.Error()) {
			t.Errorf("OxfordService.FindDefinitions(%q) == %q, %q , want %q, %q", c.in, got, err, c.want, c.err)
		}
	}

	//Test if invalid API ID, Key
	service := &OxfordService{
		appId:          "dummy",
		appKey:         "dummy",
		endpointPrefix: "https://od-api.oxforddictionaries.com",
	}

	word := "dummy"
	wantErr := "Authentication failed"
	_, err := service.FindDefinitions(word)
	if err.Error() != "Authentication failed" {
		t.Errorf("OxfordService.FindDefinitions(%q) == %q , want %q", word, err, wantErr)
	}

	//Test if error by HTTP client
	service = &OxfordService{
		appId:          "dummy",
		appKey:         "dummy",
		endpointPrefix: "https://dummydomain",
	}
	_, err = service.FindDefinitions(word)
	wantErr = "no such host"
	if !strings.Contains(err.Error(), wantErr) {
		t.Errorf("OxfordService.FindDefinitions(%q) == %q , want %q", word, err, wantErr)
	}
}

func TestOxfordServiceFindSynonyms(t *testing.T) {
	cases := []struct {
		in   string
		want string
		err  error
	}{
		{
			"line",
			"bar, dash, rule, score and underline",
			nil,
		},
		{
			"square",
			"close, courtyard, marketplace, quad and quadrangle",
			nil,
		},
		{
			"choopong",
			"No synonyms for 'choopong'.",
			nil,
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:          os.Getenv("OXFORD_API_ID"),
			appKey:         os.Getenv("OXFORD_API_KEY"),
			endpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got, err := service.FindSynonyms(c.in)
		if got != c.want || (c.err == nil && err != nil) || (c.err != nil && c.err.Error() != err.Error()) {
			t.Errorf("OxfordService.FindSynonyms(%q) == %q, %q , want %q, %q", c.in, got, err, c.want, c.err)
		}
	}
	//Test if invalid API ID, Key
	service := &OxfordService{
		appId:          "dummy",
		appKey:         "dummy",
		endpointPrefix: "https://od-api.oxforddictionaries.com",
	}

	word := "dummy"
	wantErr := "Authentication failed"
	_, err := service.FindSynonyms(word)
	if err.Error() != "Authentication failed" {
		t.Errorf("OxfordService.FindSynonyms(%q) == %q , want %q", word, err, wantErr)
	}

	//Test if error by HTTP client
	service = &OxfordService{
		appId:          "dummy",
		appKey:         "dummy",
		endpointPrefix: "https://dummydomain",
	}
	_, err = service.FindSynonyms(word)
	wantErr = "no such host"
	if !strings.Contains(err.Error(), wantErr) {
		t.Errorf("OxfordService.FindSynonyms(%q) == %q , want %q", word, err, wantErr)
	}
}

type MockDickService struct {
}

func (this *MockDickService) FindDefinitions(word string) (string, error) {
	if word == "delay_word" {
		time.Sleep(10 * time.Millisecond)
	} else if word == "error_word" {
		return "", errors.New("DummyError")
	}
	return "a long, narrow mark or band", nil
}
func (this *MockDickService) FindSynonyms(word string) (string, error) {
	if word == "delay_word" {
		time.Sleep(10 * time.Millisecond)
	} else if word == "error_word" {
		return "", errors.New("DummyError")
	}
	return "bar, dash, rule, score and underline", nil
}

func TestServiceControllerFindDefinitionsAndSynonyms(t *testing.T) {
	dictService := &MockDickService{}
	serviceController := NewServiceController(dictService, 60)
	word := "delay_word"
	wantDefinistions := "a long, narrow mark or band"
	wantSynonyms := "bar, dash, rule, score and underline"
	go func() {
		definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms("dummy_user", word)
		if definistions != wantDefinistions || synonyms != wantSynonyms {
			t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q", "dummy_user", word, definistions, synonyms, err, wantDefinistions, wantSynonyms)
		}
	}()

	// Multiple user at the time
	go func() {
		definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms("dummy_user2", word)
		if definistions != wantDefinistions || synonyms != wantSynonyms {
			t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q", "dummy_user", word, definistions, synonyms, err, wantDefinistions, wantSynonyms)
		}
	}()

	// Same user at the time
	time.Sleep(5 * time.Millisecond)
	wantErr := "You're too fast, please slow down."
	definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms("dummy_user", word)
	if err == nil || err.Error() != wantErr {
		t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q", "dummy_user", word, definistions, synonyms, err, wantErr)
	}

	//Same user after previous result
	time.Sleep(2 * time.Second)
	definistions, synonyms, err = serviceController.FindDefinitionsAndSynonyms("dummy_user", word)
	if definistions != wantDefinistions || synonyms != wantSynonyms {
		t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q", "dummy_user", word, definistions, synonyms, err, wantDefinistions, wantSynonyms)
	}

	// Some error on Dict API
	word = "error_word"
	wantErr = "There was error on DictService: DummyError"
	definistions, synonyms, err = serviceController.FindDefinitionsAndSynonyms("dummy_user3", word)
	if err == nil || err.Error() != wantErr {
		t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q", "dummy_user", word, definistions, synonyms, err, wantErr)
	}
}

func TestConcurrentLoadServiceControllerFindDefinitionsAndSynonyms(t *testing.T) {
	concurrent := 10000
	dictService := &MockDickService{}
	serviceController := NewServiceController(dictService, concurrent)
	word := "line"
	resultCh := make(chan bool)
	for i := 0; i < concurrent; i++ {
		go func(i int) {
			userID := "user" + strconv.Itoa(i)
			wantDefinistions := "a long, narrow mark or band"
			wantSynonyms := "bar, dash, rule, score and underline"
			definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms(userID, word)
			if definistions != wantDefinistions || synonyms != wantSynonyms || err != nil {
				t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q", userID, word, definistions, synonyms, err, wantDefinistions, wantSynonyms)
			}
			resultCh <- true
		}(i)
	}

	// Wait for all results
	for i := 0; i < concurrent; i++ {
		<-resultCh
	}

	// More than limit
	userID := "dummy_user1"
	wantDefinistions := ""
	wantSynonyms := ""
	wantErr := "Sorry, we've reached the number of requests limit, please wait for 1 minute and try again."
	definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms(userID, word)
	if definistions != wantDefinistions || synonyms != wantSynonyms || err == nil || err.Error() != wantErr {
		t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q %q", userID, word, definistions, synonyms, err, wantDefinistions, wantSynonyms, wantErr)
	}

	//Wait for more than one minute, the limit should be reset
	time.Sleep(60 * time.Second)
	userID = "dummy_user2"
	word = "line"
	wantDefinistions = "a long, narrow mark or band"
	wantSynonyms = "bar, dash, rule, score and underline"
	definistions, synonyms, err = serviceController.FindDefinitionsAndSynonyms(userID, word)
	if definistions != wantDefinistions || synonyms != wantSynonyms || err != nil {
		t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q", userID, word, definistions, synonyms, err, wantDefinistions, wantSynonyms)
	}
}
