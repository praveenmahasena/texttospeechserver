package transcribe

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/spf13/viper"
)

type Transcript struct {
	File      []byte
	Token     string
	Uploadurl string `json:"upload_url"`
	Text      string
	ID        string `json:"id"`
}

const (
	tokenParam = "TOKEN"
)

func New(file []byte) *Transcript {
	token := viper.GetString(tokenParam)

	t := Transcript{
		File:  file,
		Token: token,
	}

	return &t
}

func (t *Transcript) UploadMediaFile() error {

	req, reqErr := http.NewRequest("POST", "https://api.assemblyai.com/v2/upload", bytes.NewBuffer(t.File))

	if reqErr != nil {
		return reqErr
	}

	defer req.Body.Close()

	req.Header.Set("Authorization", t.Token)
	req.Header.Set("Content-Type", "application/octet-stream")

	client := &http.Client{}

	res, resErr := client.Do(req)

	if resErr != nil {
		fmt.Println(res)
		return resErr
	}

	defer res.Body.Close()

	err := json.NewDecoder(res.Body).Decode(&t)

	fmt.Println(t.Token)

	return err
}

func (t *Transcript) TranscribeRes() error {
	b, bErr := json.Marshal(map[string]string{
		"audio_url": t.Uploadurl,
	})

	if bErr != nil {
		return bErr
	}

	req, reqErr := http.NewRequest("POST", "https://api.assemblyai.com/v2/transcript", bytes.NewBuffer(b))

	req.Header.Set("Authorization", t.Token)
	req.Header.Set("Content-Type", "application/json")

	if reqErr != nil {
		return reqErr
	}

	defer req.Body.Close()

	client := &http.Client{}

	res, resErr := client.Do(req)

	if resErr != nil {
		return resErr
	}

	defer res.Body.Close()

	err := json.NewDecoder(res.Body).Decode(t)
	return err
}

func (t *Transcript) GetTransStr() (string, error) {
	req, reqErr := http.NewRequest("GET", "https://api.assemblyai.com/v2/transcript/"+t.ID, nil)

	if reqErr != nil {
		return "", reqErr
	}

	req.Header.Set("Authorization", t.Token)
	req.Header.Set("Content-Type", "application/json")

	ticker := time.NewTicker(time.Duration(5*2) * time.Second)

	for range ticker.C {
		client := http.Client{}
		res, resErr := client.Do(req)

		if resErr != nil {
			ticker.Stop()
			return "", resErr
		}

		defer res.Body.Close()

		result := map[string]any{}

		err := json.NewDecoder(res.Body).Decode(&result)

		if err != nil {
			ticker.Stop()
			return "", err
		}

		if result["status"] == "completed" {
			t.Text = result["text"].(string)
			ticker.Stop()
			return t.Text, nil
		}

	}

	return "", fmt.Errorf("server error")

}
