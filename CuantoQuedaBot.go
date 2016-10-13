package main

// taken from https://github.com/tucnak/telebot
import (

    "time"
    "os"
    "fmt"
    "strings"
    "strconv"
    "net/http"
    "crypto/tls"
    "encoding/json"
    "io/ioutil"
    "bytes"
    "math"

    "github.com/Sirupsen/logrus"
    "github.com/ripcurld00d/logrus-logzio-hook"
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
var opicionesText bytes.Buffer
var opicionesNumberText bytes.Buffer
var ahora = time.Now()
var fechas []time.Time
var log = logrus.New()

func init() {

	// Log as JSON instead of the default ASCII formatter.

	logrus.SetFormatter(&logrus.JSONFormatter{})
	name, _ := os.Hostname()
	// Declare logrus plugin
	if os.Getenv("LOGZ_TOKEN") != "" {
		fields := logrus.Fields{
			"ID": os.Getenv("LOGZ_TOKEN"),
			"Host": os.Getenv("HOST"),
			"Username": os.Getenv("USER"),
		}
		tr := &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		}
		httpsClient := &http.Client{Transport: tr}

		hook := logzio.New(os.Getenv("LOGZ_HOST"), "CuantoQuedaBot", fields)
		hook.SetClient(httpsClient)
		logrus.AddHook(hook)

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

		//adding opciones to choose 
		opicionesText.WriteString(strconv.Itoa(i) + "-" + hito.Title + "\r\n")
		opicionesNumberText.WriteString(strconv.Itoa(i) + "\r\n")

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
				Text: strings.Join( []string{"Hito ", strconv.Itoa(i), ":", hito.Title, " =>", this_url }, " "),
				DisablePreview: false,
			},
		}
		results = append( results, article )

	}
	log.Info("Opciones:\r\n"+ opicionesText.String())
//	fmt.Printf(" Results %v", results );

}

func botHito(context telebot.Context){
 	hito_n, err := strconv.Atoi(context.Args["n"])
 	if err!=nil {  
	 	//no tiene parametross
	 	log.Info("args blank")
	 	botOptions(context)
 	}else{
	 	hito, _ := results[hito_n].(*telebot.InlineQueryResultArticle)
    	text, _ := hito.InputMessageContent.(*telebot.InputTextMessageContent)
    	bot.SendMessage(context.Message.Chat, fmt.Sprintf("Hito %d\n\t %s", hito_n, text.Text ), nil)	
 	}   
}

func botHelp(context telebot.Context){
    bot.SendMessage(context.Message.Chat, "Órdenes:\n\t/hito <número> ⇒ Describe hito\n\t/cuanto_queda <número> ⇒ Horas hasta entrega", nil)
    botOptions(context)
}

func botOptions(context telebot.Context){
 	bot.SendMessage(context.Message.Chat, "Opiciones:\n\t/hito <número> ⇒ Describe hito\n\t/cuanto_queda <número> ⇒ Horas hasta entrega", nil)
    bot.SendMessage(context.Message.Chat, "elegir entre:\n\t"+opicionesText.String(), nil)
    bot.SendMessage(context.Message.Chat, "ejemplo : \"/hito 1\" o \"/cuanto_queda 1\"", nil)
}

func botCuantoQueda(context telebot.Context) {
	
	 hito_n, err := strconv.Atoi(context.Args["n"])
	 if err!=nil {  
	 	//no tiene parametross
	 	log.Info("args blank")
	 	botOptions(context)
	 }else{
	 	queda := fechas[hito_n].Sub(time.Now())
     	response := getResponse(hito_n, queda)
	 	bot.SendMessage(context.Message.Chat, response, nil)
	 }
	
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
    	//criando em start para dar opciones para o user
     bot.Handle("/start", func (context telebot.Context) {
     	botOptions(context)
    })
   
    // routes are compiled as regexps
     //adicionando +ajuda
    bot.Handle("/ayuda", func (context telebot.Context) {
	    botHelp(context)
    })

    // named groups found in routes will get injected in the controller as arguments
    bot.Handle("/hito (?P<n>[0-9]+)", func(context telebot.Context) {
	    botHito(context)
    })

    // when no arguments given, return all posible messages. 
    bot.Handle("/hito", func(context telebot.Context) {
	    for hito_n,_ := range results {
	    hito := results[hito_n].(*telebot.InlineQueryResultArticle)
	    text := hito.InputMessageContent.(*telebot.InputTextMessageContent)
	    bot.SendMessage(context.Message.Chat, fmt.Sprintf("Hito %d\n\t %s", hito_n, text.Text ), nil)
	    }
    })

    // named groups found in routes will get injected in the controller as arguments
    bot.Handle("/cuanto_queda (?P<n>[0-9]+)", func(context telebot.Context) {
	    botCuantoQueda(context)
    })

 	//blank path hito
    bot.Handle("/hito", func(context telebot.Context) {
		botCuantoQueda(context)
    })

    //blank path
    bot.Handle("/cuanto_queda", func(context telebot.Context) {
		botCuantoQueda(context)
    })
    //any search, all matchs in the end show options.
    bot.Handle("(([A-Za-z1234567890])+)", func(context telebot.Context) {
		log.Error("blank")
		botOptions(context)
    })
 
	// when no arguments given, return all posible messages. 
    bot.Handle("/cuanto_queda", func(context telebot.Context){
    	for hito_n,_ := range results {
    		queda := fechas[hito_n].Sub(time.Now())
    		response := getResponse(hito_n, queda)
    		bot.SendMessage(context.Message.Chat, response, nil)
    	}
    })

    bot.Messages = make(chan telebot.Message, 1000)
    bot.Queries = make(chan telebot.Query, 1000)

    go messages()
    go queries()

    bot.Start(1 * time.Second)
}

func getResponse(hito_n int, queda time.Duration) string {
	var response bytes.Buffer
	var string_hito string
	var string_tiempo string
	
	queda_minutos := queda.Minutes()
	
	if queda_minutos < 0 {
		string_hito = fmt.Sprintf("Hito %d finalizado hace ", hito_n)
		queda_minutos = queda_minutos*(-1)
	} else {
		string_hito = fmt.Sprintf("Hito %d :\n\tQuedan ", hito_n)
	}
	
	response.WriteString(string_hito)
	
	switch {
	case queda_minutos > 1440: // More than 1 day
		div := float64(math.Abs(queda.Hours()))/24.0
		dias := math.Floor(div)
		resto := div - dias
		
		div = resto * 24
		horas := math.Floor(div)
		resto = div - horas
		
		minutos := math.Floor(resto * 60)
		
		string_tiempo = fmt.Sprintf("%.0f días, %.0f horas y %.0f minutos.", dias, horas, minutos)
		
	case queda_minutos > 60: // More than 1 hour
		queda_horas := math.Abs(queda.Hours())
		horas := math.Floor(queda_horas)
		resto := queda_horas - horas
		minutos := math.Floor(resto * 60)
		
		string_tiempo = fmt.Sprintf("%.0f horas y %.0f minutos.", horas, minutos)
		
	default:
		minutos := math.Floor(math.Abs(queda.Minutes()))
		string_tiempo = fmt.Sprintf("%.0f minutos.", minutos)
	}
	
	response.WriteString(string_tiempo)
	
	return response.String()
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
