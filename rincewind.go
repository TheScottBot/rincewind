package rincewind

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/viper"
)

func CallApiFromWeb(w http.ResponseWriter, req *http.Request) {

	translationRequest := translationRequest{
		translateText:  req.FormValue("text"),
		sourceLanguage: req.FormValue("source"),
		targetLanguage: req.FormValue("target"),
	}

	callDeepLApi(w, translationRequest)
}

func callDeepLApi(w http.ResponseWriter, request translationRequest) {
	fmt.Println("Calling API...")
	client := &http.Client{}

	form := url.Values{}
	form.Add("text", request.translateText)
	form.Add("source_lang", request.sourceLanguage)
	form.Add("target_lang", request.targetLanguage)

	req, err := http.NewRequest("POST", "https://api.deepl.com/v2/translate", strings.NewReader(form.Encode()))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", Key)

	if err != nil {
		fmt.Print(err.Error())
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
	}

	var responseObject TranslationResponse
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response %+v\n", responseObject)

	w.Write(bodyBytes)
}

func GetKey() {
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()
	Key = viper.GetString("Key")
	Port = viper.GetString("Port")
}

type translationRequest struct {
	translateText  string
	sourceLanguage string
	targetLanguage string
}

type TranslationResponse struct {
	Translations []Translations `json:"translations"`
}

type Translations struct {
	LanguageSource string `json:"detected_language_source"`
	Text           string `json:"text"`
}

var Key string
var Port string
