package config

import (
	"log"
	"os"
	"time"
)

func CheckingConfigs(configChangeChan chan struct{}) {
	lastModifiedTime := time.Now()

	for {
		fileInfo, err := os.Stat("config.json")
		if err != nil {
			log.Println("❌ checking error config.json", err)
			time.Sleep(5 * time.Second)
			continue
		}

		if fileInfo.ModTime().After(lastModifiedTime) {
			lastModifiedTime = fileInfo.ModTime()
			log.Println("🔄 Config.json change detected, restarting...")
			configChangeChan <- struct{}{} // Отправляем сигнал
		}

		time.Sleep(5 * time.Second)
	}
}
