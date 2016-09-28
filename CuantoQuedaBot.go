package main

// taken from https://github.com/tucnak/telebot
import (

    "time"
    "os"
    "fmt" 
    "strings"
    "strconv" 

    "encoding/json"
    "io/ioutil"

    "github.com/Sirupsen/logrus"   
    "github.com/bshuster-repo/logrus-logstash-hook"
    "gopkg.in/polds/logrus-papertrail-hook.v2"

    "github.com/JJ/telebot"
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
var log = logrus.New()

func init() {

	// Log as JSON instead of the default ASCII formatter.
	
	logrus.SetFormatter(&logrus.JSONFormatter{})
	name, _ := os.Hostname()
	// Declare logrus plugin
	if os.Getenv("LOGZ_TOKEN") != "" {

		hook, err := logrus_logstash.NewHookWithFields("https", os.Getenv("LOGZ_HOST"), "CuantoQuedaBot", logrus.Fields{
			"hostname":    name,
			"serviceName": "CuantoQuedaBot",
			"token": os.Getenv("LOGZ_TOKEN"),
		})
		if  err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Hook error")
		}
		log.Hooks.Add(hook)

	}

	if os.Getenv("PAPERTRAIL_HOST") != "" {

		udp_port, _ := strconv.Atoi(os.Getenv("PAPERTRAIL_PORT"))
		hook, err := logrus_papertrail.NewPapertrailHook(&logrus_papertrail.Hook{
			Host: os.Getenv("PAPERTRAIL_HOST"),
			Port: udp_port,
			Hostname: name,
			Appname: "CuantoQuedaBot",
		})		
		if  err != nil {
			log.WithFields(logrus.Fields{
				"error": err,
			}).Fatal("Hook error")
		}
		log.Hooks.Add(hook)
	}
	
	// Load milestones array
	file, e := ioutil.ReadFile("./hitos.json")
	if e != nil {
		log.WithFields(logrus.Fields{
			"error": e,
		}).Fatal("File error", e)
		os.Exit(1)
	}
	var hitos_data Data
	if err := json.Unmarshal(file,&hitos_data); err != nil {
		log.WithFields(logrus.Fields{
			"error": err,
		}).Fatal("JSON error")
	}

	for i,hito := range hitos_data.Hitos {
		this_url := strings.Join( []string{"https://jj.github.io/IV/documentos/proyecto/",hito.File}, "/")
		d := strings.Split(hito.Date,"/")
		this_day, _ := strconv.Atoi(d[0])
		this_month, _ := strconv.Atoi(d[1])
		this_year, _ := strconv.Atoi(d[2])
		fechas = append( fechas, 
			time.Date(this_year, time.Month(this_month), this_day,
				12,30,0,0, time.Local))
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

    // routes are compiled as regexps
    bot.Handle("/ayuda", func (context telebot.Context) {
	    bot.SendMessage(context.Message.Chat, "Órdenes:\n\t/hito <número> ⇒ Describe hito\n\t/cuanto_queda <número> ⇒ Horas hasta entrega", nil)
    })

    // named groups found in routes will get injected in the controller as arguments
    bot.Handle("/hito (?P<n>[0-9]+)", func(context telebot.Context) {
	    hito_n, _ := strconv.Atoi(context.Args["n"])
	    hito, _ := results[hito_n].(*telebot.InlineQueryResultArticle)
	    text, _ := hito.InputMessageContent.(*telebot.InputTextMessageContent)
	    bot.SendMessage(context.Message.Chat, fmt.Sprintf("Hito %d\n\t %s", hito_n, text.Text ), nil)
    })

    // named groups found in routes will get injected in the controller as arguments
    bot.Handle("/cuanto_queda (?P<n>[0-9]+)", func(context telebot.Context) {
	    hito_n, _ := strconv.Atoi(context.Args["n"])
	    queda := fechas[hito_n].Sub(time.Now())
	    bot.SendMessage(context.Message.Chat, fmt.Sprintf("Hito %d\n\t Quedan %f horas", hito_n, queda.Hours() ), nil)
    })

    bot.Messages = make(chan telebot.Message, 1000)
    bot.Queries = make(chan telebot.Query, 1000)

    go messages()
    go queries()

    bot.Start(1 * time.Second)
}

func messages() {
    for message := range bot.Messages {
	    if handler, args := bot.Route(&message); handler != nil {
		    handler(telebot.Context{Message: &message, Args: args})
	    }
	    log.WithFields(logrus.Fields {
		    "type": "message",
		    "username": message.Sender.Username,
		    "text": message.Text }).Info("Message received")
    }
}

func queries() {
    for query := range bot.Queries {

        log.WithFields(logrus.Fields {
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
            log.WithFields(logrus.Fields {
		    "type": "error",
		    "query": query,
		    "error": err,
	    }).Error("Failed to respond to query:")
        }
    }
}

