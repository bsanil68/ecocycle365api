package controller

import (
	"ecocycleapis/logger"
	"ecocycleapis/utils"
	"encoding/json"
	"net/http"

	"context"

	"fmt"
	"io"

	"github.com/gin-gonic/gin"
	shell "github.com/ipfs/go-ipfs-api"
)

// Replace with your actual contract address
//var contractAddress = common.HexToAddress("0x48b4de6e7d269f3f3db63b640ff4920c30e094fd")

// DeviceData defines the structure for JSON payload
type ViewParams struct {
	DeviceID string `json:"deviceID"`
}
type JsonFilesData struct {
	FileData string `json:"filedata"`
}

// ParamData represents the structure with two parameters
type ParamData struct {
	Param1 string `json:"param1"`
	Param2 string `json:"param2"`
}

// StoreHandler handles the POST request to store device data on IPFS and the hash on the blockchain
func ViewStoredHashes(c *gin.Context) {

	// Bind JSON payload
	var viewparams ViewParams

	if err := c.BindJSON(&viewparams); err != nil {

		logger.Error("An error occurred bindjson:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	// Serialize JSON data for IPFS
	jsonData, err := json.Marshal(viewparams)
	if err != nil {
		logger.Error("An error occurred in marshalling:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize JSON: " + err.Error()})
		return
	}

	// Upload to IPFS
	ipfsHash, err := utils.UploadDataToIPFS(jsonData, deviceData.DeviceID)
	if err != nil {
		logger.Error("An error occurred in ipfshashing function:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to IPFS: " + err.Error()})
		return
	}

	transvalue, err := saveIOTdataInContract("https://testnet.hashio.io/api", ipfsHash, deviceData.DeviceID)
	if err != nil {
		logger.Error("An error occurred in storing to hedera  function:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add to  smart contract: " + transvalue})
		return
	}

}

// ReadJSONFromIPFS reads and parses a JSON file from IPFS containing two parameters
func ReadJSONFromIPFS(cidHash string, ipfsURL string) (*ParamData, error) {
	// Create a new IPFS shell instance
	sh := shell.NewShell(ipfsURL)
	if sh == nil {
		return nil, fmt.Errorf("failed to create IPFS shell")
	}

	// Create a context with timeout
	ctx := context.Background()

	// Get the file from IPFS
	reader, err := sh.Cat(cidHash)
	if err != nil {
		return nil, fmt.Errorf("failed to read from IPFS: %v", err)
	}
	defer reader.Close()

	// Read all content
	content, err := io.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read content: %v", err)
	}

	// Parse JSON
	var data ParamData
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &data, nil
}

// ReadJSONFromGateway reads and parses a JSON file from an IPFS gateway
func ReadJSONFromGateway(cidHash string) (*ParamData, error) {
	// Use public IPFS gateway
	gatewayURL := fmt.Sprintf("https://ipfs.io/ipfs/%s", cidHash)

	// Create HTTP client and make request
	resp, err := http.Get(gatewayURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch from gateway: %v", err)
	}
	defer resp.Body.Close()

	// Read response body
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %v", err)
	}

	// Parse JSON
	var data ParamData
	if err := json.Unmarshal(content, &data); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	return &data, nil
}
