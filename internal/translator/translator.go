package translator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"tg_bot/config"
)

type YdxAITranslator struct {
	config *config.AITranslatorConfig
}

func New(config *config.AITranslatorConfig) *YdxAITranslator {
	return &YdxAITranslator{
		config: config,
	}
}

type TranslateRequest struct {
	FolderId           string   `json:"folderId"`
	Texts              []string `json:"texts"`
	TargetLanguageCode string   `json:"targetLanguageCode"`
}

type TranslateResponse struct {
	Translations []Translation `json:"translations"`
}

type Translation struct {
	Text                 string `json:"text"`
	DetectedLanguageCode string `json:"detectedLanguageCode"`
}

func (t *YdxAITranslator) TranslateText(text string, target string) Translation {
	url := t.config.URL
	bearerToken := t.config.APIKey
	folderID := t.config.FolderID

	payLoad := TranslateRequest{
		FolderId:           folderID,
		Texts:              []string{text},
		TargetLanguageCode: target,
	}
	jsonData, err := json.Marshal(&payLoad)
	if err != nil {
		fmt.Println("Error marshalling JSON:", err)
		os.Exit(1)
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Ошибка при создании запроса:", err)
		os.Exit(1)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+bearerToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response:", err)
		os.Exit(1)
	}

	trResp := TranslateResponse{}

	err = json.Unmarshal(body, &trResp)
	if err != nil {
		fmt.Println("Error unmarshalling response:", err)
		os.Exit(1)
	}

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Request succeeded")
	} else {
		fmt.Printf("Request failed with status code: %d\n", resp.StatusCode)
	}

	return trResp.Translations[0]

}
