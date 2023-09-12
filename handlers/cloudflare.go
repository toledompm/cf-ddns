package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type ICloudflare interface {
	UpdateRecord(ip string, recordName string, recordZone string, proxyEnabled bool, recordType string) error
}

func NewCloudflare(token string) ICloudflare {
	return &CloudflareHandler{
		headers: map[string][]string{
			"Content-Type":  {"application/json"},
			"Authorization": {fmt.Sprintf("Bearer %s", token)},
		},
	}
}

type CloudflareHandler struct {
	headers map[string][]string
}

func (c *CloudflareHandler) UpdateRecord(ip string, recordName string, recordZone string, proxyEnabled bool, recordType string) error {
	fmt.Printf("Updating record %s with ip %s in zone %s\n", recordName, ip, recordZone)

	recordId, recordValue, err := c.headRecord(recordName, recordZone, recordType)
	if err != nil {
		return err
	}

	if recordId == "" {
		fmt.Printf("New record detected, creating...\n")
		err = c.createNewRecord(ip, recordName, recordZone, proxyEnabled, recordType)
	} else if recordValue != ip {
		fmt.Printf("Record found in cloudflare, updating existing record\n")
		err = c.updateExistingRecord(ip, recordId, recordZone, proxyEnabled)
	}

	if err != nil {
		return err
	}

	return err
}

func (c *CloudflareHandler) headRecord(recordName string, recordZone string, recordType string) (id string, value string, err error) {
	requestURL := &url.URL{
		Scheme:   "https",
		Host:     "api.cloudflare.com",
		Path:     fmt.Sprintf("/client/v4/zones/%s/dns_records", recordZone),
		RawQuery: fmt.Sprintf("name=%s,type=%s", recordName, recordType),
	}

	request, err := http.NewRequest("GET", requestURL.String(), nil)
	if err != nil {
		return "", "", err
	}

	request.Header = c.headers

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", "", err
	}

	responseBodyBytes, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return "", "", err
	}

	response := &struct {
		Result []struct {
			Name    string `json:"name"`
			Content string `json:"content"`
			ID      string `json:"id"`
		} `json:"result"`
	}{}

	if err = json.Unmarshal(responseBodyBytes, response); err != nil {
		return "", "", err
	}

	if len(response.Result) == 1 {
		return response.Result[0].ID, response.Result[0].Content, nil
	}

	if len(response.Result) > 1 {
		err = fmt.Errorf("found %d records, expected 1", len(response.Result))
	}

	return "", "", err
}

func (c *CloudflareHandler) updateExistingRecord(ip string, recordId string, recordZone string, proxyEnabled bool) error {
	requestURL := &url.URL{
		Scheme: "https",
		Host:   "api.cloudflare.com",
		Path:   fmt.Sprintf("/client/v4/zones/%s/dns_records/%s", recordZone, recordId),
	}

	requestBody := map[string]interface{}{
		"content": ip,
		"proxied": proxyEnabled,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBodyBytes)

	request, err := http.NewRequest("PATCH", requestURL.String(), requestBodyBuffer)
	if err != nil {
		return err
	}

	request.Header = c.headers

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	responseBodyBytes, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return err
	}

	response := &struct {
		Errors  []map[string]interface{} `json:"errors"`
		Success bool                     `json:"success"`
	}{}

	if err = json.Unmarshal(responseBodyBytes, response); err != nil {
		return err
	}

	if !response.Success {
		errorJson, jsonErr := json.MarshalIndent(response.Errors, "", "  ")
		if jsonErr != nil {
			fmt.Printf("error marshalling cloudflare errors: %v", jsonErr)
			return err
		}
		err = fmt.Errorf("error updating record: %s", errorJson)
		return err
	}

	return nil
}

func (c *CloudflareHandler) createNewRecord(ip string, recordName string, recordZone string, proxyEnabled bool, recordType string) error {
	requestURL := &url.URL{
		Scheme: "https",
		Host:   "api.cloudflare.com",
		Path:   fmt.Sprintf("/client/v4/zones/%s/dns_records", recordZone),
	}

	requestBody := map[string]interface{}{
		"type":    recordType,
		"name":    recordName,
		"content": ip,
		"proxied": proxyEnabled,
	}

	requestBodyBytes, err := json.Marshal(requestBody)
	if err != nil {
		return err
	}

	requestBodyBuffer := bytes.NewBuffer(requestBodyBytes)

	request, err := http.NewRequest("POST", requestURL.String(), requestBodyBuffer)
	if err != nil {
		return err
	}

	request.Header = c.headers

	httpResponse, err := http.DefaultClient.Do(request)
	if err != nil {
		return err
	}

	responseBodyBytes, err := io.ReadAll(httpResponse.Body)

	if err != nil {
		return err
	}

	response := &struct {
		Errors  []map[string]interface{} `json:"errors"`
		Success bool                     `json:"success"`
	}{}

	if err = json.Unmarshal(responseBodyBytes, response); err != nil {
		return err
	}

	if !response.Success {
		errorJson, jsonErr := json.MarshalIndent(response.Errors, "", "  ")
		if jsonErr != nil {
			fmt.Printf("error marshalling cloudflare errors: %v", jsonErr)
			return err
		}
		err = fmt.Errorf("error updating record: %s", errorJson)
		return err
	}

	return nil
}
