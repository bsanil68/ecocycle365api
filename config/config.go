package config

import "os"

// Config holds the configuration values
type Config struct {
	ContractAddress string
	PrivateKey      string
	WalletAddress   string
	RPCEndpoint     string
}

// LoadConfig initializes and returns the Config struct
func LoadConfig() *Config {
	return &Config{
		ContractAddress: os.Getenv("CONTRACT_ADDRESS"),
		PrivateKey:      os.Getenv("PRIVATE_KEY"),
		WalletAddress:   os.Getenv("WALLET_ADDRESS"),
		RPCEndpoint:     os.Getenv("RPC_ENDPOINT"),
	}
}
