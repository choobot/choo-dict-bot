package controller

import (
	"errors"
	"strings"
	"sync"
	"time"

	"github.com/choobot/choo-dict-bot/app/service"
)

type ServiceController interface {
	FindDefinitionsAndSynonyms(userID string, word string) (string, string, error)
}

type DictServiceController struct {
	dictService     service.DictService
	userProgress    map[string]int
	concurrent      int
	maxPerMinute    int
	concurrentMux   sync.Mutex
	userProgressMux sync.Mutex
}

func (this *DictServiceController) FindDefinitionsAndSynonyms(userID string, word string) (string, string, error) {
	if this.concurrent >= this.maxPerMinute {
		return "", "", errors.New("Sorry, we've reached the number of requests limit, please wait for 1 minute and try again.")
	}
	this.userProgressMux.Lock()
	progress := this.userProgress[userID]
	this.userProgressMux.Unlock()
	if progress > 0 {
		return "", "", errors.New("You're too fast, please slow down.")
	} else {
		this.concurrentMux.Lock()
		this.concurrent++
		this.concurrentMux.Unlock()
		this.userProgressMux.Lock()
		this.userProgress[userID] = 2
		this.userProgressMux.Unlock()
		word = strings.Split(word, " ")[0]
		definistionsCh := make(chan string)
		synonymsCh := make(chan string)
		errorCh := make(chan error)
		go func() {
			time.Sleep(time.Duration(60000/this.maxPerMinute) * time.Millisecond)
			res, err := this.dictService.FindDefinitions(word)
			if err != nil {
				errorCh <- err
			} else {
				definistionsCh <- res
			}

		}()
		go func() {
			time.Sleep(time.Duration(60000/this.maxPerMinute) * time.Millisecond)
			res, err := this.dictService.FindSynonyms(word)
			if err != nil {
				errorCh <- err
			} else {
				synonymsCh <- res
			}

		}()
		definistions, synonyms := "", ""
	Loop:
		for {
			select {
			case definistions = <-definistionsCh:
				this.userProgressMux.Lock()
				this.userProgress[userID]--
				userProgress := this.userProgress[userID]
				this.userProgressMux.Unlock()
				if userProgress == 0 {
					break Loop
				}
			case synonyms = <-synonymsCh:
				this.userProgressMux.Lock()
				this.userProgress[userID]--
				userProgress := this.userProgress[userID]
				this.userProgressMux.Unlock()
				if userProgress == 0 {
					break Loop
				}

			case error := <-errorCh:
				this.userProgressMux.Lock()
				this.userProgress[userID] = 0
				this.userProgressMux.Unlock()
				return "", "", errors.New("There was error on DictService: " + error.Error())
			}
		}
		return definistions, synonyms, nil
	}
}

func NewServiceController(dictService service.DictService, maxPerMinute int) *DictServiceController {
	serviceController := DictServiceController{
		userProgress: map[string]int{},
		concurrent:   0,
		maxPerMinute: maxPerMinute,
		dictService:  dictService,
	}
	limiter := time.Tick(time.Minute)
	go func() {
		for {
			<-limiter
			serviceController.concurrentMux.Lock()
			serviceController.concurrent = 0
			serviceController.concurrentMux.Unlock()
		}
	}()
	return &serviceController
}
