package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/joho/godotenv"
)

type Media struct {
	Type    string `json:"type"`
	Mid     string `json:"mid"`
	Label   string `json:"label"`
	Port    int    `json:"port"`
	Pt      int    `json:"pt"`
	Codec   string `json:"codec"`
	Record  bool   `json:"record"`
	Recfile string `json:"recfile"`
}

type Payload struct {
	Janus       string  `json:"janus"`
	Transaction string  `json:"transaction"`
	AdminSecret string  `json:"admin_secret"`
	Id          int     `json:"id"`
	Type        string  `json:"type"`
	Description string  `json:"description"`
	Media       []Media `json:"media"`
}

func main() {
	//Load env
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	JANUS_ADMIN_URL := os.Getenv("JANUS_SERVER")
	JANUS_ADMIN_KEY := os.Getenv("JANUS_ADMIN_KEY")
	var STREAM_ID = 100

	payload := Payload{
		Janus:       "create",
		Transaction: fmt.Sprintf("%d", time.Now().UnixNano()),
		AdminSecret: JANUS_ADMIN_KEY,
		Id:          STREAM_ID,
		Type:        "rtp",
		Description: "Stream-100",
		Media: []Media{
			{
				Type:    "video",
				Mid:     "v100",
				Label:   "100",
				Port:    5100,
				Pt:      100,
				Codec:   "h264",
				Record:  true,
				Recfile: "/opt/janus/recordings/stream-5100-%Y%m%d%H%M%S.mjr",
			},
		},
	}

	data, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error encoding payload:", err)
		return
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	req, err := http.NewRequest("POST", JANUS_ADMIN_URL, bytes.NewBuffer(data))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error contacting Janus:", err)
		return
	}
	defer resp.Body.Close()

	// Handle HTTP status codes
	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("HTTP error %d: %s\n", resp.StatusCode, string(body))
		return
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		fmt.Println("Error decoding response:", err)
		return
	}

	fmt.Println("Janus response:", result)
}
