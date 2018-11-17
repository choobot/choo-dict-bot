package bot

import (
	"github.com/choobot/choo-dict-bot/app/controller"
	"github.com/line/line-bot-sdk-go/linebot"
)

type DictBot struct {
	ServiceController controller.ServiceController
	Client            *linebot.Client
}

func (this *DictBot) Response(events []*linebot.Event) error {
	for _, event := range events {
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				word := message.Text
				definistions, synonyms, err := this.ServiceController.FindDefinitionsAndSynonyms(event.Source.UserID, word)
				if err != nil {
					if _, err = this.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(err.Error())).Do(); err != nil {
						return err
					}
				} else {
					if _, err = this.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(definistions), linebot.NewTextMessage(synonyms)).Do(); err != nil {
						return err
					}
				}
			}
		} else if event.Type == linebot.EventTypeJoin {
			replyMessage := "Thanks for adding me. I'm Choo Dict Bot, I'm here to help you to find English word definitions and synonyms. Try to send me some words."
			if _, err := this.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do(); err != nil {
				return err
			}
		}
	}
	return nil
}
