package main

// taken from https://github.com/tucnak/telebot
import (

    "time"
    "os"
    "fmt" 
    "strings"
    "strconv" 
    "log"
    "encoding/json"
    "io/ioutil"



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

func init() {

	// Load milestones array
	file, e := ioutil.ReadFile("./hitos.json")
	if e != nil {
		log.Fatal("File error", e)
		os.Exit(1)
	}
	var hitos_data Data
	if err := json.Unmarshal(file,&hitos_data); err != nil {
		log.Fatal("JSON error")
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

	
}

func main() {
	fmt.Printf(" Queda %d", fechas[0].Sub(time.Now()).Hours() );
}
