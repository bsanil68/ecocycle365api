package controller

import (
	"context"
	"ecocycleapis/logger"
	"ecocycleapis/utils"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"

	"net/http"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/gin-gonic/gin"
)

// Replace with your actual contract address
//var contractAddress = common.HexToAddress("0x48b4de6e7d269f3f3db63b640ff4920c30e094fd")

// DeviceData defines the structure for JSON payload
type DeviceData struct {
	DeviceID string `json:"deviceID"`
	Data     string `json:"data"`
}

// StoreHandler handles the POST request to store device data on IPFS and the hash on the blockchain
func StoreHandler(c *gin.Context) {

	// Bind JSON payload
	var deviceData DeviceData

	if err := c.BindJSON(&deviceData); err != nil {

		logger.Error("An error occurred bindjson:", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}

	// Serialize JSON data for IPFS
	jsonData, err := json.Marshal(deviceData)
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

// prepareTransactionOptions prepares transaction options for the Ethereum client
func saveIOTdataInContract(hederaurl string, ipfshash string, deviceid string) (string, error) {

	// Load private key (replace with actual private key)
	privateKeyHex := os.Getenv("Hedera_Priv_Key")
	logger.Info(" pvkeyhex:", privateKeyHex)
	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	logger.Info("priv :", privateKey)
	if err != nil {
		log.Fatalf("Failed to convert private key: %v", err)
		logger.Error("An error occurred in storing to hedera  function:", err)
	}

	// Connect to the hedera client (similar to connecting to Ethereum, adjust as necessary)
	client, err := ethclient.Dial(hederaurl) // Replace with actual Neon RPC URL
	logger.Info("urlhederaconnected:", hederaurl)
	if err != nil {
		log.Fatalf("Failed to connect to the hedera client: %v", err)
		logger.Error("An error occurred connecting hedera url:", err)
	}

	// Get chain ID
	chainID := big.NewInt(296) // For hedera testnet
	//chainID, err := client.NetworkID(nil)
	//if err != nil {
	//	log.Fatalf("Failed to get chain ID: %v", err)
	//}

	// Create transactor with private key and chain ID
	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
	logger.Info("authvalye :", auth.Value)
	if err != nil {
		log.Fatalf("Failed to create transactor: %v", err)
		logger.Error("An error occurred in creating transactor:", err)
	}
	nonce, err := client.PendingNonceAt(context.Background(), auth.From)
	if err != nil {
		log.Fatalf("Failed to fetch nonce: %v", err)
		logger.Error("An error occurred in fetching nonce:", err)
	}
	// Set gas price (optional, or it can be fetched dynamically)
	//auth.GasPrice = big.NewInt(20000000000) // 20 Gwei (example)
	//auth.GasLimit = uint64(300000)          // Gas limit (adjust as needed)

	// Smart contract address (replace with actual contract address)
	contractAddress := common.HexToAddress(os.Getenv("Hedera_Store_IOT"))
	logger.Info("contract is ", contractAddress)

	// Load ABI from file
	abiFilePath := filepath.Join("resources", "iotdevicedata.json")
	abiData, err := os.ReadFile(abiFilePath)
	logger.Info("abidata is ", abiData)
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
		logger.Error("An error occurred in parsing abi file:", err)

	}
	// ABI of the contract (replace with actual ABI of your contract)
	abiString := string(abiData)

	contractABI, err := abi.JSON(strings.NewReader(abiString))
	if err != nil {
		log.Fatalf("Failed to parse ABI: %v", err)
		logger.Error("An error occurred in reading abistring:", err)
	}

	// Prepare the data for the contract method call
	data, err := contractABI.Pack("storeHash", ipfshash, deviceid)
	if err != nil {
		log.Fatalf("Failed to pack data: %v", err)
		logger.Error("An error occurred in pack call of contractabi:", err)
	}

	// Correct usage
	value := big.NewInt(0)
	gasLimit := uint64(150000) // Adjust as needed based on testing
	gasPrice := big.NewInt(300000000000)

	tx := types.NewTransaction(
		nonce,
		contractAddress,
		value,
		gasLimit,
		gasPrice,
		data,
	)

	// Sign the transaction
	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatalf("Failed to sign transaction: %v", err)
		logger.Error("An error occurred in signing transaction:", err)
	}

	// Send the transaction to the network
	err = client.SendTransaction(context.TODO(), signedTx)
	if err != nil {
		log.Fatalf("Failed to send transaction: %v", err)
		logger.Error("An error occurred in sending transaction:", err)
	}

	// Output the transaction hash
	fmt.Printf("Transaction sent! Hash: %s\n", signedTx.Hash().Hex())

	logger.Info("transa sent and its hash:", signedTx.Hash().Hex())
	return signedTx.Hash().Hex(), nil
}

// StoreHandler handles the POST request to store device data on IPFS and the hash on the blockchain
func StoreDeviceMeasuerments(ctx *gin.Context) {

	// Controller holds the configuration for the controller

	// Bind JSON payload
	var deviceData DeviceData
	if err := ctx.BindJSON(&deviceData); err != nil {
		log.Printf("Error parsing JSON payload: %v", err)
		logger.Error("An error occurred in parsing json payload input:", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload: " + err.Error()})
		return
	}
	log.Println("JSON payload successfully parsed.")

	// Serialize JSON data for IPFS
	jsonData, err := json.Marshal(deviceData)
	if err != nil {
		log.Printf("Error serializing JSON: %v", err)
		logger.Error("An error occurred in serilaisng json:", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to serialize JSON: " + err.Error()})
		return
	}
	log.Println("JSON data serialized successfully.")
	logger.Info("JSON data serialized successfully.", jsonData)

	// Upload to IPFS
	ipfsHash, err := utils.UploadDataToIPFS(jsonData, deviceData.DeviceID)
	logger.Info("JSON data stored in ipfs successfully.", ipfsHash)
	if err != nil {
		logger.Error("JSON data error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload to IPFS: " + err.Error()})
		return
	}
	//   store  hash  and deviceid  in hedera  smart  contract
	transvalue, err := saveIOTdataInContract(os.Getenv("Hedera_Testnet_Endpoint"), ipfsHash, deviceData.DeviceID)
	logger.Info("JSON data stored in hedera successfully.", transvalue)
	if err != nil {
		logger.Error("hedera function call error", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save hash in contract: " + transvalue})
		return
	}
}
