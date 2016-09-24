package main

// taken from https://github.com/tucnak/telebot
import (

    "time"
    "os"

    "encoding/json"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/tucnak/telebot"
)

var bot *telebot.Bot
type Hito struct {
	File string `json:"file"`
	Title string `json:"title"`
}

type Data struct {
	Comment string `json:"comment"`
	hitos Hito
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Load milestones array
	file, e := ioutil.ReadFile("./hitos.json")
	if e != nil {
		log.Error("File error: %v\n", e)
		os.Exit(1)
	}
	var hitos Data
	if err := json.Unmarshal(file,&hitos); err != nil {
		log.Error(err)
	}
}

func main() {
    var err error
    bot, err = telebot.NewBot(os.Getenv("BOT_TOKEN"))
    if err != nil {
        log.Error(err)
    }

    bot.Messages = make(chan telebot.Message, 1000)
    bot.Queries = make(chan telebot.Query, 1000)

    go messages()
    go queries()

    bot.Start(1 * time.Second)
}

func messages() {
    for message := range bot.Messages {
	    log.WithFields(log.Fields {
		    "type": "message",
		    "username": message.Sender.Username,
		    "text": message.Text }).Info("Message received")
    }
}

func queries() {
    for query := range bot.Queries {
        log.WithFields(log.Fields {
		"type": "query",
		"from": query.From.Username,
		"text": query.Text }).Info("New query")

        // Create an article (a link) object to show in our results.
        article := &telebot.InlineQueryResultArticle{
            Title: "Telegram bot framework written in Go",
            URL:   "https://github.com/tucnak/telebot",
            InputMessageContent: &telebot.InputTextMessageContent{
                Text:           "Telebot is a convenient wrapper to Telegram Bots API, written in Golang.",
                DisablePreview: false,
            },
        }

        // Build the list of results. In this instance, just our 1 article from above.
        results := []telebot.InlineQueryResult{article}

        // Build a response object to answer the query.
        response := telebot.QueryResponse{
            Results:    results,
            IsPersonal: true,
        }

        // And finally send the response.
        if err := bot.AnswerInlineQuery(&query, &response); err != nil {
            log.WithFields(log.Fields {
		    "type": "error",
		    "query": query,
		    "error": err,
	    }).Error("Failed to respond to query:")
        }
    }
}

