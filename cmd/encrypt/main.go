package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"the-engine/internal/config"
)

func main() {
	key := flag.String("key", "", "Environment variable key")
	value := flag.String("value", "", "Environment variable value")
	masterKey := flag.String("master-key", "", "Master encryption key (or set ENGINE_MASTER_KEY env var)")
	file := flag.String("file", "", "Secure environment file to update")
	flag.Parse()

	if *key == "" || *value == "" {
		fmt.Println("Usage: encrypt -key <key> -value <value> [-master-key <key>] [-file <file>]")
		os.Exit(1)
	}

	// Set master key from flag or environment
	if *masterKey != "" {
		os.Setenv("ENGINE_MASTER_KEY", *masterKey)
	}

	secure, err := config.NewSecureEnv()
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	encrypted, err := secure.Encrypt(*value)
	if err != nil {
		fmt.Printf("Error encrypting: %v\n", err)
		os.Exit(1)
	}

	if *file != "" {
		// Update secure environment file
		updateSecureFile(*file, *key, encrypted)
	} else {
		fmt.Printf("Encrypted value for %s: enc:%s\n", *key, encrypted)
	}
}

func updateSecureFile(path, key, encryptedValue string) {
	var envFile config.SecureEnvFile

	// Load existing file if it exists
	if data, err := os.ReadFile(path); err == nil {
		json.Unmarshal(data, &envFile)
	}

	// Initialize map if needed
	if envFile.Variables == nil {
		envFile.Variables = make(map[string]string)
	}

	// Update the value
	envFile.Variables[key] = encryptedValue

	// Write back to file
	data, err := json.MarshalIndent(envFile, "", "  ")
	if err != nil {
		fmt.Printf("Error marshaling: %v\n", err)
		os.Exit(1)
	}

	if err := os.WriteFile(path, data, 0600); err != nil {
		fmt.Printf("Error writing file: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Updated %s in %s\n", key, path)
}
