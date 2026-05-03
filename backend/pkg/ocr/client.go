package ocr

import (
	"bytes"
	"encoding/json"
	"fmt"
	"invoice-ocr-backend/internal/models"
	"mime/multipart"
	"net/http"
	"time"
)

type Client struct {
	BaseURL string
	HTTPClient *http.Client
}

func NewClient(baseURL string, timeout time.Duration) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: timeout,
		},
	}
}

func (c *Client) Extract(data []byte, filename string) (*models.OCRResponse, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filename)
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(data); err != nil {
		return nil, err
	}
	writer.Close()

	req, err := http.NewRequest("POST", c.BaseURL+"/extract", body)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ocr service returned status: %s", resp.Status)
	}

	var result models.OCRResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode ocr response: %w", err)
	}

	return &result, nil
}
