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

var r *Rincewind

func init() {
	r = New()
}

func New() *Rincewind {
	r := new(Rincewind)
	getApiDetails(r)
	setupDefaults(r)
	return r
}

func Translate(re TranslationRequest) (TranslationResponse, error) { return r.Translate(re) }

func (r *Rincewind) Translate(translationRequest TranslationRequest) (TranslationResponse, error) {
	fmt.Printf("%+v\n", translationRequest)

	fmt.Println("Calling API...")
	client := &http.Client{}

	form := url.Values{}
	form.Add("text", translationRequest.TranslateText)
	form.Add("source_lang", sourceOfDefault(translationRequest.SourceLanguage))
	form.Add("target_lang", targetOrDefault(translationRequest.TargetLanguage))
	form.Add("split_sentences", translationRequest.SplitSentences)
	form.Add("preserve_formatting", translationRequest.PreserveFormatting)
	form.Add("formality", translationRequest.Formality)
	form.Add("glossary_id", translationRequest.GlossaryID)
	form.Add("tag_handling", translationRequest.TagHandling)
	form.Add("non_splitting_tags", translationRequest.NonSplittingTags)
	form.Add("outline_detection", translationRequest.OutlineDetection)
	form.Add("splitting_tags", translationRequest.SplittingTags)
	form.Add("ignore_tags", translationRequest.IgnoreTags)

	req, err := http.NewRequest("POST", r.apiEndPoint, strings.NewReader(form.Encode()))

	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("Authorization", r.apiKey)

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

	if responseObject.Translations[0].Text == "" {
		return TranslationResponse{}, err
	}
	return responseObject, err
}

func SetDefaultTarget(in string) { r.SetDefaultTarget(in) }

func (r *Rincewind) SetDefaultTarget(in string) {
	if in != "" {
		r.defaultTargetLanguage = in
	}
}

func SetDefaultSource(in string) { r.SetDefaultSource(in) }

func (r *Rincewind) SetDefaultSource(in string) {
	if in != "" {
		r.defaultSourceLanguage = in
	}
}

func setupDefaults(r *Rincewind) {
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()
	r.defaultSourceLanguage = viper.GetString("DefaultSource")
	r.defaultTargetLanguage = viper.GetString("DefaultTargetLang")

	fmt.Println("Defaults set " + r.defaultSourceLanguage + " " + r.defaultTargetLanguage)
}

func sourceOfDefault(value string) string {
	return valueOrDefault(value, r.defaultSourceLanguage)
}

func targetOrDefault(value string) string {
	return valueOrDefault(value, r.defaultTargetLanguage)
}

func valueOrDefault(value string, defaultValue string) string {
	if len(value) > 0 {
		return value
	}
	return defaultValue
}

func getApiDetails(r *Rincewind) {
	viper.SetConfigFile("config.json")
	viper.ReadInConfig()
	r.apiKey = viper.GetString("Key")
	r.apiEndPoint = viper.GetString("Endpoint")
	fmt.Println("Endpoint set " + r.apiEndPoint)
}

type Rincewind struct {
	defaultSourceLanguage string
	defaultTargetLanguage string
	apiEndPoint           string
	apiKey                string
}

type TranslationRequest struct {
	TranslateText      string
	SourceLanguage     string
	TargetLanguage     string
	SplitSentences     string
	PreserveFormatting string
	Formality          string
	GlossaryID         string
	TagHandling        string
	NonSplittingTags   string
	OutlineDetection   string
	SplittingTags      string
	IgnoreTags         string
}

type TranslationResponse struct {
	Translations []translations `json:"translations"`
}

type translations struct {
	LanguageSource string `json:"detected_language_source"`
	Text           string `json:"text"`
}
