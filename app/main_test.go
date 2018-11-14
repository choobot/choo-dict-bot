package main

import (
	"errors"
	"testing"
)

var appId = "cf094011"
var appKey = "f551e904c5b65e71766ff179cf486d07"

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
			appId:  appId,
			appKey: appKey,
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
			appId:  appId,
			appKey: appKey,
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
			appId:  appId,
			appKey: appKey,
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
			"line line2",
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
			"",
			errors.New("No definition for 'choopong'."),
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:  appId,
			appKey: appKey,
		}
		got, err := service.FindDefinitions(c.in)
		if got != c.want || (c.err == nil && err != nil) || (c.err != nil && c.err.Error() != err.Error()) {
			t.Errorf("OxfordService.FindDefinitions(%q) == %q, %q , want %q, %q", c.in, got, err, c.want, c.err)
		}
	}

	service := &OxfordService{
		appId:  "dummy",
		appKey: "dummy",
	}
	_, err := service.FindDefinitions("dummy")
	if err.Error() != "Authentication failed" {
		t.Errorf("OxfordService.FindDefinitions(%q) == %q , want %q", "dummy", err, "")
	}
}

func TestOxfordServiceFindSynonyms(t *testing.T) {
	cases := []struct {
		in   string
		want string
		err  error
	}{
		{
			"line line2",
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
			"",
			errors.New("No synonyms for 'choopong'."),
		},
	}
	for _, c := range cases {
		service := &OxfordService{
			appId:  appId,
			appKey: appKey,
		}
		got, err := service.FindSynonyms(c.in)
		if got != c.want || (c.err == nil && err != nil) || (c.err != nil && c.err.Error() != err.Error()) {
			t.Errorf("OxfordService.FindSynonyms(%q) == %q, %q , want %q, %q", c.in, got, err, c.want, c.err)
		}
	}
	service := &OxfordService{
		appId:  "dummy",
		appKey: "dummy",
	}
	_, err := service.FindSynonyms("dummy")
	if err.Error() != "Authentication failed" {
		t.Errorf("OxfordService.FindSynonyms(%q) == %q , want %q", "dummy", err, "")
	}
}

func TestServiceControllerFindDefinitionsAndSynonyms(t *testing.T) {
	serviceController := &ServiceController{
		dictService: &OxfordService{},
	}
	wantDefinistions := "a long, narrow mark or band"
	wantSynonyms := "bar, dash, rule, score and underline"
	go func() {
		definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms("dummy_user", "line")
		if definistions != wantDefinistions || synonyms != wantSynonyms {
			t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q, %q", "dummy_user", "line", definistions, synonyms, err, wantDefinistions, wantSynonyms)
		}
	}()

	// wantErr := "You're too fast, please wait for the result of previous inquiry"
	// definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms("dummy_user", "line")
	// if err.Error() != wantErr {
	// 	t.Errorf("ServiceController.FindDefinitionsAndSynonyms(%q, %q) == %q, %q %q, want %q", "dummy_user", "line", definistions, synonyms, err, wantErr)
	// }

}
