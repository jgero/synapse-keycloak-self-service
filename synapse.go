package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type SynapseConnection struct {
	Host        string
	Domain      string
	AccessToken string
}

func (conn *SynapseConnection) IsUsernameAvailable(name string) (bool, error) {
	qUrl := url.URL{
		Scheme:   "https",
		Host:     conn.Host,
		Path:     "/_synapse/admin/v1/username_available",
		RawQuery: url.Values{"username": {name}}.Encode(),
	}
	req, err := http.NewRequest(http.MethodGet, qUrl.String(), nil)
	if err != nil {
		return false, err
	}
	req.Header.Set("Authorization", "Bearer "+conn.AccessToken)
	hc := http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return false, err
	}
	if res.StatusCode != http.StatusOK {
		if body, err := io.ReadAll(res.Body); err != nil {
			return false, fmt.Errorf("unexpected status code: %s", res.Status)
		} else {
			return false, fmt.Errorf("unexpected status code: %s: %s", res.Status, string(body))
		}
	}
	body, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}
	var data struct {
		Available bool `json:"available"`
	}
	err = json.Unmarshal(body, &data)
	if err != nil {
		return false, err
	}
	return data.Available, nil
}

type createAccountParams struct {
	Displayname string `json:"displayname"`
	ExternalIds []struct {
		AuthProvider string `json:"auth_provider"`
		ExternalId   string `json:"external_id"`
	} `json:"external_ids"`
}

func (conn *SynapseConnection) CreateAccount(userInfo *UserInfo, selectedUsername string) error {
	ok, err := conn.IsUsernameAvailable(selectedUsername)
	if err != nil {
		return err
	} else if !ok {
		return fmt.Errorf("username %s is already taken", selectedUsername)
	}
	qUrl := url.URL{
		Scheme: "https",
		Host:   conn.Host,
		Path:   fmt.Sprintf("/_synapse/admin/v2/users/@%s:%s", selectedUsername, conn.Domain),
	}
	body := createAccountParams{
		Displayname: userInfo.Name,
		ExternalIds: []struct {
			AuthProvider string `json:"auth_provider"`
			ExternalId   string `json:"external_id"`
		}{
			{
				AuthProvider: "oidc-keycloak",
				ExternalId:   userInfo.Uuid,
			},
		},
	}
	b, err := json.Marshal(body)
	if err != nil {
		return err
	}
	req, err := http.NewRequest(http.MethodPut, qUrl.String(), bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+conn.AccessToken)
	req.Header.Set("Content-Type", "application/json")
	hc := http.Client{}
	res, err := hc.Do(req)
	if err != nil {
		return err
	}
	if res.StatusCode != http.StatusCreated {
		if body, err := io.ReadAll(res.Body); err != nil {
			return fmt.Errorf("unexpected status code: %s", res.Status)
		} else {
			return fmt.Errorf("unexpected status code: %s: %s", res.Status, string(body))
		}
	}
	return nil
}
