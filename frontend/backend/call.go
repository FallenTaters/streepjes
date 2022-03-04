package backend

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/FallenTaters/streepjes/frontend/events"
)

func checkStatus(code int) error {
	if code == http.StatusOK {
		return nil
	}

	if code == http.StatusUnauthorized {
		events.Trigger(events.Unauthorized)
		return ErrUnauthorized
	}

	if code == http.StatusForbidden {
		return ErrForbidden
	}

	return fmt.Errorf(`%w: %d`, ErrStatus, code)
}

func get(path string, dst interface{}) error {
	resp, err := http.Get(settings.URL() + path)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return err
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}

func post(path string, payload interface{}, dst interface{}) error {
	var body []byte

	if payload != nil {
		data, err := json.Marshal(payload)
		if err != nil {
			return err
		}

		body = data
	}

	resp, err := http.Post(settings.URL()+path, ``, bytes.NewReader(body))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if err := checkStatus(resp.StatusCode); err != nil {
		return err
	}

	if dst == nil {
		return nil
	}

	return json.NewDecoder(resp.Body).Decode(dst)
}
