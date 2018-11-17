package service

import (
	"errors"
	"io/ioutil"
	"net/http"
	"sort"

	"github.com/buger/jsonparser"
)

type DictService interface {
	FindDefinitions(word string) (string, error)
	FindSynonyms(word string) (string, error)
}

type OxfordService struct {
	AppId          string
	AppKey         string
	EndpointPrefix string
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
	req, err := http.NewRequest("GET", this.EndpointPrefix+"/api/v1/entries/en/"+word, nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", this.AppId)
	req.Header.Add("app_key", this.AppKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return "No definition for '" + word + "'.", nil
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
	req, err := http.NewRequest("GET", this.EndpointPrefix+"/api/v1/entries/en/"+word+"/synonyms", nil)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("app_id", this.AppId)
	req.Header.Add("app_key", this.AppKey)
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()
	if res.StatusCode == http.StatusNotFound {
		return "No synonyms for '" + word + "'.", nil
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
