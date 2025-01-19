package utils

import (
	"bytes"

	"ecocycleapis/logger"

	"encoding/json"
	"errors"
	"fmt"
	"io"

	"mime/multipart"
	"net/http"
)

func UploadDataToIPFS(data []byte, deviceID string) (string, error) {
	// Create a buffer to hold the form data
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	// Create a form file with the provided device ID
	formFile, err := writer.CreateFormFile("file", deviceID)
	if err != nil {
		logger.Error("Unable to create form file", err)
		return "", fmt.Errorf("failed to create form file: %w", err)
	}

	// Write data to the form file
	if _, err = formFile.Write(data); err != nil {
		logger.Error("Failed to write data to form file", err)
		return "", fmt.Errorf("failed to write data: %w", err)
	}

	// Finalize the form writer
	if err = writer.Close(); err != nil {
		logger.Error("Failed to close writer", err)
		return "", fmt.Errorf("failed to close writer: %w", err)
	}

	// IPFS gateway URL
	url := "http://localhost:5001/api/v0/add"

	// Perform POST request
	resp, err := http.Post(url, writer.FormDataContentType(), body)
	if err != nil {
		logger.Error("Failed to send POST request to IPFS", err)
		return "", fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != http.StatusOK {
		logger.Error("IPFS gateway returned an error", err)
		return "", fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	// Parse response body
	var buffer bytes.Buffer
	if _, err = io.Copy(&buffer, resp.Body); err != nil {
		logger.Error("Failed to read response body", err)
		return "", fmt.Errorf("failed to read response body: %w", err)
	}

	// Extract CID (assuming response is JSON)
	var result map[string]interface{}
	if err = json.Unmarshal(buffer.Bytes(), &result); err != nil {
		logger.Error("Failed to parse response JSON", err)
		return "", fmt.Errorf("failed to parse response JSON: %w", err)
	}

	// Ensure CID exists
	cid, ok := result["Hash"].(string)
	if !ok {
		logger.Error("CID not found in IPFS response")
		return "", errors.New("CID not found in IPFS response")
	}

	logger.Info("Data uploaded to IPFS", cid)
	return cid, nil
}
