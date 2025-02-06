package config

import (
	"api-gateway/internal/models"
	"encoding/json"
	"log"
	"os"
)

var CFG *models.Config

func LoadConfigs() {
	file, err := os.ReadFile("config.json")
	if err != nil {
		log.Fatalf("❌ Error reading config.json: %v", err)
	}

	err = json.Unmarshal(file, &CFG)
	if err != nil {
		log.Fatalf("❌ Error parsing config.json: %v", err)
	}

	log.Println("✅ Configs are ready!")
}
