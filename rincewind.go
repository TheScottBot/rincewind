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
	getKey()
	setupDefaults()

	fmt.Printf("%+v\n", translationRequest)

	fmt.Println("Calling API...")
	client := &http.Client{}

	form := url.Values{}
	form.Add("text", translationRequest.TranslateText)
	form.Add("source_lang", sourceOfDefault(translationRequest.SourceLanguage))
	form.Add("target_lang", targetOrDefault(translationRequest.TargetLanguage))

	req, err := http.NewRequest("POST", "https://api.deepl.com/v2/translate", strings.NewReader(form.Encode()))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", key)

	if err != nil {
		fmt.Print(err.Error())
		return TranslationResponse{}, err
	}

	resp, err := client.Do(req)

	if err != nil {
		fmt.Print(err.Error())
		return TranslationResponse{}, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		fmt.Print(err.Error())
		return TranslationResponse{}, err
	}

	var responseObject TranslationResponse
	json.Unmarshal(bodyBytes, &responseObject)
	fmt.Printf("API Response %+v\n", responseObject)

	return responseObject, err
}

func setupDefaults() {
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()
	defaultSource = viper.GetString("DefaultSource")
	defaultTarget = viper.GetString("DefaultTargetLang")

	fmt.Println("Defaults set " + defaultSource + " " + defaultTarget)
}

func sourceOfDefault(value string) string {
	return valueOrDefault(value, defaultSource)
}

func targetOrDefault(value string) string {
	return valueOrDefault(value, defaultTarget)
}

func valueOrDefault(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

func getKey() {
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()
	key = viper.GetString("Key")
}

type TranslationRequest struct {
	TranslateText  string
	SourceLanguage string
	TargetLanguage string
}

type TranslationResponse struct {
	Translations []translations `json:"translations"`
}

type translations struct {
	LanguageSource string `json:"detected_language_source"`
	Text           string `json:"text"`
}

var key, defaultSource, defaultTarget string
