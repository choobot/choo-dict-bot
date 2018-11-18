package bot

import (
	"errors"
	"os"
	"testing"

	"github.com/line/line-bot-sdk-go/linebot"
)

type mockServiceController struct {
}

func (this mockServiceController) FindDefinitionsAndSynonyms(userID string, word string) (string, string, error) {
	if word == "error_word" {
		return "", "", errors.New("dummy")
	}
	return "dummy", "dummy", nil
}

func TestDictBotResponse(t *testing.T) {
	wantErr := errors.New("linebot: APIError 400 Invalid reply token")
	client, _ := linebot.New(os.Getenv("LINE_BOT_SECRET"), os.Getenv("LINE_BOT_TOKEN"))
	serviceController := mockServiceController{}
	bot := &DictBot{
		ServiceController: serviceController,
		Client:            client,
	}

	events := []*linebot.Event{}
	err := bot.Response(events)
	if err != nil {
		t.Errorf("DictBot.Response(%v) == %v, want %v", events, err, nil)
	}

	event := linebot.Event{
		Type: linebot.EventTypeMessage,
		Message: &linebot.TextMessage{
			Text: "dummy",
		},
		Source: &linebot.EventSource{
			UserID: "dummy",
		},
		ReplyToken: "dummy",
	}
	events = append(events, &event)
	bot.Response(events)
	err = bot.Response(events)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("DictBot.Response(%v) == %v, want %v", events, err, wantErr)
	}

	event = linebot.Event{
		Type: linebot.EventTypeMessage,
		Message: &linebot.TextMessage{
			Text: "error_word",
		},
		Source: &linebot.EventSource{
			UserID: "dummy",
		},
		ReplyToken: "dummy",
	}
	events = append(events, &event)
	err = bot.Response(events)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("DictBot.Response(%v) == %v, want %v", events, err, wantErr)
	}

	event = linebot.Event{
		Type: linebot.EventTypeJoin,
		Message: &linebot.TextMessage{
			Text: "error_word",
		},
		Source: &linebot.EventSource{
			UserID: "dummy",
		},
		ReplyToken: "dummy",
	}
	events = append(events, &event)
	err = bot.Response(events)
	if err == nil || err.Error() != wantErr.Error() {
		t.Errorf("DictBot.Response(%v) == %v, want %v", events, err, wantErr)
	}
}
