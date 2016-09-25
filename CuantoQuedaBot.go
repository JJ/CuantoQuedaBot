package main

// taken from https://github.com/tucnak/telebot
import (

    "time"
    "os"

    "strings"
    "strconv" 

    "encoding/json"
    "io/ioutil"

    log "github.com/Sirupsen/logrus"
    "github.com/tucnak/telebot"
)

var bot *telebot.Bot
type Hito struct {
	File string `json:"file"`
	Title string `json:"title"`
	Date string `json:"fecha"`
}

type Data struct {
	Comment string `json:"comment"`
	Hitos []Hito `json:"hitos"`
}

var hitos []Hito
var results []telebot.InlineQueryResult
var ahora = time.Now()
var fechas []time.Time
 
func init() {
	// Log as JSON instead of the default ASCII formatter.
	log.SetFormatter(&log.JSONFormatter{})

	// Load milestones array
	file, e := ioutil.ReadFile("./hitos.json")
	if e != nil {
		log.WithFields(log.Fields{
			"error": e,
		}).Fatal("File error", e)
		os.Exit(1)
	}
	var hitos_data Data
	if err := json.Unmarshal(file,&hitos_data); err != nil {
		log.WithFields(log.Fields{
			"error": err,
		}).Fatal("JSON error")
	}

	for i,hito := range hitos_data.Hitos {
		this_url := strings.Join( []string{"https://jj.github.io/IV/documentos/proyecto/",hito.File}, "/")
		d := strings.Split(hito.Date,"/")
		this_year, _ := strconv.Atoi(d[0])
		this_month, _ := strconv.Atoi(d[1])
		this_day, _ := strconv.Atoi(d[2])
		fechas = append( fechas, 
			time.Date(this_day, time.Month(this_month), this_year,
				12,30,0,0, time.UTC))
		article := &telebot.InlineQueryResultArticle{
			Title: hito.Title,
			URL:   this_url,
			InputMessageContent: &telebot.InputTextMessageContent{
				Text:            strings.Join( []string{"Hito ", strconv.Itoa(i), ":", hito.Title, " =>", this_url }, " "),
				DisablePreview: false,
			},
		}
		results = append( results, article )

	}
//	fmt.Printf(" Results %v", results );
	
}

func main() {
    var err error
    if os.Getenv("BOT_TOKEN") == "" {
	    log.Fatal("No se ha definido el token del bot")
    }
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

