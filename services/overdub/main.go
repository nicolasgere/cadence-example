package overdub

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

const url = "https://descriptapi.com/v1/overdub/"
const token = "YOUR_TOKEN"
const voiceId = "VOICE_ID"

type OverdubResponse struct {
	Id    string `json:"id"`
	State string `json:"state"`
	Url   string `json:"url"`
}

func StartOverdubWithApi(text string) (response *OverdubResponse, err error) {
	data := map[string]string{
		"text":     text,
		"voice_id": voiceId,
	}
	payload, err := json.Marshal(data)

	req, _ := http.NewRequest("POST", url+"generate_async", bytes.NewReader(payload))

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)
	if res.StatusCode != 201 {
		return nil, errors.New("Overdub request fail: " + string(body))
	}
	response = &OverdubResponse{}
	err = json.Unmarshal(body, response)
	if err != nil {
		return
	}
	return
}

func GetOverdubWithApi(id string) (response *OverdubResponse, err error) {
	req, _ := http.NewRequest("GET", url+"generate_async/"+id, nil)

	req.Header.Add("content-type", "application/json")
	req.Header.Add("authorization", "Bearer "+token)

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	err = json.Unmarshal(body, &response)
	if err != nil {
		return
	}
	return
}
