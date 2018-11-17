package controller

import (
	"errors"
	"strconv"
	"testing"
	"time"
)

type mockDictService struct {
}

func (this *mockDictService) FindDefinitions(word string) (string, error) {
	if word == "delay_word" {
		time.Sleep(10 * time.Millisecond)
	} else if word == "error_word" {
		return "", errors.New("DummyError")
	}
	return "a long, narrow mark or band", nil
}
func (this *mockDictService) FindSynonyms(word string) (string, error) {
	if word == "delay_word" {
		time.Sleep(10 * time.Millisecond)
	} else if word == "error_word" {
		return "", errors.New("DummyError")
	}
	return "bar, dash, rule, score and underline", nil
}

func TestServiceControllerFindDefinitionsAndSynonyms(t *testing.T) {
	dictService := &mockDictService{}
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
	dictService := &mockDictService{}
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
