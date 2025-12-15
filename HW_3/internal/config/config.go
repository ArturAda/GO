package config

import "os"

type Config struct {
	HTTPAddress string
	StorageType string
	JSONPath    string
}

func CreateNewConfigFrom() Config {
	newAddress := os.Getenv("HTTP_ADDRESS")
	if newAddress == "" {
		newAddress = ":8080"
	}
	newStorageType := os.Getenv("STORAGE_TYPE")
	if newStorageType == "" {
		newStorageType = "memory"
	}
	newJSONPath := os.Getenv("JSON_PATH")
	if newJSONPath == "" {
		newJSONPath = "./all_balances.json"
	}
	return Config{HTTPAddress: newAddress, StorageType: newStorageType, JSONPath: newJSONPath}
}
