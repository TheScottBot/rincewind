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

func Translate(translationRequest TranslationRequest) (TranslationResponse, error) {
	GetKey()

	fmt.Printf("%+v\n", translationRequest)

	fmt.Println("Calling API...")
	client := &http.Client{}

	form := url.Values{}
	form.Add("text", translationRequest.TranslateText)
	form.Add("source_lang", translationRequest.SourceLanguage)
	form.Add("target_lang", translationRequest.TargetLanguage)

	req, err := http.NewRequest("POST", "https://api.deepl.com/v2/translate", strings.NewReader(form.Encode()))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", translationRequest.Key)

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

	return responseObject, err
}

func GetKey() {
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()
	Key = viper.GetString("Key")
}

type TranslationRequest struct {
	TranslateText  string
	SourceLanguage string
	TargetLanguage string
	Key            string
}

type TranslationResponse struct {
	Translations []Translations `json:"translations"`
}

type Translations struct {
	LanguageSource string `json:"detected_language_source"`
	Text           string `json:"text"`
}

var Key string
