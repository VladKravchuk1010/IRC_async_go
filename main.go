package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

func randomStatus() float64 {
	rand.Seed(time.Now().UnixNano())
	return 100 + rand.Float64()*(1000-100)
}

func main() {
	r := gin.Default()

	r.POST("/set_status", func(c *gin.Context) {

		const secretToken = "SECRET_KEY1227"

		token := c.GetHeader("X-Async-Token")

		if token != secretToken {
			c.JSON(http.StatusForbidden, gin.H{"error": "Invalid token"})
			return
		}

		pk := c.PostForm("pk")
		if pk == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "ID (pk) is required"})
			return
		}

		go SendStatus(pk)

		c.JSON(http.StatusOK, gin.H{"message": "Расчет запущен (Go)"})
	})

	fmt.Println("🚀 Go сервис слушает порт :8081")
	r.Run(":8081")
}

func SendStatus(pk string) {
	fmt.Printf(" [ID %s] Расчет начат...\n", pk)
	time.Sleep(10 * time.Second)

	result := randomStatus()

	url := fmt.Sprintf("http://localhost:8000/api/reagent_calculations/%s/update_result/", pk)

	payload := map[string]interface{}{
		"result": result,
	}
	jsonData, _ := json.Marshal(payload)

	req, _ := http.NewRequest("PUT", url, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	req.Header.Set("X-Async-Token", "SECRET_KEY1227")

	client := &http.Client{}
	resp, err := client.Do(req)

	if err != nil {
		log.Printf("Ошибка отправки: %v", err)
		return
	}
	defer resp.Body.Close()
	fmt.Printf(" [ID %s] Статус обновлен в Django. Код: %d\n", pk, resp.StatusCode)
}
