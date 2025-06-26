package dsl

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
)

const PUMPFUN_METADATA_URL = "https://pump.fun/api/ipfs"

type NewTokenRequest struct {
	FilePath    string
	Name        string
	Symbol      string
	Description string
	Twitter     string
	Telegram    string
	Website     string
}

func (request *NewTokenRequest) toRequest() (*bytes.Buffer, string, error) {
	if _, err := os.Stat(request.FilePath); errors.Is(err, os.ErrNotExist) {
		return nil, "", fmt.Errorf("file does not exist: %w", err)
	}

	var b bytes.Buffer
	writer := multipart.NewWriter(&b)

	part, err := writer.CreateFormFile("file", "image.png")
	if err != nil {
		return nil, "", err
	}

	file, err := os.Open(request.FilePath)
	if err != nil {
		return nil, "", err
	}
	defer file.Close()

	_, err = io.Copy(part, file)
	if err != nil {
		return nil, "", err
	}

	writer.WriteField("name", request.Name)
	writer.WriteField("symbol", request.Symbol)
	writer.WriteField("description", request.Description)
	writer.WriteField("twitter", request.Twitter)
	writer.WriteField("telegram", request.Telegram)
	writer.WriteField("website", request.Website)
	writer.WriteField("showName", "true")

	err = writer.Close()
	if err != nil {
		return nil, "", err
	}

	return &b, writer.FormDataContentType(), nil
}

type NewTokenResponse struct {
	Name        string `json:"name"`
	Symbol      string `json:"symbol"`
	Description string `json:"description"`
	ShowName    bool   `json:"showName"`
	CreatedOn   string `json:"createdOn"`
	Twitter     string `json:"twitter"`
	Telegram    string `json:"telegram"`
	Website     string `json:"website"`

	Image       string `json:"image"`
	MetadataUri string `json:"metadataUri"`
}

func Upload(client *http.Client, request NewTokenRequest) (*NewTokenResponse, error) {
	b, contentType, err := request.toRequest()
	if err != nil {
		return nil, fmt.Errorf("failed to create writer: %w", err)
	}

	req, err := http.NewRequest("POST", PUMPFUN_METADATA_URL, b)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", contentType)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	fmt.Println("response: ", string(body))

	var response NewTokenResponse
	json.Unmarshal(body, &response)

	return &response, nil
}
