package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type KeycloakConnector struct {
	Host     string
	Realm    string
	ClientId string
}

func (conn *KeycloakConnector) GetAuthUrl() string {
	kUrl := url.URL{
		Scheme: "https",
		Host:   conn.Host,
		Path:   fmt.Sprintf("/realms/%s/protocol/openid-connect/auth", conn.Realm),
		RawQuery: url.Values{
			"scope":         {"openid"},
			"response_type": {"code"},
			"client_id":     {conn.ClientId},
			"redirect_uri":  {"/"},
		}.Encode(),
	}
	return kUrl.String()
}

type Token struct {
	ExpiresIn   int64  `json:"expires_in"`
	AccessToken string `json:"access_token"`
}

func (conn *KeycloakConnector) GetAccesToken(authCode string) (token *Token, err error) {
	kUrl := url.URL{
		Scheme: "https",
		Host:   conn.Host,
		Path:   fmt.Sprintf("/realms/%s/protocol/openid-connect/token", conn.Realm),
	}
	if err != nil {
		return
	}
	res, err := http.PostForm(kUrl.String(), url.Values{
		"grant_type":   {"authorization_code"},
		"code":         {authCode},
		"redirect_uri": {"/"},
		"client_id":    {"matrix-self-service"},
	})
	if err != nil {
		return
	}
	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error getting access token: status '%s'", res.Status)
	}
	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &token)
	if err != nil {
		return
	}
	return
}

type UserInfo struct {
	Uuid              string `json:"sub"`
	Name              string `json:"name"`
	PreferredUsername string `json:"preferred_username"`
}

func (conn *KeycloakConnector) GetUserinfo(token string) (data UserInfo, err error) {
	kUrl := url.URL{
		Scheme: "https",
		Host:   conn.Host,
		Path:   fmt.Sprintf("/realms/%s/protocol/openid-connect/userinfo", conn.Realm),
	}
	if err != nil {
		return
	}
	req, err := http.NewRequest(http.MethodPost, kUrl.String(), nil)
	if err != nil {
		return
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	hc := &http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return
	}
	defer res.Body.Close()
	var body []byte
	body, err = io.ReadAll(res.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return
	}
	return
}
