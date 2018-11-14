package main

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sort"
	"strings"
	"sync"

	"github.com/buger/jsonparser"
)

type DictService interface {
	FindDefinitions(word string) (string, error)
	FindSynonyms(word string) (string, error)
}

type OxfordService struct {
	appId  string
	appKey string
}

type ServiceController struct {
	dictService DictService
}

func main() {
	fmt.Print("Please enter word: ")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		word := scanner.Text()
		serviceController := &ServiceController{
			dictService: &OxfordService{
				appId:  os.Getenv("OXFORD_API_ID"),
				appKey: os.Getenv("OXFORD_API_KEY"),
			},
		}
		definistions, synonyms, err := serviceController.FindDefinitionsAndSynonyms("dummy_user", word)
		if err != nil {
			fmt.Println(err.Error())
		} else {
			fmt.Println(definistions)
			fmt.Println(synonyms)
		}
		fmt.Print("Please enter word: ")
	}
}

func (this *OxfordService) UnmarshallDefinitions(data []byte) string {
	found := false
	values := ""
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						if found {
							return
						}
						found = true
						values = string(value)
					}, "definitions")
					jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
							if found {
								return
							}
							found = true
							values = string(value)
						}, "definitions")
					}, "subsenses")
				}, "senses")
			}, "entries")
		}, "lexicalEntries")
	}, "results")
	return values
}

func (this *OxfordService) FindDefinitions(word string) (string, error) {
	req, err := http.NewRequest("GET", "https://od-api.oxforddictionaries.com/api/v1/entries/en/"+word, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", this.appId)
	req.Header.Add("app_key", this.appKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return "", errors.New("No definition for '" + word + "'.")
	} else if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return this.UnmarshallDefinitions(body), nil
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		return "", errors.New(string(body))
	}
}

func (this *OxfordService) UnmarshallSynonyms(data []byte) string {
	count := 0
	values := map[string]int{}
	jsonparser.ArrayEach(data, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
			jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
				jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
					jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						if count >= 5 {
							return
						}
						val, err := jsonparser.GetString(value, "text")
						if err == nil {
							values[val] = 0
							count++
						}
					}, "synonyms")
					jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
						jsonparser.ArrayEach(value, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
							if count >= 5 {
								return
							}
							val, err := jsonparser.GetString(value, "text")
							if err == nil {
								values[val] = 0
								count++
							}
						}, "synonyms")
					}, "subsenses")
				}, "senses")
			}, "entries")
		}, "lexicalEntries")
	}, "results")
	return this.MapToString(values)
}

func (this *OxfordService) FindSynonyms(word string) (string, error) {
	req, err := http.NewRequest("GET", "https://od-api.oxforddictionaries.com/api/v1/entries/en/"+word+"/synonyms", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", this.appId)
	req.Header.Add("app_key", this.appKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println(err)
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return "", errors.New("No synonyms for '" + word + "'.")
	} else if res.StatusCode == http.StatusOK {
		body, _ := ioutil.ReadAll(res.Body)
		return this.UnmarshallSynonyms(body), nil
	} else {
		body, _ := ioutil.ReadAll(res.Body)
		return "", errors.New(string(body))
	}
}

func (this *OxfordService) MapToString(values map[string]int) string {
	text := ""
	i := 0
	keys := []string{}
	for k := range values {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if i == len(keys)-1 {
			text += " and "
		} else if i != 0 {
			text += ", "
		}
		text += k
		i++
	}

	return text
}

var userProgress = map[string]int{}

func (this *ServiceController) FindDefinitionsAndSynonyms(userID string, word string) (string, string, error) {
	var mutex = &sync.Mutex{}
	mutex.Lock()
	progress := userProgress[userID]
	if progress == 0 {
		userProgress[userID] = 2
		mutex.Unlock()
		word = strings.Split(word, " ")[0]
		fmt.Println("Looking for: " + word)
		definistionsCh := make(chan string)
		synonymsCh := make(chan string)
		go func() {
			res, err := this.dictService.FindDefinitions(word)
			if err != nil {
				definistionsCh <- err.Error()
			} else {
				definistionsCh <- res
			}

		}()
		go func() {
			res, err := this.dictService.FindSynonyms(word)
			if err != nil {
				synonymsCh <- err.Error()
			} else {
				synonymsCh <- res
			}

		}()
		definistions, synonyms := "", ""
	Loop:
		for {
			select {
			case definistions = <-definistionsCh:
				mutex.Lock()
				userProgress[userID]--
				mutex.Unlock()
			case synonyms = <-synonymsCh:
				mutex.Lock()
				userProgress[userID]--
				mutex.Unlock()
			default:
				mutex.Lock()
				if userProgress[userID] == 0 {
					break Loop
				}
				mutex.Unlock()
			}
		}
		return definistions, synonyms, nil
	} else {
		mutex.Unlock()
		return "", "", errors.New("You're too fast, please wait for the result of previous inquiry")
	}
}
