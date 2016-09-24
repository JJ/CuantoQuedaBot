package main

// taken from https://github.com/tucnak/telebot
import (

    "time"
    "os"

    "github.com/op/go-logging"
    "github.com/tucnak/telebot"
)

var bot *telebot.Bot

var log = logging.MustGetLogger("example")

// Example format string. Everything except the message has a custom color
// which is dependent on the log level. Many fields have a custom output
// formatting too, eg. the time returns the hour down to the milli second.
var format = logging.MustStringFormatter(
    `%{color}%{time:15:04:05.000} %{shortfunc} ▶ %{level:.4s} %{id:03x}%{color:reset} %{message}`,
)

func main() {
    var err error
    bot, err = telebot.NewBot(os.Getenv("BOT_TOKEN"))
    if err != nil {
        log.Critical(err)
    }

    bot.Messages = make(chan telebot.Message, 1000)
    bot.Queries = make(chan telebot.Query, 1000)

    backend := logging.NewLogBackend(os.Stderr, "", 0)
    backend1Leveled := logging.AddModuleLevel(backend)
    backend1Leveled.SetLevel(logging.INFO, "")

    go messages()
    go queries()

    bot.Start(1 * time.Second)
}

func messages() {
    for message := range bot.Messages {
        log.Info("Received a message from %s with the text: %s\n", message.Sender.Username, message.Text)
    }
}

func queries() {
    for query := range bot.Queries {
        log.Info("--- new query ---")
        log.Info("from:", query.From.Username)
        log.Info("text:", query.Text)

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
            log.Info("Failed to respond to query:", err)
        }
    }
}

