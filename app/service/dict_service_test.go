package service

import (
	"os"
	"strings"
	"testing"
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
			AppId:          os.Getenv("OXFORD_API_ID"),
			AppKey:         os.Getenv("OXFORD_API_KEY"),
			EndpointPrefix: "https://od-api.oxforddictionaries.com",
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
			AppId:          os.Getenv("OXFORD_API_ID"),
			AppKey:         os.Getenv("OXFORD_API_KEY"),
			EndpointPrefix: "https://od-api.oxforddictionaries.com",
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
			AppId:          os.Getenv("OXFORD_API_ID"),
			AppKey:         os.Getenv("OXFORD_API_KEY"),
			EndpointPrefix: "https://od-api.oxforddictionaries.com",
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
			AppId:          os.Getenv("OXFORD_API_ID"),
			AppKey:         os.Getenv("OXFORD_API_KEY"),
			EndpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got, err := service.FindDefinitions(c.in)
		if got != c.want || (c.err == nil && err != nil) || (c.err != nil && c.err.Error() != err.Error()) {
			t.Errorf("OxfordService.FindDefinitions(%q) == %q, %q , want %q, %q", c.in, got, err, c.want, c.err)
		}
	}

	//Test if invalid API ID, Key
	service := &OxfordService{
		AppId:          "dummy",
		AppKey:         "dummy",
		EndpointPrefix: "https://od-api.oxforddictionaries.com",
	}

	word := "dummy"
	wantErr := "Authentication failed"
	_, err := service.FindDefinitions(word)
	if err.Error() != "Authentication failed" {
		t.Errorf("OxfordService.FindDefinitions(%q) == %q , want %q", word, err, wantErr)
	}

	//Test if error by HTTP client
	service = &OxfordService{
		AppId:          "dummy",
		AppKey:         "dummy",
		EndpointPrefix: "https://dummydomain",
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
			AppId:          os.Getenv("OXFORD_API_ID"),
			AppKey:         os.Getenv("OXFORD_API_KEY"),
			EndpointPrefix: "https://od-api.oxforddictionaries.com",
		}
		got, err := service.FindSynonyms(c.in)
		if got != c.want || (c.err == nil && err != nil) || (c.err != nil && c.err.Error() != err.Error()) {
			t.Errorf("OxfordService.FindSynonyms(%q) == %q, %q , want %q, %q", c.in, got, err, c.want, c.err)
		}
	}
	//Test if invalid API ID, Key
	service := &OxfordService{
		AppId:          "dummy",
		AppKey:         "dummy",
		EndpointPrefix: "https://od-api.oxforddictionaries.com",
	}

	word := "dummy"
	wantErr := "Authentication failed"
	_, err := service.FindSynonyms(word)
	if err.Error() != "Authentication failed" {
		t.Errorf("OxfordService.FindSynonyms(%q) == %q , want %q", word, err, wantErr)
	}

	//Test if error by HTTP client
	service = &OxfordService{
		AppId:          "dummy",
		AppKey:         "dummy",
		EndpointPrefix: "https://dummydomain",
	}
	_, err = service.FindSynonyms(word)
	wantErr = "no such host"
	if !strings.Contains(err.Error(), wantErr) {
		t.Errorf("OxfordService.FindSynonyms(%q) == %q , want %q", word, err, wantErr)
	}
}
