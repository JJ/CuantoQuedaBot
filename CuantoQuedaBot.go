package main

// taken from https://github.com/tucnak/telebot
import (
    "log"
    "time"
    "os"
    "github.com/tucnak/telebot"
)

var bot *telebot.Bot

func main() {
    var err error
    bot, err = telebot.NewBot(os.Getenv("BOT_TOKEN"))
    if err != nil {
        log.Fatalln(err)
    }

    bot.Messages = make(chan telebot.Message, 1000)
    bot.Queries = make(chan telebot.Query, 1000)

    go messages()
    go queries()

    bot.Start(1 * time.Second)
}

func messages() {
    for message := range bot.Messages {
        log.Printf("Received a message from %s with the text: %s\n", message.Sender.Username, message.Text)
    }
}

func queries() {
    for query := range bot.Queries {
        log.Println("--- new query ---")
        log.Println("from:", query.From.Username)
        log.Println("text:", query.Text)

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
            log.Println("Failed to respond to query:", err)
        }
    }
}

