package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
)

type IntrospectResponse struct {
	Active     bool       `json:"active"`
	Credential Credential `json:"credential"`
}

type GetTokenResponse struct {
	Token string `json:"token"`
}

type Credential struct {
	Id         string `json:"id"`
	UserId     string `json:"user_id"`
	UserName   string `json:"user_name"`
	CreateDate string `json:"create_date"`
	ExpiresIn  int    `json:"expires_in"`
}

func SessionValidate(sessionId string) (*IntrospectResponse, error) {
	form := fmt.Sprintf("session_id=%s", sessionId)
	payload := strings.NewReader(form)
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/ac/session/validate", *ac), payload)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("Session validate endpoint response: %s", string(body))
	r := &IntrospectResponse{}
	_ = json.Unmarshal(body, r)
	return r, nil
}

func MakeGetTokenRequest(key string) (*GetTokenResponse, error) {
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/idc/session/fetch/%s/%s", *idc, key, *appId), nil)
	if err != nil {
		return nil, err
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, err
	}

	body, _ := io.ReadAll(resp.Body)
	log.Printf("IDC Token endpoint response: %s", string(body))
	r := &GetTokenResponse{}
	_ = json.Unmarshal(body, r)
	return r, nil
}
