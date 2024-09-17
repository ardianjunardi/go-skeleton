package onesignal

import (
	"bytes"
	"encoding/json"
	"errors"
	"go-skeleton/lib/utils"
	"io/ioutil"
	"net/http"
)

type (
	// Config ...
	config struct {
		conf utils.Config
	}

	OneSignal interface {
		PushNotification(xPlayer, title, description string) error
	}
)

const (
	UrlHost = "https://onesignal.com/api/v1/notifications"
)

func New(conf utils.Config) OneSignal {
	return &config{conf: conf}
}

func (lib *config) PushNotification(xPlayer, title, description string) error {
	var (
		err error
	)

	bodyData := map[string]interface{}{
		"app_id":             lib.conf.GetString("onesignal.app_id"),
		"headings":           map[string]interface{}{"en": title},
		"contents":           map[string]interface{}{"en": description},
		"include_player_ids": []string{xPlayer},
	}

	payload, err := json.Marshal(bodyData)
	if err != nil {
		return err
	}

	// Populate Http Request
	requestData, err := http.NewRequest(http.MethodPost, UrlHost, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	requestData.Header.Set("Authorization", "Basic "+lib.conf.GetString("onesignal.app_key"))
	requestData.Header.Set("Content-Type", "application/json; charset=utf-8")

	client := &http.Client{}

	// Do Http Request
	resp, err := client.Do(requestData)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	_, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	if resp.StatusCode != 200 {
		return errors.New("OneSignal  error: " + resp.Status)
	}

	return nil
}
