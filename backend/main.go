package main

import (
	"encoding/json"
	"fmt"
	"github.com/spf13/viper"
	"log"
	"net/http"
)

type Config struct {
	KeaControlAgentURL string
}

var (
	config Config
)

func leases(w http.ResponseWriter, req *http.Request) {
	result := updateLeases()

	response, err := json.Marshal(result)
	if err != nil {
		log.Fatalf("Error marshalling JSON response: %v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*") // URL of the Kea control agent API endpoint
	w.Write(response)
}

func main() {
	viper.SetEnvPrefix("DHCPBROWSER")
	viper.BindEnv("KEACONTROLAGENTURL")
	viper.AutomaticEnv()

	viper.SetConfigName("config")
	viper.AddConfigPath(".")
	viper.SetConfigType("yaml")

	err := viper.ReadInConfig() // Find and read the config file

	err = viper.Unmarshal(&config)

	if err != nil {
		log.Fatalf("Error unmarshalling config: %v", err)
	}

	fmt.Println("Kea Control Agent URL: ", config.KeaControlAgentURL)

	http.HandleFunc("/api/ipv4-leases", leases)

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/", fs)

	http.ListenAndServe(":8090", nil)
}
