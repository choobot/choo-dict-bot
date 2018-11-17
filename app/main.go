package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/choobot/choo-dict-bot/app/bot"
	"github.com/choobot/choo-dict-bot/app/service"

	"github.com/choobot/choo-dict-bot/app/controller"
	"github.com/line/line-bot-sdk-go/linebot"
)

func main() {
	client, err := linebot.New(os.Getenv("LINE_BOT_SECRET"), os.Getenv("LINE_BOT_TOKEN"))
	if err != nil {
		log.Fatal(err)
	}
	dictService := &service.OxfordService{
		AppId:          os.Getenv("OXFORD_API_ID"),
		AppKey:         os.Getenv("OXFORD_API_KEY"),
		EndpointPrefix: "https://od-api.oxforddictionaries.com",
	}
	serviceController := controller.NewServiceController(dictService, 30)
	bot := &bot.DictBot{
		ServiceController: serviceController,
		Client:            client,
	}
	http.HandleFunc("/callback", func(w http.ResponseWriter, r *http.Request) {
		events, err := client.ParseRequest(r)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(http.StatusBadRequest)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
			return
		}
		if err := bot.Response(events); err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusInternalServerError)
		}

	})
	port := os.Getenv("PORT")
	if port == "" {
		port = "80"
	}
	fmt.Println("Runnign at :" + port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}
